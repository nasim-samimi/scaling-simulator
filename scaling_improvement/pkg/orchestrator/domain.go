package orchestrator

type DomainID string
type Nodes map[NodeName]*Node

type Domain struct {
	AllNodes          Nodes
	ActiveNodes       Nodes
	InactiveNodes     Nodes
	DomainID          DomainID
	AllocatedServices Services
	AlwaysActiveNodes []NodeName
}

type Domains map[DomainID]*Domain

func NewDomain(nodes Nodes, reservedNodes Nodes, domainID DomainID) *Domain {
	// fmt.Println("Active Nodes in new domain: ", nodes)
	// fmt.Println("Reserved Nodes in new domain: ", reservedNodes)
	// fmt.Println("Domain ID in new domain: ", domainID)
	AlwaysActiveNodes := make([]NodeName, 0)
	for nodeName := range nodes {
		AlwaysActiveNodes = append(AlwaysActiveNodes, nodeName)
	}
	return &Domain{
		AllNodes:          nodes,
		ActiveNodes:       nodes,
		InactiveNodes:     reservedNodes,
		DomainID:          domainID,
		AllocatedServices: nil,
		AlwaysActiveNodes: AlwaysActiveNodes,
	}
}
