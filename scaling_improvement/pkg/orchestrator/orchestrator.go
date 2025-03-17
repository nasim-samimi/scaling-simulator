package orchestrator

import (
	"os"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Automatically runs when the package is imported
func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)       // Log to console
	log.SetLevel(logrus.InfoLevel) // Set log level
}

type Location string

const (
	cloudLoc Location = "cloud"
	edgeLoc  Location = "edge"
	bothLoc  Location = "both"
)

const (
	Min  cnfg.Heuristic = "Min"
	Max  cnfg.Heuristic = "Max"
	MmRB cnfg.Heuristic = "MmRB"
	mMRB cnfg.Heuristic = "mMRB"
	MMRB cnfg.Heuristic = "MMRB"
	mmRB cnfg.Heuristic = "mmRB"
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
	log.Info("Domain cost:", domainCost)
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

func (o *Orchestrator) decreaseQoS(qos QoS) {
	o.QoS = o.QoS - qos
	if o.QoS < 0 {
		log.Error("QoS is negative")
	}
}

func (o *Orchestrator) allocateEdge(service *Service, node *Node, eventID ServiceID, cpuThreshold float64) (bool, error) {
	log.Info("Allocating standard service: ", service.serviceID, " to node: ", node.NodeName)
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
	log.Info("inside split scheduling")
	log.Info("reduced mode of the service", service.ReducedMode)
	sortedNodes, _ := o.sortNodes(o.Domains[domainID].ActiveNodes, service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge)
	edgeAllocated := false
	cloudAllocated := false
	var svcEdge, svcCloud *Service
	var edgeNodeName NodeName
	var cloudNodeName NodeName

	for _, eNodeName := range sortedNodes {
		potentialNode := o.Domains[domainID].ActiveNodes[eNodeName]
		potentialCores, _ := potentialNode.NodeAdmission.Admission(service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge, potentialNode.Cores, cpuThreshold)
		if potentialCores == nil {
			continue
		}
		edgeNodeName = eNodeName
		break
	}
	if !edgeAllocated {
		return false, false, &Service{}, nil
	}

	sortedNodes, _ = o.sortNodes(o.Cloud.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	for _, cNodeName := range sortedNodes {
		potentialNode := o.Cloud.ActiveNodes[cNodeName]
		potentialCores, _ := potentialNode.NodeAdmission.Admission(service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud, potentialNode.Cores, cloudCpuThreshold)
		if potentialCores == nil {
			continue
		}
		cloudNodeName = cNodeName
		break

	}
	if edgeNodeName == "" || cloudNodeName == "" {
		return false, false, &Service{}, nil
	}
	// real allocation
	edgeAllocated, svcEdge, _ = service.ReducedMode.ServiceAllocate(service, o.Domains[domainID].ActiveNodes[edgeNodeName], edgeLoc, eventID, cpuThreshold)
	cloudAllocated, svcCloud, _ = service.ReducedMode.ServiceAllocate(service, o.Cloud.ActiveNodes[cloudNodeName], cloudLoc, eventID, cloudCpuThreshold)

	if !edgeAllocated || !cloudAllocated {
		return false, false, &Service{}, nil
	}
	log.Info("show svc in split scheduling", svcEdge)
	log.Info("show svc in split scheduling", svcCloud)
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
		AverageConsumedBandwidth: svcEdge.AverageConsumedBandwidth,
		TotalConsumedBandwidth:   svcEdge.TotalConsumedBandwidth,
		StandardQoS:              svcEdge.StandardQoS,
		ReducedQoS:               svcCloud.ReducedQoS,
	}
	svcEdge = nil
	svcCloud = nil
	o.RunningServices[eventID] = newSvc
	for _, n := range o.Cloud.ActiveNodes {
		log.Info("average consumed bandwidth in cloud nodes, after allocating", n.AverageConsumedBandwidth)
	}
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

	sortedNodesNoFilter, _ := o.sortNodesNoFilter(domain.ActiveNodes, Max)

	intraDomainHelp := true
	// reallocationHelp := false
	// reducedHelpful := false
	// removedHelpful := false
	for _, nodeName := range sortedNodesNoFilter {
		node := domain.ActiveNodes[nodeName]
		reallocatedEventID, secondReallocatedEventID, thirdreallocatedEventID, err := o.getReallocatedService(node, service, o.Config.IntraNodeReallocHeu)
		if err != nil {
			continue
		}

		reallocationHelp, NewCores, _ := ReallocateTest(service, reallocatedEventID, *node)
		selectedEventID := reallocatedEventID
		// if service.StandardQoS-o.RunningServices[selectedEventID].StandardQoS+o.RunningServices[selectedEventID].ReducedQoS > service.ReducedQoS {
		// 	reducedHelpful = true
		// }
		// if service.StandardQoS-o.RunningServices[selectedEventID].StandardQoS > service.ReducedQoS {
		// 	removedHelpful = true
		// }
		if !reallocationHelp {
			reallocationHelp, NewCores, _ = ReallocateTest(service, secondReallocatedEventID, *node)
			selectedEventID = secondReallocatedEventID
		}
		if !reallocationHelp {
			reallocationHelp, NewCores, _ = ReallocateTest(service, thirdreallocatedEventID, *node)
			selectedEventID = thirdreallocatedEventID
		}
		log.Info("selected event id for reallocation:", selectedEventID)

		if reallocationHelp {
			ctx := ReallocContext{
				Service:            service,
				Node:               node,
				Domain:             domain,
				SortedNodes:        sortedNodes,
				EventID:            eventID,
				ReallocatedEventID: selectedEventID,
				NewCores:           NewCores,
			}
			if o.Config.IntraNodeRealloc {
				allocated, err := o.intraNodeRealloc(ctx)
				if allocated {
					o.QoS = o.QoS + service.StandardQoS
					return allocated, nil
				}
				if err != nil {
					log.Info("Error in intra node reallocation: ", err)
				}
			}
			if o.Config.IntraDomainRealloc {
				log.Info("going to intra domain reallocation")
				if intraDomainHelp {
					allocated, err := o.intraDomainRealloc(ctx)
					if allocated {
						o.QoS = o.QoS + service.StandardQoS
						log.Info("allocated with intra domain reallocation")
						return allocated, nil
					}
					if err != nil {
						log.Info("Error in intra domain reallocation:", err)
					}
				}
			}
		} else {
			log.Info("Reallocated event not helpful")
		}
	}

	for _, nodeName := range sortedNodesNoFilter {
		node := domain.ActiveNodes[nodeName]
		reallocatedEventID, secondReallocatedEventID, thirdreallocatedEventID, err := o.getReallocatedService(node, service, o.Config.IntraNodeReducedHeu)
		if err != nil {
			continue
		}
		reallocationHelp, NewCores, _ := ReallocateTest(service, reallocatedEventID, *node)
		selectedEventID := reallocatedEventID

		if !reallocationHelp {
			reallocationHelp, NewCores, _ = ReallocateTest(service, secondReallocatedEventID, *node)
			selectedEventID = secondReallocatedEventID
		}
		if !reallocationHelp {
			reallocationHelp, NewCores, _ = ReallocateTest(service, thirdreallocatedEventID, *node)
			selectedEventID = thirdreallocatedEventID
		}
		log.Info("selected event id for reallocation:", selectedEventID)

		if reallocationHelp {
			ctx := ReallocContext{
				Service:            service,
				Node:               node,
				Domain:             domain,
				SortedNodes:        sortedNodes,
				EventID:            eventID,
				ReallocatedEventID: selectedEventID,
				NewCores:           NewCores,
			}
			if o.Config.IntraNodeReduced {
				log.Info("going to intra node cloud reallocation")
				otherEvent := o.RunningServices[ctx.ReallocatedEventID]
				if service.StandardQoS-otherEvent.StandardQoS+otherEvent.ReducedQoS > service.ReducedQoS {
					allocated, _ = o.IntraNodeReduced(ctx)
					if allocated {

						o.QoS = o.QoS + service.StandardQoS - otherEvent.StandardQoS + otherEvent.ReducedQoS
						log.Info("allocated with intra node cloud reallocation")
						return allocated, nil
					}
				}
			}
			if o.Config.IntraNodeRemoved {
				log.Info("going to intra node cloud reallocation")
				otherEvent := o.RunningServices[ctx.ReallocatedEventID]
				if service.StandardQoS-otherEvent.StandardQoS > service.ReducedQoS {
					allocated, _ = o.intraNodeRemove(ctx)
					if allocated {
						o.QoS = o.QoS + service.StandardQoS - otherEvent.StandardQoS
						log.Info("allocated with intra node cloud reallocation")
						return allocated, nil
					}
				}
			}
		} else {
			log.Info("Reallocated event not helpful")
		}

	}

	// reduced mode edge only

	// sortedNodes, _ = o.sortNodes(domain.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	// for _, nodeName := range sortedNodes {
	// 	allocated, _ := o.allocateEdgeReduced(service, domain.ActiveNodes[nodeName], eventID)
	// 	if allocated {
	// 		log.Info("qos before allocation", o.QoS)
	// 		log.Info("qos of the service", service.ReducedQoS)
	// 		log.Info("event id for edge reduced allocation:", eventID)
	// 		o.QoS = o.QoS + service.EdgeReducedQoS
	// 		return allocated, nil
	// 	}
	// }
	// sortedNodes, _ = o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	// for _, nodeName := range sortedNodes {
	// 	allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID, 100)
	// 	if allocated {
	// 		o.QoS = o.QoS + service.StandardQoS
	// 		return allocated, nil
	// 	}
	// }

	edgeAllocated, cloudAllocated, svc, _ := o.SplitSched(service, domainID, eventID)
	if edgeAllocated && cloudAllocated {
		log.Info("show svc in split scheduling", svc)
		o.QoS = o.QoS + service.ReducedQoS
		return true, nil
	} else {
		log.Info("the split scheduling didn't work. powering on some nodes")
		if !edgeAllocated {
			success, nodeName := o.edgePowerOnNode(domainID)
			log.Info("Edge node powered on, node name: ", nodeName)
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
		o.QoS = o.QoS + service.ReducedQoS
		return true, nil
	}
	// if !allocated {
	// 	o.QoS = o.QoS - QoS(service.StandardQoS/10)
	// }
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
		log.Info("Deallocating standard service: ", service.serviceID, " from node: ", node.NodeName)
		_, err := service.StandardMode.ServiceDeallocate(eventID, node)
		if err != nil {
			log.Info("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - serviceQoS
		}

	}
	if allocatedMode == ReducedMode {
		serviceQoS = service.ReducedQoS
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		cloudNode := o.Cloud.ActiveNodes[service.AllocatedNodeCloud]
		for _, n := range o.Cloud.ActiveNodes {
			log.Info("average consumed bandwidth in cloud nodes, before deallocating", n.AverageConsumedBandwidth)
		}
		_, err := service.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
		_, err = service.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
		if err != nil {
			log.Info("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - service.ReducedQoS
		}
		for _, n := range o.Cloud.ActiveNodes {
			log.Info("average consumed bandwidth in cloud nodes, after deallocating", n.AverageConsumedBandwidth)
		}

	}
	if allocatedMode == EdgeReducedMode {
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		_, err := service.ReducedMode.EdgeServiceDeallocate(eventID, edgeNode)
		if err != nil {
			log.Info("Error in deallocation: ", err)
		} else {
			o.QoS = o.QoS - service.EdgeReducedQoS
		}

	}
	// for nodeName, node := range domain.ActiveNodes {
	// 	if node.AverageConsumedBandwidth == 0 && node.TotalConsumedBandwidth == 0 {
	// 		o.edgePowerOffNode(domainID, nodeName)
	// 		log.Info("Edge node powered off: ", nodeName)
	// 	}
	// }
	// for nodeName, node := range o.Cloud.ActiveNodes {
	// 	if len(o.Cloud.ActiveNodes) == 1 {
	// 		log.Info("Cannot power off the last cloud node")
	// 		break
	// 	}
	// 	if node.AverageConsumedBandwidth == 0 && node.TotalConsumedBandwidth == 0 {
	// 		o.cloudPowerOffNode(nodeName)
	// 		log.Info("Cloud node powered off: ", nodeName)
	// 	}
	// }

	delete(o.RunningServices, eventID)

	return true
}
