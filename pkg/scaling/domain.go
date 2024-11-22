package scaling

type DomainID string
type Nodes map[NodeName]*Node

type Domain struct {
	AllNodes          Nodes
	OnNodes           Nodes
	DomainID          DomainID
	AllocatedServices Services
}

type Domains map[DomainID]*Domain

func NewDomain(nodes Nodes, domainID DomainID) *Domain {
	return &Domain{
		AllNodes:          nodes,
		OnNodes:           make(Nodes),
		DomainID:          domainID,
		AllocatedServices: nil,
	}
}
