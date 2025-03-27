package orchestrator

import (
	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

type NodeName string
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
	AllocatedServices        Services
	AverageConsumedBandwidth float64
	TotalConsumedBandwidth   float64
	Status                   NodeStatus
}

func NewNode(cores Cores, heuristic cnfg.Heuristic, nodeName NodeName, domainID DomainID) *Node {
	admissionTest := NewAdmissionTest(cores, heuristic)
	// log.Info("Node Admission Test Cores: ", admissionTest.Cores[CoreID("core-0")])
	return &Node{
		Cores:                    cores,
		ReallocHeuristic:         heuristic,
		NodeName:                 nodeName,
		NodeAdmission:            admissionTest,
		Status:                   Inactive,
		AverageConsumedBandwidth: 0,
		TotalConsumedBandwidth:   0,
		AllocatedServices:        make(Services),
		DomainID:                 domainID,
		numCores:                 uint64(len(cores)),
	}
}

func (n *Node) NodeAllocate(reqCpus uint64, reqBandwidth float64, service *Service, eventID ServiceID, cpuThreshold float64) ([]CoreID, error) {
	selectedCpus, err := n.NodeAdmission.Admission(reqCpus, reqBandwidth, n.Cores, cpuThreshold)
	if err != nil {
		return selectedCpus, err
	}
	TotalConsumedBandwidth := n.TotalConsumedBandwidth

	for _, coreID := range selectedCpus {
		n.Cores[coreID].ConsumedBandwidth += reqBandwidth
		TotalConsumedBandwidth += reqBandwidth
	}

	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))

	return selectedCpus, nil
}

func ReallocateTest(newService *Service, oldServiceID ServiceID, n Node) (bool, Cores, error) {
	if oldServiceID == "" {
		return false, nil, nil
	}
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}

	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.BandwidthEdge
	newBandwidth := newService.StandardMode.BandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth

	}
	possibleCores, err := n.NodeAdmission.Admission(newService.StandardMode.CpusEdge, newService.StandardMode.BandwidthEdge, NewCores, 100.0)

	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		return true, NewCores, nil
	}
	return false, nil, err
}

func ReallocateTestReduced(newService *Service, oldServiceID ServiceID, n Node) (bool, Cores, error) {
	if oldServiceID == "" {
		return false, nil, nil
	}
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.BandwidthEdge
	newBandwidth := newService.ReducedMode.bandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth

	}
	possibleCores, err := n.NodeAdmission.Admission(newService.ReducedMode.cpusEdge, newService.ReducedMode.bandwidthEdge, NewCores, 100.0)

	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		return true, NewCores, nil
	}
	return false, nil, err
}

func IntraDomainReallocateTest(newService *Service, oldServiceID ServiceID, n *Node) (bool, error) {
	NewCores := make(Cores, len(n.Cores))
	for k, v := range n.Cores {
		NewCores[k] = &Core{
			ConsumedBandwidth: v.ConsumedBandwidth,
			ID:                v.ID,
		}
	}
	// n.Cores = NewCores
	oldServiceCores := n.AllocatedServices[oldServiceID].AllocatedCoresEdge
	bandwidth := n.AllocatedServices[oldServiceID].StandardMode.BandwidthEdge
	for _, coreID := range oldServiceCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth
	}
	_, err := n.NodeAdmission.Admission(newService.StandardMode.CpusEdge, newService.StandardMode.BandwidthEdge, NewCores, 100.0)

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
	// log.Info("allocated services:", n.AllocatedServices)
	cores := n.AllocatedServices[eventID].AllocatedCoresEdge
	mode := n.AllocatedServices[eventID].AllocationMode

	bandwidth := n.AllocatedServices[eventID].StandardMode.BandwidthEdge
	switch mode {
	case StandardMode:
		bandwidth = n.AllocatedServices[eventID].StandardMode.BandwidthEdge
	case ReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthEdge
	case EdgeReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthCloud

	}

	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		TotalConsumedBandwidth -= bandwidth
	}

	delete(n.AllocatedServices, eventID)
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))
	if n.TotalConsumedBandwidth < 0 {
		return false
	}

	return true
}

func (n *Node) CloudNodeDeallocate(eventID ServiceID) bool {

	cores := n.AllocatedServices[eventID].AllocatedCoresCloud

	bandwidth := n.AllocatedServices[eventID].ReducedMode.bandwidthCloud
	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		TotalConsumedBandwidth -= bandwidth
	}
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	delete(n.AllocatedServices, eventID)
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))
	if n.TotalConsumedBandwidth < 0 {
		return false
	}

	return true
}

func (n *Node) Upgraded(event *Service) {
	cores := event.AllocatedCoresEdge
	bandwidth := event.ReducedMode.bandwidthEdge

	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		TotalConsumedBandwidth -= bandwidth
	}
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))
}
