package controller

type ClusterController interface {
	TurnOff() error
	TurnOn() error
	Close()
}

type NodePoolInfo struct {
	NodePool  string
	ProjectID string
	Cluster   string
}

type ClusterStore interface {
	StoreNodePoolInfo(NodePoolInfo, int) error
	GetNodePoolInfo(NodePoolInfo) (int, error)
}
