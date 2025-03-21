package orchestrator

type Cloud struct {
	AllNodes      Nodes
	ActiveNodes   Nodes
	InactiveNodes Nodes
}

func NewCloud(nodes Nodes, reservedNodes Nodes) *Cloud {
	return &Cloud{
		AllNodes:      nodes,
		ActiveNodes:   make(Nodes),
		InactiveNodes: reservedNodes,
	}
}
