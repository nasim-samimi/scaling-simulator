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

func NewDomain(nodes Nodes, reservedNodes Nodes, domainID DomainID) *Domain {
	// fmt.Println("Active Nodes in new domain: ", nodes)
	// fmt.Println("Reserved Nodes in new domain: ", reservedNodes)
	// fmt.Println("Domain ID in new domain: ", domainID)
	return &Domain{
		AllNodes:          nodes,
		ActiveNodes:       nodes,
		InactiveNodes:     reservedNodes,
		DomainID:          domainID,
		AllocatedServices: nil,
	}
}
