package scaling

type DomainID string
type Nodes map[NodeName]*Node

type Domain struct {
	AllNodes       Nodes
	OnNodes        Nodes
	DomainID       DomainID
	AllocatedTasks []*Task
}

type Domains map[DomainID]*Domain
