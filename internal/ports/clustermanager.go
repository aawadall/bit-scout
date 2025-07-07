package ports

// ClusterManagerPort defines the interface for cluster management (driven port, optional)
type ClusterManagerPort interface {
	RegisterNode(nodeID string, address string) error
	DeregisterNode(nodeID string) error
	ListNodes() ([]string, error)
}
