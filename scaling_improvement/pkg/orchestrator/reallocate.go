package orchestrator

import (
	"fmt"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

const (
	HBI   cnfg.Heuristic = "HBI"
	HCI   cnfg.Heuristic = "HCI"
	HBCI  cnfg.Heuristic = "HBCI"
	HB    cnfg.Heuristic = "HB"
	HC    cnfg.Heuristic = "HC"
	HBC   cnfg.Heuristic = "HBC"
	LB    cnfg.Heuristic = "LB"
	LC    cnfg.Heuristic = "LC"
	LBC   cnfg.Heuristic = "LBC"
	HCLI  cnfg.Heuristic = "HCLI"
	HBLI  cnfg.Heuristic = "HBLI"
	HBIcC cnfg.Heuristic = "HBIcC"
	LBCI  cnfg.Heuristic = "LBCI"
	LBI   cnfg.Heuristic = "LBI"
	LCI   cnfg.Heuristic = "LCI"
	LRED  cnfg.Heuristic = "LRED"
	LREM  cnfg.Heuristic = "LREM"
	LI    cnfg.Heuristic = "LI"
)

type ReallocContext struct {
	Service            *Service
	Node               *Node
	Domain             *Domain
	SortedNodes        []NodeName
	EventID            ServiceID
	ReallocatedEventID ServiceID
	NewCores           Cores
}

func (o *Orchestrator) getReallocatedService(node *Node, t *Service) (ServiceID, ServiceID, ServiceID, error) {
	var selectedEventID ServiceID
	var bestScore float64

	calculateScore := func(service *Service, heuristic cnfg.Heuristic) float64 {
		switch heuristic {
		case LI:
			return 1 / service.ImportanceFactor
		case HB:
			return (service.StandardMode.bandwidthEdge)
		case HC:
			return float64(service.StandardMode.cpusEdge)
		case HBC:
			return (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
		case LB:
			return (1 / service.StandardMode.bandwidthEdge)
		case LC:
			return 1 / float64(service.StandardMode.cpusEdge)
		case LBC:
			return 1 / (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
		case LBI:
			return 1 / (service.ImportanceFactor * (service.StandardMode.bandwidthEdge))
		case LCI:
			return 1 / (service.ImportanceFactor * (float64(service.StandardMode.cpusEdge)))
		case HBI:
			return (service.ImportanceFactor * service.StandardMode.bandwidthEdge)
		case HCI:
			return service.ImportanceFactor * float64(service.StandardMode.cpusEdge)
		case HBCI:
			return service.ImportanceFactor * (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
		case HBLI:
			return (service.StandardMode.bandwidthEdge * float64(1/service.ImportanceFactor))
		case HCLI:
			return float64(service.StandardMode.cpusEdge) * float64(1/service.ImportanceFactor)
		case LBCI:
			return 1 / (service.ImportanceFactor * (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge)))
		case LRED:
			if (service.StandardQoS - service.ReducedQoS) < (t.StandardQoS - t.ReducedQoS) {
				return 1 / float64(service.StandardQoS)
			}
		case LREM:
			if service.StandardQoS < (t.StandardQoS - t.ReducedQoS) {
				return 1 / float64(service.StandardQoS)
			}
		case HBIcC:
			if service.StandardMode.cpusEdge >= t.StandardMode.cpusEdge {
				return service.StandardMode.bandwidthEdge
			}
			// if service.StandardMode.bandwidthEdge >= t.StandardMode.bandwidthEdge {
			// 	return 1 / float64(service.StandardMode.cpusEdge)
			// }
		}
		return 0
	}

	// for eventID, service := range node.AllocatedServices {
	// 	if service.AllocationMode == StandardMode {
	// 		score := calculateScore(service, o.Config.ReallocationHeuristic)
	// 		if score > bestScore {
	// 			bestScore = score
	// 			selectedEventID = eventID
	// 		}
	// 	}
	// }
	var secondBestScore, thirdBestScore float64
	var secondSelectedEventID, thirdSelectedEventID ServiceID

	for eventID, service := range node.AllocatedServices {
		if service.AllocationMode == StandardMode {
			score := calculateScore(service, o.Config.ReallocationHeuristic)
			if score > bestScore {
				thirdBestScore = secondBestScore
				thirdSelectedEventID = secondSelectedEventID
				secondBestScore = bestScore
				secondSelectedEventID = selectedEventID
				bestScore = score
				selectedEventID = eventID
			} else if score > secondBestScore {
				thirdBestScore = secondBestScore
				thirdSelectedEventID = secondSelectedEventID
				secondBestScore = score
				secondSelectedEventID = eventID
			} else if score > thirdBestScore {
				thirdBestScore = score
				thirdSelectedEventID = eventID
			}
		}
	}

	if selectedEventID == "" && secondSelectedEventID == "" && thirdSelectedEventID == "" {
		return "", "", "", fmt.Errorf("no suitable services found for reallocation using heuristic %s", o.Config.ReallocationHeuristic)
	}

	return selectedEventID, secondSelectedEventID, thirdSelectedEventID, nil

}

func (o *Orchestrator) intraNodeRealloc(ctx ReallocContext) (bool, error) {
	const cpuThreshold = 100.0
	if ctx.ReallocatedEventID == "" {
		return false, nil
	}
	oldBandwidth := o.RunningServices[ctx.ReallocatedEventID].StandardMode.bandwidthEdge
	_, err := ctx.Node.NodeAdmission.Admission(ctx.Node.AllocatedServices[ctx.ReallocatedEventID].StandardMode.cpusEdge, oldBandwidth, ctx.NewCores, cpuThreshold)

	if err != nil {
		return false, err
	}
	fmt.Println("show the status of the intra node reallocation:", true)

	reallocatedService := ctx.Node.AllocatedServices[ctx.ReallocatedEventID]
	_, err = ctx.Service.StandardMode.ServiceDeallocate(ctx.ReallocatedEventID, ctx.Node)
	if err != nil {
		return false, err
	}
	_, newSvc, _ := ctx.Service.StandardMode.ServiceAllocate(ctx.Service, ctx.Node, ctx.EventID, cpuThreshold)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", ctx.Node.AverageResidualBandwidth)
	_, oldSvc, _ := reallocatedService.StandardMode.ServiceAllocate(reallocatedService, ctx.Node, ctx.ReallocatedEventID, cpuThreshold)
	fmt.Println("Allocated services in the end: ", ctx.Node.AllocatedServices)
	// delete(o.RunningServices, reallocatedEventID)
	o.RunningServices[ctx.ReallocatedEventID] = oldSvc
	o.RunningServices[ctx.EventID] = newSvc
	fmt.Println("intra node reallocation completed")

	return true, nil
}

func (o *Orchestrator) intraDomainRealloc(ctx ReallocContext) (bool, error) {
	const cpuThreshold = 100.0
	if ctx.ReallocatedEventID == "" {
		return false, nil
	}

	reallocated := false

	otherService := ctx.Node.AllocatedServices[ctx.ReallocatedEventID]
	fmt.Println("inside the intra domain reallocation")
	for _, nodeName := range ctx.SortedNodes {
		if nodeName == ctx.Node.NodeName {
			continue
		}
		otherNode := ctx.Domain.ActiveNodes[nodeName]
		for _, core := range otherNode.Cores {
			fmt.Println("cores of the other node:", core)
		}

		allocatedCore, _ := otherNode.NodeAdmission.Admission(otherService.StandardMode.cpusEdge, otherService.StandardMode.bandwidthEdge, otherNode.Cores, cpuThreshold)

		fmt.Println("allocated core:", allocatedCore)
		if allocatedCore != nil {
			fmt.Println("reallocation successful")

			otherService.StandardMode.ServiceDeallocate(ctx.ReallocatedEventID, ctx.Node)
			_, newSvc, _ := ctx.Service.StandardMode.ServiceAllocate(ctx.Service, ctx.Node, ctx.EventID, cpuThreshold)
			_, oldSvc, _ := otherService.StandardMode.ServiceAllocate(otherService, otherNode, ctx.ReallocatedEventID, cpuThreshold)
			o.RunningServices[ctx.ReallocatedEventID] = oldSvc
			o.RunningServices[ctx.EventID] = newSvc
			reallocated = true
			fmt.Println("intra domain reallocation completed")
			return reallocated, nil
		}
	}
	return reallocated, fmt.Errorf("intra domain reallocation failed")
}

func (o *Orchestrator) IntraNodeCloudRealloc(ctx ReallocContext) (bool, error) {
	const cpuThreshold = 100.0
	if ctx.ReallocatedEventID == "" {
		return false, nil
	}
	otherService := ctx.Node.AllocatedServices[ctx.ReallocatedEventID]
	oldBandwidth := o.RunningServices[ctx.ReallocatedEventID].ReducedMode.bandwidthCloud
	oldCpus := o.RunningServices[ctx.ReallocatedEventID].ReducedMode.cpusCloud
	selectedCPUs, err := ctx.Node.NodeAdmission.Admission(oldCpus, oldBandwidth, ctx.NewCores, cpuThreshold)

	if err != nil || selectedCPUs == nil {
		return false, err
	}
	// cloudAllocated := false
	var svcEdge, svcCloud *Service
	var cloudNode NodeName
	sortedNodes, _ := o.sortNodes(o.Cloud.ActiveNodes, otherService.ReducedMode.cpusCloud, otherService.ReducedMode.bandwidthCloud)
	for _, cloudNodeName := range sortedNodes {
		selectedCPUs, err := ctx.Node.NodeAdmission.Admission(otherService.ReducedMode.cpusCloud, otherService.ReducedMode.bandwidthCloud, o.Cloud.ActiveNodes[cloudNodeName].Cores, cpuThreshold)
		if err == nil && selectedCPUs != nil {
			cloudNode = cloudNodeName
			break
		}
		//
		// if cloudAllocated {
		// 	break
		// }
	}
	if cloudNode == "" {
		return false, fmt.Errorf("cloud not suitable for reallocation")
	}

	fmt.Println("show the status of the intra node cloud reallocation:", true)

	_, err = otherService.StandardMode.ServiceDeallocate(ctx.ReallocatedEventID, ctx.Node)
	if err != nil {
		return false, err
	}
	allocated, newSvc, _ := ctx.Service.StandardMode.ServiceAllocate(ctx.Service, ctx.Node, ctx.EventID, cpuThreshold)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", ctx.Node.AverageResidualBandwidth)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation")
	}
	allocated, svcEdge, _ = otherService.ReducedMode.ServiceAllocate(otherService, ctx.Node, edgeLoc, ctx.ReallocatedEventID, cpuThreshold)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation in edge")
	}
	allocated, svcCloud, _ = otherService.ReducedMode.ServiceAllocate(otherService, o.Cloud.ActiveNodes[cloudNode], cloudLoc, ctx.ReallocatedEventID, cpuThreshold)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation in cloud")
	}

	fmt.Println("Allocated services in the end: ", ctx.Node.AllocatedServices)
	// delete(o.RunningServices, reallocatedEventID)
	oldSvc := &Service{
		StandardMode:             svcEdge.StandardMode,
		ReducedMode:              svcEdge.ReducedMode,
		ImportanceFactor:         svcEdge.ImportanceFactor,
		serviceID:                otherService.serviceID,
		AllocatedCoresEdge:       svcEdge.AllocatedCoresEdge,
		AllocatedCoresCloud:      svcCloud.AllocatedCoresCloud,
		AllocatedNodeEdge:        svcEdge.AllocatedNodeEdge,
		AllocatedNodeCloud:       svcCloud.AllocatedNodeCloud,
		AllocatedDomain:          svcEdge.AllocatedDomain,
		AllocationMode:           ReducedMode,
		AverageResidualBandwidth: svcEdge.AverageResidualBandwidth,
		TotalResidualBandwidth:   svcEdge.TotalResidualBandwidth,
		StandardQoS:              otherService.StandardQoS,
		ReducedQoS:               otherService.ReducedQoS,
	}
	ctx.Node.AllocatedServices[ctx.ReallocatedEventID] = oldSvc
	svcEdge = nil
	svcCloud = nil
	o.RunningServices[ctx.ReallocatedEventID] = oldSvc
	o.RunningServices[ctx.EventID] = newSvc
	fmt.Println("intra node cloud reallocation completed")
	return true, nil
}

func (o *Orchestrator) intraNodeRemove(ctx ReallocContext) (bool, error) {
	const cpuThreshold = 100.0
	if ctx.ReallocatedEventID == "" {
		return false, nil
	}
	oldBandwidth := o.RunningServices[ctx.ReallocatedEventID].StandardMode.bandwidthEdge
	_, err := ctx.Node.NodeAdmission.Admission(ctx.Node.AllocatedServices[ctx.ReallocatedEventID].StandardMode.cpusEdge, oldBandwidth, ctx.NewCores, cpuThreshold)

	if err != nil {
		return false, err
	}
	fmt.Println("show the status of the intra node reallocation:", true)

	_, err = ctx.Service.StandardMode.ServiceDeallocate(ctx.ReallocatedEventID, ctx.Node)
	if err != nil {
		return false, err
	}
	delete(o.RunningServices, ctx.ReallocatedEventID)
	_, newSvc, _ := ctx.Service.StandardMode.ServiceAllocate(ctx.Service, ctx.Node, ctx.EventID, cpuThreshold)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", ctx.Node.AverageResidualBandwidth)
	o.RunningServices[ctx.EventID] = newSvc
	fmt.Println("intra node reallocation completed")

	return true, nil
}

type ReallocHandler func(ReallocContext) (bool, error)

func (o *Orchestrator) ReallocStrategies() []ReallocHandler {
	var strategies []ReallocHandler

	if o.Config.IntraNodeRealloc {
		strategies = append(strategies, o.intraNodeRealloc)
	}
	if o.Config.IntraDomainRealloc {
		strategies = append(strategies, o.intraDomainRealloc)
	}
	if o.Config.IntraNodeReduced {
		strategies = append(strategies, o.IntraNodeCloudRealloc)
	}
	return strategies
}
