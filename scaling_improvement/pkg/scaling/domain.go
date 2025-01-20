package scaling

type DomainID string
type Nodes map[NodeName]*Node

type Domain struct {
	AllNodes          Nodes
	ActiveNodes       Nodes
	InactiveNodes     Nodes
	DomainID          DomainID
	AllocatedServices Services
}

type Domains map[DomainID]*Domain

func NewDomain(nodes Nodes, domainID DomainID) *Domain {
	inactiveNodes := make(Nodes)
	for nodeN := range nodes {
		inactiveNodes[nodeN+"r"] = NewNode(nodes[nodeN].Cores, nodes[nodeN].ReallocHeuristic, nodeN+"r")
	}
	return &Domain{
		AllNodes:          nodes,
		ActiveNodes:       nodes,
		InactiveNodes:     inactiveNodes,
		DomainID:          domainID,
		AllocatedServices: nil,
	}
}
