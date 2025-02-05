package orchestrator

import (
	"fmt"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

type NodeName string
type AllocatedServices map[ServiceID]*Service
type NodeStatus string

const (
	Active   NodeStatus = "Active"
	Inactive NodeStatus = "Inactive"
)

type Node struct {
	Cores                    Cores
	numCores                 uint64
	ReallocHeuristic         cnfg.Heuristic
	NodeName                 NodeName
	NodeAdmission            *AdmissionTest
	Location                 Location
	DomainID                 DomainID
	AllocatedServices        AllocatedServices
	AverageResidualBandwidth float64
	TotalResidualBandwidth   float64
	Status                   NodeStatus
}

func NewNode(cores Cores, heuristic cnfg.Heuristic, nodeName NodeName, domainID DomainID) *Node {
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
		DomainID:                 domainID,
		numCores:                 uint64(len(cores)),
	}
}

func (n *Node) NodeAllocate(reqCpus uint64, reqBandwidth float64, service *Service, eventID ServiceID, cpuThreshold float64) ([]CoreID, error) {
	selectedCpus, err := n.NodeAdmission.Admission(reqCpus, reqBandwidth, n.Cores, cpuThreshold)
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

func ReallocateTest(newService *Service, oldServiceID ServiceID, n Node) (bool, Cores, error) {
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	fmt.Println("Reallocate Test")
	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.bandwidthEdge
	newBandwidth := newService.StandardMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth

	}
	possibleCores, err := n.NodeAdmission.Admission(newService.StandardMode.cpusEdge, newService.StandardMode.bandwidthEdge, NewCores, 100.0)

	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		return true, NewCores, nil
	}
	return false, nil, err
}

func IntraDomainReallocateTest(newService *Service, oldServiceID ServiceID, n *Node) (bool, error) {
	fmt.Println("Intra Domain Reallocate Test")
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	// n.Cores = NewCores
	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth
	}
	_, err := n.NodeAdmission.Admission(newService.StandardMode.cpusEdge, newService.StandardMode.bandwidthEdge, NewCores, 100.0)
	for _, core := range NewCores {
		fmt.Println("cores after intra domain realocate test newcores:", core)
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

// func IntraNodeReallocateTest(newService *Service, oldEventID ServiceID, n Node, NewCores Cores) (bool, error) {
// 	// NewCores := make(Cores, len(n.Cores))
// 	// for k, v := range n.Cores {
// 	// 	NewCores[k] = &Core{
// 	// 		ConsumedBandwidth: v.ConsumedBandwidth,
// 	// 		ID:                v.ID,
// 	// 	}
// 	// }
// 	// // n.Cores = NewCores

// 	// oldServiceCores := n.AllocatedServices[oldEventID].AllocatedCoresEdge
// 	oldBandwidth := n.AllocatedServices[oldEventID].StandardMode.bandwidthEdge
// 	// for _, coreID := range oldServiceCores {
// 	// 	NewCores[coreID].ConsumedBandwidth -= oldBandwidth
// 	// }

// 	// possibleCores, err := n.NodeAdmission.Admission(newService.StandardMode.cpusEdge, newService.StandardMode.bandwidthEdge, NewCores)
// 	_, err := n.NodeAdmission.Admission(n.AllocatedServices[oldEventID].StandardMode.cpusEdge, oldBandwidth, NewCores)
// 	if err == nil {
// 		return true, nil
// 	}
// 	return false, err
// }

func (n *Node) NodeDeallocate(eventID ServiceID) bool {
	fmt.Println("going to Node Deallocate", n)
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
	case EdgeReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthCloud

	}
	totalResidualBandwidth := n.TotalResidualBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		totalResidualBandwidth -= bandwidth
	}
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))
	if n.TotalResidualBandwidth < 0 {
		return false
	}
	delete(n.AllocatedServices, eventID)
	return true
}

func (n *Node) CloudNodeDeallocate(eventID ServiceID) bool {
	fmt.Println("going to Node Deallocate", n)
	fmt.Println("allocated services:", n.AllocatedServices)
	fmt.Println("eventID:", eventID)
	cores := n.AllocatedServices[eventID].AllocatedCoresCloud
	fmt.Println("cores for deallocating in node deallocate:", cores)

	bandwidth := n.AllocatedServices[eventID].ReducedMode.bandwidthCloud

	totalResidualBandwidth := n.TotalResidualBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		totalResidualBandwidth -= bandwidth
	}
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))
	if n.TotalResidualBandwidth < 0 {
		return false
	}
	delete(n.AllocatedServices, eventID)
	return true
}

func (n *Node) Upgraded(event *Service) {
	cores := event.AllocatedCoresEdge
	bandwidth := event.ReducedMode.bandwidthEdge

	totalResidualBandwidth := n.TotalResidualBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		totalResidualBandwidth -= bandwidth
	}
	n.TotalResidualBandwidth = totalResidualBandwidth
	n.AverageResidualBandwidth = totalResidualBandwidth / float64(len(n.Cores))
}
