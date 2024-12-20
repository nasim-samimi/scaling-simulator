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

func (n *Node) NodeAllocate(reqCpus uint64, reqBandwidth float64, service *Service) ([]CoreID, error) {
	selectedCpus, err := n.NodeAdmission.Admission(reqCpus, reqBandwidth, n.Cores)
	fmt.Println("Selected CPUs: ", selectedCpus)
	if err != nil {
		return selectedCpus, err
	}
	for _, coreID := range selectedCpus {
		core := n.Cores[coreID]
		core.ConsumedBandwidth += reqBandwidth
		n.Cores[coreID] = core
	}
	fmt.Println("Service: ", service)
	n.AllocatedServices[service.serviceID] = service
	n.AverageResidualBandwidth = n.AverageResidualBandwidth + reqBandwidth*float64(reqCpus)/float64(len(n.Cores))
	n.TotalResidualBandwidth = n.TotalResidualBandwidth + reqBandwidth*float64(reqCpus)
	return selectedCpus, nil
}

func (n *Node) IntraDomainReallocateTest(newService *Service, oldServiceID ServiceID) (bool, error) {
	NewCores := n.Cores
	oldServiceCores := n.AllocatedServices[oldServiceID].allocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].standardMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth
	}

	_, err := n.NodeAdmission.Admission(newService.standardMode.cpusEdge, newService.standardMode.bandwidthEdge, NewCores)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (n *Node) IntraNodeReallocateTest(newService *Service, oldServiceID ServiceID) (bool, error) {
	NewCores := n.Cores
	oldServiceCores := n.AllocatedServices[oldServiceID].allocatedCoresEdge
	oldBandwidth := n.AllocatedServices[oldServiceID].standardMode.bandwidthEdge
	newBandwidth := newService.standardMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= oldBandwidth
	}

	possibleCores, err := n.NodeAdmission.Admission(newService.standardMode.cpusEdge, newService.standardMode.bandwidthEdge, NewCores)
	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		_, err = n.NodeAdmission.Admission(n.AllocatedServices[oldServiceID].standardMode.cpusEdge, oldBandwidth, NewCores)
		if err == nil {
			return true, nil
		}
	}
	return false, err
}

func (n *Node) NodeDeallocate(serviceID ServiceID) bool {
	cores := n.AllocatedServices[serviceID].allocatedCoresEdge
	mode := n.AllocatedServices[serviceID].allocationMode
	for _, core := range cores {
		switch mode {
		case StandardMode:
			n.Cores[core].ConsumedBandwidth -= n.AllocatedServices[serviceID].standardMode.bandwidthEdge
		case ReducedMode:
			n.Cores[core].ConsumedBandwidth -= n.AllocatedServices[serviceID].reducedMode.bandwidthEdge
		}
	}
	delete(n.AllocatedServices, serviceID)
	return true
}
