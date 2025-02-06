package orchestrator

import (
	"fmt"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

type Location string

const (
	cloudLoc Location = "cloud"
	edgeLoc  Location = "edge"
	bothLoc  Location = "both"
)

const (
	MinMin cnfg.Heuristic = "MinMin"
	MaxMax cnfg.Heuristic = "MaxMax"
)

type QoS int

type EventID string

type Orchestrator struct {
	Domains         Domains
	Cloud           *Cloud
	AllServices     Services
	RunningServices Services // change name of service to service
	Cost            cnfg.Cost
	QoS             QoS
	Config          *cnfg.OrchestratorConfig
}

func NewOrchestrator(config *cnfg.OrchestratorConfig, cloud *Cloud, domains Domains, services Services) *Orchestrator {
	domainCost := cnfg.Cost(0)
	for _, d := range domains {
		for _, n := range d.ActiveNodes {
			domainCost += cnfg.Cost(n.numCores) * config.EdgeNodeCost
		}
	}
	fmt.Println("Domain cost:", domainCost)
	cloudCost := cnfg.Cost(len(cloud.ActiveNodes)) * cnfg.Cost(config.CloudNodeCost)

	cost := domainCost + cloudCost
	o := &Orchestrator{

		Domains:         domains,
		Cloud:           cloud,
		AllServices:     services,
		Cost:            cost,
		QoS:             0,
		RunningServices: make(Services),
		Config:          config,
	}

	o.cloudPowerOnNode()
	// for domainID := range o.Domains {
	// 	o.edgePowerOnNode(domainID)
	// }
	return o

}

func (o *Orchestrator) allocateEdge(service *Service, node *Node, eventID ServiceID, cpuThreshold float64) (bool, error) {
	fmt.Println("Allocating standard service: ", service.serviceID, " to node: ", node.NodeName)
	allocated, svc, err := service.StandardMode.ServiceAllocate(service, node, eventID, cpuThreshold)
	if allocated {
		o.RunningServices[eventID] = svc
	}

	return allocated, err
}

func (o *Orchestrator) allocateEdgeReduced(service *Service, node *Node, eventID ServiceID) (bool, error) {
	const cpuThreshold = 100.0 //80.0

	allocated, svc, err := service.ReducedMode.EdgeServiceAllocate(service, node, eventID, cpuThreshold)
	if allocated {
		o.RunningServices[eventID] = svc
	}

	return allocated, err
}

func (o *Orchestrator) SplitSched(service *Service, domainID DomainID, eventID ServiceID) (bool, bool, *Service, error) {
	// edge-cloud split (has qos degradation) -- there is no cloud only apparently
	const cpuThreshold = 100.0
	const cloudCpuThreshold = 100.0
	fmt.Println("inside split scheduling")
	sortedNodes, _ := o.sortNodes(o.Domains[domainID].ActiveNodes, service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge)
	edgeAllocated := false
	cloudAllocated := false
	var svcEdge, svcCloud *Service

	for _, edgeNodeName := range sortedNodes {
		edgeAllocated, svcEdge, _ = service.ReducedMode.ServiceAllocate(service, o.Domains[domainID].ActiveNodes[edgeNodeName], edgeLoc, eventID, cpuThreshold)
		if edgeAllocated {
			break
		}
	}

	sortedNodes, _ = o.sortNodes(o.Cloud.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	for _, cloudNodeName := range sortedNodes {
		cloudAllocated, svcCloud, _ = service.ReducedMode.ServiceAllocate(service, o.Cloud.ActiveNodes[cloudNodeName], cloudLoc, eventID, cloudCpuThreshold)
		if cloudAllocated {
			break
		}
	}
	if !edgeAllocated || !cloudAllocated {
		return false, false, &Service{}, nil
	}
	fmt.Println("show svc in split scheduling", svcEdge)
	fmt.Println("show svc in split scheduling", svcCloud)
	newSvc := &Service{
		StandardMode:             svcEdge.StandardMode,
		ReducedMode:              svcEdge.ReducedMode,
		ImportanceFactor:         svcEdge.ImportanceFactor,
		serviceID:                svcEdge.serviceID,
		AllocatedCoresEdge:       svcEdge.AllocatedCoresEdge,
		AllocatedCoresCloud:      svcCloud.AllocatedCoresCloud,
		AllocatedNodeEdge:        svcEdge.AllocatedNodeEdge,
		AllocatedNodeCloud:       svcCloud.AllocatedNodeCloud,
		AllocatedDomain:          svcEdge.AllocatedDomain,
		AllocationMode:           ReducedMode,
		AverageResidualBandwidth: svcEdge.AverageResidualBandwidth,
		TotalResidualBandwidth:   svcEdge.TotalResidualBandwidth,
		StandardQoS:              svcEdge.StandardQoS,
		ReducedQoS:               svcCloud.ReducedQoS,
	}
	svcEdge = nil
	svcCloud = nil
	o.RunningServices[eventID] = newSvc
	return edgeAllocated, cloudAllocated, newSvc, nil
}

func (o *Orchestrator) Allocate(domainID DomainID, serviceID ServiceID, eventID ServiceID) (bool, error) {
	allocated := false
	domain := o.Domains[domainID]
	service := o.AllServices[serviceID]

	sortedNodes, _ := o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID, o.Config.DomainNodeThreshold)
		if allocated {
			o.QoS = o.QoS + service.StandardQoS
			return allocated, nil
		}
	}

	sortedNodesNoFilter, _ := o.sortNodesNoFilter(domain.ActiveNodes)

	intraDomainHelp := true
	for _, nodeName := range sortedNodesNoFilter {
		node := domain.ActiveNodes[nodeName]
		reallocatedEventID, _ := o.getReallocatedService(node, service)

		if node.AllocatedServices[reallocatedEventID] == nil {
			continue
		}
		reallocationHelp, NewCores, _ := ReallocateTest(service, reallocatedEventID, *node)
		if reallocationHelp {
			ctx := ReallocContext{
				Service:            service,
				Node:               node,
				Domain:             domain,
				SortedNodes:        sortedNodes,
				EventID:            eventID,
				ReallocatedEventID: reallocatedEventID,
				NewCores:           NewCores,
			}
			if o.Config.IntraNodeRealloc {
				allocated, err := o.intraNodeRealloc(ctx)
				if allocated {
					o.QoS = o.QoS + service.StandardQoS
					return allocated, nil
				}
				if err != nil {
					fmt.Println("Error in intra node reallocation: ", err)
				}
			}
			if o.Config.IntraDomainRealloc {
				fmt.Println("going to intra domain reallocation")
				if intraDomainHelp {
					allocated, err := o.intraDomainRealloc(ctx)
					if allocated {
						o.QoS = o.QoS + service.StandardQoS
						fmt.Println("allocated with intra domain reallocation")
						return allocated, nil
					}
					if err != nil {
						fmt.Println("Error in intra domain reallocation:", err)
					}
				}
			}
			// intraDomainHelp = false
			if o.Config.IntraNodeReduced {
				fmt.Println("going to intra node cloud reallocation")
				allocated, _ = o.IntraNodeCloudRealloc(ctx)
				if allocated {
					otherEvent := o.RunningServices[reallocatedEventID]
					o.QoS = o.QoS + service.StandardQoS - otherEvent.StandardQoS + otherEvent.ReducedQoS
					fmt.Println("allocated with intra node cloud reallocation")
					return allocated, nil
				}
			}
			if o.Config.IntraNodeRemoved {
				fmt.Println("going to intra node cloud reallocation")
				otherEvent := o.RunningServices[reallocatedEventID]
				allocated, _ = o.intraNodeRemove(ctx)
				if allocated {
					o.QoS = o.QoS + service.StandardQoS - otherEvent.StandardQoS
					fmt.Println("allocated with intra node cloud reallocation")
					return allocated, nil
				}
			}
		} else {
			fmt.Println("Reallocated event not helpful")
		}

	}

	// reduced mode edge only

	// sortedNodes, _ = o.sortNodes(domain.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	// for _, nodeName := range sortedNodes {
	// 	allocated, _ := o.allocateEdgeReduced(service, domain.ActiveNodes[nodeName], eventID)
	// 	if allocated {
	// 		fmt.Println("qos before allocation", o.QoS)
	// 		fmt.Println("qos of the service", service.ReducedQoS)
	// 		fmt.Println("event id for edge reduced allocation:", eventID)
	// 		o.QoS = o.QoS + service.EdgeReducedQoS
	// 		return allocated, nil
	// 	}
	// }
	sortedNodes, _ = o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID, 100)
		if allocated {
			o.QoS = o.QoS + service.StandardQoS
			return allocated, nil
		}
	}

	edgeAllocated, cloudAllocated, svc, _ := o.SplitSched(service, domainID, eventID)
	if edgeAllocated && cloudAllocated {
		fmt.Println("show svc in split scheduling", svc)
		o.QoS = o.QoS + service.ReducedQoS
		return true, nil
	} else {
		fmt.Println("the split scheduling didn't work. powering on some nodes")
		if !edgeAllocated {
			success, nodeName := o.edgePowerOnNode(domainID)
			fmt.Println("Edge node powered on, node name: ", nodeName)
			if success {
				allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID, o.Config.DomainNodeThreshold)
				if allocated {
					o.QoS = o.QoS + service.StandardQoS
					return allocated, nil
				}
			}
		}

		if !cloudAllocated {
			o.cloudPowerOnNode()
		}
	}
	edgeAllocated, cloudAllocated, svc, _ = o.SplitSched(service, domainID, eventID)
	if edgeAllocated && cloudAllocated {
		allocated = true
		o.QoS = o.QoS + service.ReducedQoS
	} else {
		return false, nil
	}

	return allocated, nil
}

func (o *Orchestrator) Deallocate(domainID DomainID, serviceID ServiceID, eventID ServiceID) bool {
	domain := o.Domains[domainID]
	service := o.RunningServices[eventID]
	allocatedMode := o.RunningServices[eventID].AllocationMode
	var serviceQoS QoS

	if allocatedMode == StandardMode {
		serviceQoS = service.StandardQoS
		nodeN := service.AllocatedNodeEdge
		node := domain.ActiveNodes[nodeN]
		_, err := service.StandardMode.ServiceDeallocate(eventID, node)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - serviceQoS
		}

	}
	if allocatedMode == ReducedMode {
		serviceQoS = service.ReducedQoS
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		cloudNode := o.Cloud.ActiveNodes[service.AllocatedNodeCloud]
		_, err := service.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
		_, err = service.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - service.ReducedQoS
		}

	}
	if allocatedMode == EdgeReducedMode {
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		_, err := service.ReducedMode.EdgeServiceDeallocate(eventID, edgeNode)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - service.EdgeReducedQoS
		}

	}
	for nodeName, node := range domain.ActiveNodes {
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.edgePowerOffNode(domainID, nodeName)
		}
	}
	for nodeName, node := range o.Cloud.ActiveNodes {
		if len(o.Cloud.ActiveNodes) == 1 {
			break
		}
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.cloudPowerOffNode(nodeName)
		}
	}

	delete(o.RunningServices, eventID)

	return true
}
