package orchestrator

import (
	"sort"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

const (
	HQ    cnfg.Heuristic = "HQ"
	HQcC  cnfg.Heuristic = "HQcC"
	HQcB  cnfg.Heuristic = "HQcB"
	HQcCB cnfg.Heuristic = "HQcCB"
)

func (o *Orchestrator) BasicNodeReclaim(domainID DomainID) {
	domain := o.Domains[domainID]
	for nodeName, node := range domain.ActiveNodes {
		if node.AverageConsumedBandwidth == 0 && node.TotalConsumedBandwidth == 0 {
			o.edgePowerOffNode(domainID, nodeName)
			log.Info("node powered off:", nodeName)
		}
	}

	for nodeName, node := range o.Cloud.ActiveNodes {
		if len(o.Cloud.ActiveNodes) == 1 {
			break
		}
		if node.AverageConsumedBandwidth == 0 && node.TotalConsumedBandwidth == 0 {
			o.cloudPowerOffNode(nodeName)
			log.Info("node powered off:", nodeName)
		}
	}
}

func (o *Orchestrator) NodeReclaim(domainID DomainID) {
	const cpuThreshold = 100.0
	o.BasicNodeReclaim(domainID)
	domain := o.Domains[domainID]
	// advanced node reclaim
	totalUnderloadedNodes := 0
	var underloadedNodes []NodeName
	for nodeName, node := range domain.ActiveNodes {
		if node.AverageConsumedBandwidth < 0.5 {
			totalUnderloadedNodes++
			underloadedNodes = append(underloadedNodes, nodeName)
		}
	}
	sortedNodes, _ := o.sortNodesNoFilter(domain.ActiveNodes, Min)
	i := 0
	l := len(sortedNodes)
	// j := l - 1
	nodeToPowerOff := make([]NodeName, 0)
	if l == 1 {
		return
	}
	allAllocated := true
	for _, nodeName := range sortedNodes {
		node := domain.ActiveNodes[nodeName]
		// for _, nn := range domain.AlwaysActiveNodes {
		// 	if nn == nodeName {
		// 		continue
		// 	}
		// }
		// j = l - 1
		if node.AverageConsumedBandwidth < 0.4 {
			// for _, otherNodeName := range domain.AlwaysActiveNodes {
			log.Info("nodes underloaded:", nodeName)
			for j := l - 1; j > i; j-- {
				otherNodeName := sortedNodes[j]
				otherNode := domain.ActiveNodes[otherNodeName]
				if otherNode.AverageConsumedBandwidth < 0.5 {
					log.Info("other node underloaded:", otherNodeName)
					allocatedService := node.AllocatedServices
					sortedServices := o.sortServicesBW(allocatedService)
					for _, eventID := range sortedServices {
						service := allocatedService[eventID]
						if service.AllocationMode == StandardMode {
							selectedCpus, err := node.NodeAdmission.Admission(service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge, otherNode.Cores, cpuThreshold)
							if err != nil || selectedCpus == nil {
								log.Info("Error in admission test for node reclaim: ", err)
								continue
							}
							service.StandardMode.ServiceDeallocate(eventID, node)

							allocated, svc, _ := service.StandardMode.ServiceAllocate(service, otherNode, eventID, cpuThreshold)
							if allocated {
								o.RunningServices[eventID] = svc
							} else {
								allAllocated = false
								log.Info("service was deallocated and not allocated to other node for node reclaim")
							}
						}
						if service.AllocationMode == ReducedMode {
							selectedCpus, err := node.NodeAdmission.Admission(service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge, otherNode.Cores, cpuThreshold)
							if err != nil || selectedCpus == nil {
								log.Info("Error in admission test for node reclaim: ", err)
								continue
							}
							service.ReducedMode.ServiceDeallocate(eventID, node, edgeLoc)

							allocated, svc, _ := service.ReducedMode.ServiceAllocate(service, otherNode, edgeLoc, eventID, cpuThreshold)
							if allocated {
								o.RunningServices[eventID] = svc
							} else {
								allAllocated = false
								log.Info("service was deallocated and not allocated to other node for node reclaim")
							}
						}
					}
					if allAllocated {
						if node.AllocatedServices == nil {
							nodeToPowerOff = append(nodeToPowerOff, nodeName)
							log.Info("node to power off:", nodeName)
						} else {
							log.Info("not all services were deallocated for node reclaim")
						}
						break
					}
				}
			}
			i++
			if i == l-1 {
				break
			}
		}
	}
	for _, nodeName := range nodeToPowerOff {
		o.edgePowerOffNode(domainID, nodeName)
	}

}

func (o *Orchestrator) UpgradeService(Heu cnfg.Heuristic, svc Service, domainID DomainID) error {
	sortedEventIDs := SortServicesForUpgrade(o.RunningServices, Heu, svc.StandardMode.bandwidthEdge, svc.StandardMode.cpusEdge, domainID)
	for _, eventID := range sortedEventIDs {
		event := o.RunningServices[eventID]
		domain := o.Domains[event.AllocatedDomain]
		if event.AllocationMode == ReducedMode {
			sortedNodes, _ := o.sortNodes(domain.ActiveNodes, event.StandardMode.cpusEdge, event.StandardMode.bandwidthEdge)
			edgeNode := domain.ActiveNodes[event.AllocatedNodeEdge]
			cloudNode := o.Cloud.ActiveNodes[event.AllocatedNodeCloud]
			oldEvent := event
			for _, nodeName := range sortedNodes {
				node := domain.ActiveNodes[nodeName]
				selectedCPUs, err := node.NodeAdmission.Admission(event.StandardMode.cpusEdge, event.StandardMode.bandwidthEdge, node.Cores, 100.0)
				if err != nil || selectedCPUs == nil {
					log.Info("Error in admission test for upgrading: ", err)
					continue
				}

				_, err = oldEvent.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
				_, err = oldEvent.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
				if err != nil {
					log.Info("Error in deallocation: ", err)
				}
				_, svc, err := event.StandardMode.ServiceAllocate(event, domain.ActiveNodes[nodeName], eventID, 100)
				domain.ActiveNodes[nodeName].AllocatedServices[eventID] = svc
				if err != nil {
					log.Info("Error in allocation upgrade: ", err)
				}
				o.QoS = o.QoS - event.ReducedQoS + event.StandardQoS
				o.RunningServices[eventID] = svc
				log.Info("the upgraded service:", svc)
				oldEvent = nil
				log.Info("upgrade successful")
				return nil

			}
		}
	}
	return nil
}

func (o *Orchestrator) UpgradeServiceIfEnabled(Heu cnfg.Heuristic, svc Service, domainID DomainID) {
	if o.Config.UpgradeService {
		o.UpgradeService(Heu, svc, domainID)
	}
}

func (o *Orchestrator) NodeReclaimIfEnabled(domainID DomainID) {
	if o.Config.NodeReclaim {
		o.NodeReclaim(domainID)
	} else {
		o.BasicNodeReclaim(domainID)
	}
}

func SortServicesForUpgrade(services Services, upgradeHeu cnfg.Heuristic, BW float64, m uint64, domainID DomainID) []ServiceID {
	sortedEventIDs := make([]ServiceID, 0, len(services))
	for id, svc := range services {
		if svc.AllocatedDomain == domainID {
			sortedEventIDs = append(sortedEventIDs, id)
		}
	}

	switch upgradeHeu {
	case HQ:
		sort.Slice(sortedEventIDs, func(i, j int) bool {
			return float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS) >
				float64(services[sortedEventIDs[j]].StandardQoS-services[sortedEventIDs[j]].ReducedQoS)
		})

	case HQcC:
		sort.Slice(sortedEventIDs, func(i, j int) bool {
			cpuI := services[sortedEventIDs[i]].StandardMode.cpusEdge <= m
			cpuJ := services[sortedEventIDs[j]].StandardMode.cpusEdge <= m

			if cpuI && cpuJ {
				return float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS) >
					float64(services[sortedEventIDs[j]].StandardQoS-services[sortedEventIDs[j]].ReducedQoS)
			}
			return cpuI
		})

	case HQcB:
		sort.Slice(sortedEventIDs, func(i, j int) bool {
			bwI := services[sortedEventIDs[i]].StandardMode.bandwidthEdge <= BW
			bwJ := services[sortedEventIDs[j]].StandardMode.bandwidthEdge <= BW

			if bwI && bwJ {
				return float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS) >
					float64(services[sortedEventIDs[j]].StandardQoS-services[sortedEventIDs[j]].ReducedQoS)
			}
			return bwI
		})

	case HQcCB:
		sort.Slice(sortedEventIDs, func(i, j int) bool {
			cpuI := services[sortedEventIDs[i]].StandardMode.cpusEdge <= m
			cpuJ := services[sortedEventIDs[j]].StandardMode.cpusEdge <= m
			bwI := services[sortedEventIDs[i]].StandardMode.bandwidthEdge <= BW
			bwJ := services[sortedEventIDs[j]].StandardMode.bandwidthEdge <= BW

			if cpuI && cpuJ && bwI && bwJ {
				return float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS) >
					float64(services[sortedEventIDs[j]].StandardQoS-services[sortedEventIDs[j]].ReducedQoS)
			}
			if cpuI && bwI {
				return true
			}
			if cpuJ && bwJ {
				return false
			}
			return false
		})
	}

	return sortedEventIDs
}

// type AlgorithmStep func(*src.Orchestrator, src.DomainID, src.ServiceID, string)

// func withUpgradeService(next AlgorithmStep) AlgorithmStep {
//     return func(orchestrator *src.Orchestrator, domainID src.DomainID, serviceID src.ServiceID, eventID string) {
//         if orchestrator.Config.UpgradeService {
//             orchestrator.UpgradeService()
//         }
//         next(orchestrator, domainID, serviceID, eventID)
//     }
// }

// func withNodeReclaim(next AlgorithmStep) AlgorithmStep {
//     return func(orchestrator *src.Orchestrator, domainID src.DomainID, serviceID src.ServiceID, eventID string) {
//         if orchestrator.Config.NodeReclaim {
//             orchestrator.NodeReclaim(domainID)
//         }
//         next(orchestrator, domainID, serviceID, eventID)
//     }
// }
