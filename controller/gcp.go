package controller

import (
	"broomstick/logger"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"google.golang.org/api/option"
)

type gcp struct {
	projectID string
	location  string
	cluster   string
	client    *container.ClusterManagerClient
	store     ClusterStore
	logger    logger.Logger
}

type GCPClusterConfig struct {
	ProjectID            string
	Location             string
	Cluster              string
	Base64ServiceAccount string
	Store                ClusterStore
	Logger               logger.Logger
}

func NewGCPClusterController(config GCPClusterConfig) (ClusterController, error) {

	data, err := base64.StdEncoding.DecodeString(config.Base64ServiceAccount)
	if err != nil {
		config.Logger.Error("[CONTROLLER]", config.ProjectID, config.Cluster, err)
		return nil, err
	}

	c, err := container.NewClusterManagerClient(context.TODO(), option.WithCredentialsJSON(data))
	if err != nil {
		config.Logger.Error("[CONTROLLER]", config.ProjectID, config.Cluster, err)
		return nil, err
	}

	return &gcp{
		projectID: config.ProjectID,
		location:  config.Location,
		cluster:   config.Cluster,
		client:    c,
		store:     config.Store,
		logger:    config.Logger,
	}, nil
}

func (g *gcp) Close() {
	g.client.Close()
}

func (g *gcp) TurnOff() error {

	ctx := context.TODO()

	pools, err := g.client.ListNodePools(ctx, &containerpb.ListNodePoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", g.projectID, g.location, g.cluster),
	})
	if err != nil {
		g.logger.Error("[CONTROLLER]", g.projectID, g.cluster, err)
		return err
	}

	for _, nodePool := range pools.NodePools {

		if nodePool.InitialNodeCount == 0 {
			continue
		}

		g.logger.Info("[CONTROLLER]", g.cluster, nodePool.Name, "turning off")

		if err := g.store.StoreNodePoolInfo(NodePoolInfo{
			NodePool:  nodePool.Name,
			ProjectID: g.projectID,
			Cluster:   g.cluster,
		}, int(nodePool.InitialNodeCount)); err != nil {
			g.logger.Error("[CONTROLLER]", g.projectID, g.cluster, err)
			continue
		}

		operation, err := g.client.SetNodePoolSize(ctx, &containerpb.SetNodePoolSizeRequest{
			NodeCount: 0,
			Name:      fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", g.projectID, g.location, g.cluster, nodePool.GetName()),
		})
		if err != nil {

			if strings.Contains(err.Error(), "incompatible operation") {
				continue
			}

			g.logger.Error("[CONTROLLER]", g.cluster, nodePool.Name, err)
			continue
		}

		for operation.GetStatus() != containerpb.Operation_DONE && operation.GetStatus() != containerpb.Operation_ABORTING {

			time.Sleep(3 * time.Second)

			operation, err = g.client.GetOperation(ctx, &containerpb.GetOperationRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", g.projectID, g.location, operation.GetName()),
			})
			if err != nil {
				g.logger.Error("[CONTROLLER]", g.cluster, nodePool.Name, err)
				return err
			}
		}
	}

	return nil
}

func (g *gcp) TurnOn() error {

	ctx := context.TODO()

	pools, err := g.client.ListNodePools(ctx, &containerpb.ListNodePoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", g.projectID, g.location, g.cluster),
	})
	if err != nil {
		g.logger.Error("[CONTROLLER]", g.cluster, err)
		return err
	}

	for _, nodePool := range pools.NodePools {

		g.logger.Info("[CONTROLLER]", g.cluster, nodePool.Name, "turning on")

		size, err := g.store.GetNodePoolInfo(NodePoolInfo{
			NodePool:  nodePool.GetName(),
			ProjectID: g.projectID,
			Cluster:   g.cluster,
		})
		if err != nil {
			g.logger.Error("[CONTROLLER]", g.cluster, nodePool.Name, err)
			continue
		}

		operation, err := g.client.SetNodePoolSize(ctx, &containerpb.SetNodePoolSizeRequest{
			NodeCount: int32(size),
			Name:      fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", g.projectID, g.location, g.cluster, nodePool.GetName()),
		})
		if err != nil {
			g.logger.Error("[CONTROLLER]", g.cluster, nodePool.Name, err)
			continue
		}

		for operation.GetStatus() != containerpb.Operation_DONE && operation.GetStatus() != containerpb.Operation_ABORTING {

			time.Sleep(3 * time.Second)

			operation, err = g.client.GetOperation(ctx, &containerpb.GetOperationRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", g.projectID, g.location, operation.GetName()),
			})
			if err != nil {
				g.logger.Error("[CONTROLLER]", g.cluster, nodePool.Name, err)
				return err
			}
		}
	}

	return nil
}
