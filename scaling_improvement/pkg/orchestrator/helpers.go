package orchestrator

import "fmt"

func (o *Orchestrator) NodeReclaim(domainID DomainID) {
	const cpuThreshold = 100.0
	domain := o.Domains[domainID]
	for nodeName, node := range domain.ActiveNodes {
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.edgePowerOffNode(domainID, nodeName)
			fmt.Println("node powered off:", nodeName)
		}
	}

	for nodeName, node := range o.Cloud.ActiveNodes {
		if len(o.Cloud.ActiveNodes) == 1 {
			break
		}
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.cloudPowerOffNode(nodeName)
			fmt.Println("node powered off:", nodeName)
		}
	}

	// advanced node reclaim
	totalUnderloadedNodes := 0
	var underloadedNodes []NodeName
	for nodeName, node := range domain.ActiveNodes {
		if node.AverageResidualBandwidth < 0.5 {
			totalUnderloadedNodes++
			underloadedNodes = append(underloadedNodes, nodeName)
		}
	}
	sortedNodes, _ := o.sortNodesNoFilter(domain.ActiveNodes, MinMin)
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
		if node.AverageResidualBandwidth < 0.4 {
			// for _, otherNodeName := range domain.AlwaysActiveNodes {
			fmt.Println("nodes underloaded:", nodeName)
			for j := l - 1; j > i; j-- {
				otherNodeName := sortedNodes[j]
				otherNode := domain.ActiveNodes[otherNodeName]
				if otherNode.AverageResidualBandwidth < 0.5 {
					fmt.Println("other node underloaded:", otherNodeName)
					allocatedService := node.AllocatedServices
					sortedServices := o.sortServicesBW(allocatedService)
					for _, eventID := range sortedServices {
						service := allocatedService[eventID]
						if service.AllocationMode == StandardMode {
							selectedCpus, err := node.NodeAdmission.Admission(service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge, otherNode.Cores, cpuThreshold)
							if err != nil || selectedCpus == nil {
								fmt.Println("Error in admission test for node reclaim: ", err)
								continue
							}
							service.StandardMode.ServiceDeallocate(eventID, node)

							allocated, svc, _ := service.StandardMode.ServiceAllocate(service, otherNode, eventID, cpuThreshold)
							if allocated {
								o.RunningServices[eventID] = svc
							} else {
								allAllocated = false
								fmt.Println("service was deallocated and not allocated to other node for node reclaim")
							}
						}
						if service.AllocationMode == ReducedMode {
							selectedCpus, err := node.NodeAdmission.Admission(service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge, otherNode.Cores, cpuThreshold)
							if err != nil || selectedCpus == nil {
								fmt.Println("Error in admission test for node reclaim: ", err)
								continue
							}
							service.ReducedMode.ServiceDeallocate(eventID, node, edgeLoc)

							allocated, svc, _ := service.ReducedMode.ServiceAllocate(service, otherNode, edgeLoc, eventID, cpuThreshold)
							if allocated {
								o.RunningServices[eventID] = svc
							} else {
								allAllocated = false
								fmt.Println("service was deallocated and not allocated to other node for node reclaim")
							}
						}
					}
					if allAllocated {
						if node.AllocatedServices == nil {
							nodeToPowerOff = append(nodeToPowerOff, nodeName)
							fmt.Println("node to power off:", nodeName)
						} else {
							fmt.Println("not all services were deallocated for node reclaim")
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

func (o *Orchestrator) UpgradeService() error {
	sortedEventIDs := o.sortServicesForUpgrade(o.RunningServices)
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
					fmt.Println("Error in admission test for upgrading: ", err)
					continue
				}

				_, err = oldEvent.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
				_, err = oldEvent.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
				if err != nil {
					fmt.Println("Error in deallocation: ", err)
				}
				_, svc, err := event.StandardMode.ServiceAllocate(event, domain.ActiveNodes[nodeName], eventID, 100)
				domain.ActiveNodes[nodeName].AllocatedServices[eventID] = svc
				if err != nil {
					fmt.Println("Error in allocation upgrade: ", err)
				}
				o.QoS = o.QoS - event.ReducedQoS + event.StandardQoS
				o.RunningServices[eventID] = svc
				fmt.Println("the upgraded service:", svc)
				oldEvent = nil
				fmt.Println("upgrade successful")
				return nil

			}
		}
	}
	return nil
}

func (o *Orchestrator) UpgradeServiceIfEnabled() {
	if o.Config.UpgradeService {
		o.UpgradeService()
	}
}

func (o *Orchestrator) NodeReclaimIfEnabled(domainID DomainID) {
	if o.Config.NodeReclaim {
		o.NodeReclaim(domainID)
	}
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
