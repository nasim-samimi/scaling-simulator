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
	return &Domain{
		AllNodes:          nodes,
		ActiveNodes:       make(Nodes),
		InactiveNodes:     nodes,
		DomainID:          domainID,
		AllocatedServices: nil,
	}
}
