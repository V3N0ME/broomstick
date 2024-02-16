package main

import (
	"broomstick/controller"
	"broomstick/scheduler"
	"broomstick/store"
	"log"
	"os"
	"os/signal"
	"syscall"

	"broomstick/logger"
)

type ClusterInfo struct {
	ServiceAccount string `yaml:"service_account"`
	ProjectID      string `yaml:"project_id"`
	Location       string `yaml:"location"`
	Cluster        string `yaml:"cluster"`
}

type Configuration struct {
	ApiVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		ID             string           `yaml:"id"`
		ServiceAccount string           `yaml:"service_account"`
		ProjectID      string           `yaml:"project_id"`
		Location       string           `yaml:"location"`
		Cluster        string           `yaml:"cluster"`
		Schedule       scheduler.Config `yaml:"schedule"`
	} `yaml:"clusters"`
}

type broomstick struct {
	clusterInfo map[string]ClusterInfo
	taskQueue   chan task
	logger      logger.Logger
	store       controller.ClusterStore
}

func NewBroomstick(config Configuration) *broomstick {

	clusterMap := map[string]ClusterInfo{}
	taskQueue := make(chan task)

	logger := logger.NewLogger()

	for _, clusterInfo := range config.Clusters {

		logger.Info("[CONFIG]", clusterInfo.ID)

		clusterMap[clusterInfo.ID] = ClusterInfo{
			ServiceAccount: clusterInfo.ServiceAccount,
			ProjectID:      clusterInfo.ProjectID,
			Location:       clusterInfo.Location,
			Cluster:        clusterInfo.Cluster,
		}

		clusterID := clusterInfo.ID

		scheduler.Schedule(clusterInfo.Schedule, func() {

			taskQueue <- task{
				ID:     clusterID,
				Action: "ON",
			}

		}, func() {

			taskQueue <- task{
				ID:     clusterID,
				Action: "OFF",
			}
		})
	}

	return &broomstick{
		clusterInfo: clusterMap,
		taskQueue:   taskQueue,
		logger:      logger,
		store:       store.NewConfigStore(),
	}
}

type task struct {
	ID     string
	Action string
}

type RunConfig struct {
	ID     string
	Action string
}

func (b *broomstick) Run(config RunConfig) {

	b.logger.Info("[RUN]", config.ID)

	clusterInfo := b.clusterInfo[config.ID]

	cluster, err := controller.NewGCPClusterController(controller.GCPClusterConfig{
		ProjectID:            clusterInfo.ProjectID,
		Location:             clusterInfo.Location,
		Cluster:              clusterInfo.Cluster,
		Base64ServiceAccount: clusterInfo.ServiceAccount,
		Store:                b.store,
		Logger:               b.logger,
	})
	if err != nil {
		b.logger.Error("[CLUSTER]", config.ID, err)
		return
	}

	switch config.Action {
	case "ON":

		err = cluster.TurnOn()

	case "OFF":

		err = cluster.TurnOff()
	}

	if err == nil {
		b.logger.Info("[CLUSTER]", config.ID, config.Action, "[COMPLETED]")
	}
}

func (b *broomstick) Schedule() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGTERM, syscall.SIGINT)

	signal := scheduler.Start()

	go func() {
		<-signalChan
		log.Println("Exiting")
		signal <- true
		close(b.taskQueue)
	}()

	for task := range b.taskQueue {

		b.logger.Info("[CRON]", task.ID)

		clusterInfo := b.clusterInfo[task.ID]

		cluster, err := controller.NewGCPClusterController(controller.GCPClusterConfig{
			ProjectID:            clusterInfo.ProjectID,
			Location:             clusterInfo.Location,
			Cluster:              clusterInfo.Cluster,
			Base64ServiceAccount: clusterInfo.ServiceAccount,
			Store:                b.store,
			Logger:               b.logger,
		})
		if err != nil {
			b.logger.Error("[CLUSTER]", task.ID, err)
			continue
		}

		switch task.Action {
		case "ON":

			err = cluster.TurnOn()

		case "OFF":

			err = cluster.TurnOff()
		}

		if err == nil {
			b.logger.Info("[CLUSTER]", task.ID, task.Action, "[COMPLETED]")
		}
	}
}
