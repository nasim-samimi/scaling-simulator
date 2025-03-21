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
	log.Info("Selected CPUs: ", selectedCpus)
	if err != nil {
		return selectedCpus, err
	}
	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	log.Info("NodeAllocate, before allocation")
	for _, core := range n.Cores {
		log.Info("core: ", core.ConsumedBandwidth)
	}
	for _, coreID := range selectedCpus {
		n.Cores[coreID].ConsumedBandwidth += reqBandwidth
		TotalConsumedBandwidth += reqBandwidth
	}
	log.Info("NodeAllocate, after allocation")
	for _, core := range n.Cores {
		log.Info("core: ", core.ConsumedBandwidth)
	}

	log.Info("Service: ", service)
	// service.AllocatedNodeEdge = n.NodeName
	log.Info("in node allocation, total residual bandwidth and average: ", n.TotalConsumedBandwidth, n.AverageConsumedBandwidth)
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))
	log.Info("in node allocation, total residual bandwidth and average after update: ", n.TotalConsumedBandwidth, n.AverageConsumedBandwidth)
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
	log.Info("Reallocate Test")
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

func IntraDomainReallocateTest(newService *Service, oldServiceID ServiceID, n *Node) (bool, error) {
	log.Info("Intra Domain Reallocate Test")
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
	for _, core := range NewCores {
		log.Info("cores after intra domain realocate test newcores:", core)
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
	log.Info("going to Node Deallocate", n)
	log.Info("allocated services:", n.AllocatedServices)
	log.Info("eventID:", eventID)
	cores := n.AllocatedServices[eventID].AllocatedCoresEdge
	log.Info("cores for deallocating in node deallocate:", cores)
	mode := n.AllocatedServices[eventID].AllocationMode
	log.Info("deallocation mode: ", mode)

	bandwidth := n.AllocatedServices[eventID].StandardMode.BandwidthEdge
	switch mode {
	case StandardMode:
		bandwidth = n.AllocatedServices[eventID].StandardMode.BandwidthEdge
	case ReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthEdge
	case EdgeReducedMode:
		bandwidth = n.AllocatedServices[eventID].ReducedMode.bandwidthCloud

	}
	log.Info("bandwidth of the service:", bandwidth)
	log.Info("total bandwidth before edge deallocation: ", n.TotalConsumedBandwidth)
	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		TotalConsumedBandwidth -= bandwidth
	}
	log.Info("total bandwidth after edge deallocation: ", TotalConsumedBandwidth)
	delete(n.AllocatedServices, eventID)
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	n.AverageConsumedBandwidth = TotalConsumedBandwidth / float64(len(n.Cores))
	if n.TotalConsumedBandwidth < 0 {
		return false
	}

	return true
}

func (n *Node) CloudNodeDeallocate(eventID ServiceID) bool {
	log.Info("going to Node Deallocate", n)
	log.Info("allocated services:", n.AllocatedServices)
	log.Info("eventID:", eventID)
	cores := n.AllocatedServices[eventID].AllocatedCoresCloud
	log.Info("cores for deallocating in cloud node deallocate:", cores)

	bandwidth := n.AllocatedServices[eventID].ReducedMode.bandwidthCloud
	log.Info("bandwidth of the service:", bandwidth)
	log.Info("total bandwidth before cloud deallocation: ", n.TotalConsumedBandwidth)
	TotalConsumedBandwidth := n.TotalConsumedBandwidth
	for _, core := range cores {
		n.Cores[core].ConsumedBandwidth -= bandwidth
		TotalConsumedBandwidth -= bandwidth
	}
	n.TotalConsumedBandwidth = TotalConsumedBandwidth
	delete(n.AllocatedServices, eventID)
	log.Info("total bandwidth after cloud deallocation: ", TotalConsumedBandwidth)
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
