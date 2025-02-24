package scaling

import "fmt"

type NodeName string
type AllocatedServices map[ServiceID]*Service
type NodeStatus string

const (
	Active   NodeStatus = "Active"
	Inactive NodeStatus = "Inactive"
)

type Node struct {
	Cores                    Cores
	ReallocHeuristic         Heuristic
	NodeName                 NodeName
	NodeAdmission            *AdmissionTest
	Location                 Location
	DomainID                 DomainID
	AllocatedServices        AllocatedServices
	AverageResidualBandwidth float64
	TotalResidualBandwidth   float64
	Status                   NodeStatus
}

func NewNode(cores Cores, heuristic Heuristic, nodeName NodeName) *Node {
	admissionTest := NewAdmissionTest(cores, heuristic)
	// fmt.Println("Node Admission Test Cores: ", admissionTest.Cores[CoreID("core-0")])
	return &Node{
		Cores:                    cores,
		ReallocHeuristic:         heuristic,
		NodeName:                 nodeName,
		NodeAdmission:            admissionTest,
		Status:                   Inactive,
		AverageResidualBandwidth: 0,
		TotalResidualBandwidth:   0,
		AllocatedServices:        make(AllocatedServices),
	}
}

func (n *Node) NodeAllocate(reqCpus uint64, reqBandwidth float64, service *Service, eventID ServiceID) ([]CoreID, error) {
	selectedCpus, err := n.NodeAdmission.Admission(reqCpus, reqBandwidth, n.Cores)
	fmt.Println("Selected CPUs: ", selectedCpus)
	if err != nil {
		return selectedCpus, err
	}
	totalResidualBandwidth := 0.0
	for _, coreID := range selectedCpus {
		core := n.Cores[coreID]
		core.ConsumedBandwidth += reqBandwidth
		n.Cores[coreID] = core
	}
	for _, core := range n.Cores {
		totalResidualBandwidth += core.ConsumedBandwidth
	}
	fmt.Println("Service: ", service)
	// service.AllocatedNodeEdge = n.NodeName
	fmt.Println("in node allocation, total residual bandwidth and average: ", n.TotalResidualBandwidth, n.AverageResidualBandwidth)
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))
	fmt.Println("in node allocation, total residual bandwidth and average after update: ", n.TotalResidualBandwidth, n.AverageResidualBandwidth)
	return selectedCpus, nil
}

func (n *Node) CloudNodeDeallocate(eventID ServiceID) bool {
	cores := n.AllocatedServices[eventID].AllocatedCoresCloud

	bandwidth := n.AllocatedServices[eventID].ReducedMode.bandwidthCloud

	totalResidualBandwidth := n.TotalResidualBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		totalResidualBandwidth -= bandwidth
	}
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))

	delete(n.AllocatedServices, eventID)
	return true
}

func IntraDomainReallocateTest(newService *Service, oldServiceID ServiceID, n Node) (bool, error) {
	fmt.Println("Intra Domain Reallocate Test")
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	n.Cores = NewCores
	fmt.Println("old service ID:", oldServiceID)
	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.bandwidthEdge
	fmt.Println("Old Service Cores: ", oldServiceCores)
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth
	}
	_, err := n.NodeAdmission.Admission(newService.StandardMode.cpusEdge, newService.StandardMode.bandwidthEdge, NewCores)
	for _, core := range NewCores {
		fmt.Println("cores after intra domain realocate test newcores:", core)
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func IntraNodeReallocateTest(newService *Service, oldEventID ServiceID, n Node) (bool, error) {
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	n.Cores = NewCores

	oldServiceCores := n.AllocatedServices[oldEventID].AllocatedCoresEdge
	oldBandwidth := n.AllocatedServices[oldEventID].StandardMode.bandwidthEdge
	newBandwidth := newService.StandardMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= oldBandwidth
	}

	possibleCores, err := n.NodeAdmission.Admission(newService.StandardMode.cpusEdge, newService.StandardMode.bandwidthEdge, NewCores)
	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		_, err = n.NodeAdmission.Admission(n.AllocatedServices[oldEventID].StandardMode.cpusEdge, oldBandwidth, NewCores)
		if err == nil {
			return true, nil
		}
	}
	return false, err
}

func (n *Node) NodeDeallocate(eventID ServiceID) bool {
	fmt.Println("going to Node Deallocate")
	fmt.Println("allocated services:", n.AllocatedServices)
	fmt.Println("eventID:", eventID)
	cores := n.AllocatedServices[eventID].AllocatedCoresEdge
	fmt.Println("cores for deallocating in node deallocate:", cores)
	mode := n.AllocatedServices[eventID].AllocationMode
	fmt.Println("deallocation mode: ", mode)
	bandwidth := n.AllocatedServices[eventID].StandardMode.bandwidthEdge
	switch mode {
	case StandardMode:
		bandwidth = n.AllocatedServices[eventID].StandardMode.bandwidthEdge
	case ReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthEdge

	}
	totalResidualBandwidth := n.TotalResidualBandwidth
	for _, core := range cores {
		fmt.Println("core:", n.Cores[core].ConsumedBandwidth)
		fmt.Println("coreID:", core)
		n.Cores[core].ConsumedBandwidth -= bandwidth
		fmt.Println("consumed bandwidth after deallocation: ", n.Cores[core].ConsumedBandwidth)
		totalResidualBandwidth -= bandwidth
	}
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))

	delete(n.AllocatedServices, eventID)
	return true
}
