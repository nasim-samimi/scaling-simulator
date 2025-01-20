package scaling

type Cloud struct {
	AllNodes      Nodes
	ActiveNodes   Nodes
	InactiveNodes Nodes
}

func NewCloud(nodes Nodes) *Cloud {
	return &Cloud{
		AllNodes:      nodes,
		ActiveNodes:   make(Nodes),
		InactiveNodes: nodes,
	}
}
