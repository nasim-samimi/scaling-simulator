package scaling

type Cloud struct {
	AllNodes Nodes
	OnNodes  Nodes
}

func NewCloud(nodes Nodes) *Cloud {
	return &Cloud{
		AllNodes: nodes,
		OnNodes:  make(Nodes),
	}
}
