package orchestrator

import (
	"sort"

	"github.com/nasim-samimi/scaling-simulator/pkg/config"
)

func (o *Orchestrator) sortNodes(nodes Nodes, serviceCpus uint64, serviceBandwidth float64) ([]NodeName, error) {
	// sort nodes according to the heuristic
	sortedNodes := []Node{}
	for _, node := range nodes {
		// must filter out the nodes that do not pass admission test
		available, _ := node.NodeAdmission.QuickFilter(serviceCpus, serviceBandwidth, node.Cores)
		if available {
			sortedNodes = append(sortedNodes, *node)
		}
		log.Info("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageConsumedBandwidth, "total residual bandwidth: ", node.TotalConsumedBandwidth)
		for _, core := range node.Cores {
			log.Info("Core: ", core.ID, " Consumed Bandwidth: ", core.ConsumedBandwidth)
		}
	}

	log.Info("inside switch case", o.Config.NodeHeuristic)
	switch o.Config.NodeHeuristic {

	case MinMin:
		// Sort by number of cores (descending) first, then by average residual bandwidth (ascending)
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageConsumedBandwidth < sortedNodes[j].AverageConsumedBandwidth
			}
			return len(sortedNodes[i].Cores) > len(sortedNodes[j].Cores)
		})

	case MaxMax:
		// Sort by number of cores (descending) first, then by average residual bandwidth (descending)
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageConsumedBandwidth > sortedNodes[j].AverageConsumedBandwidth
			}
			return len(sortedNodes[i].Cores) > len(sortedNodes[j].Cores)
		})

	case MMRB:
		// get the node with a core that the core has highest residual bandwidth among all cores in all nodes
		sort.Slice(sortedNodes, func(i, j int) bool {
			return getMaxResBWCore(sortedNodes[i]) > getMaxResBWCore(sortedNodes[j])
		})

	case mmRB:
		// get the node with a core that the core has lowest residual bandwidth among all cores in all nodes
		sort.Slice(sortedNodes, func(i, j int) bool {
			return getMinResBWCore(sortedNodes[i]) < getMinResBWCore(sortedNodes[j])
		})
	case mMRB:
		// get the node with a core that the core has highest residual bandwidth among all cores in all nodes
		sort.Slice(sortedNodes, func(i, j int) bool {
			return getMaxResBWCore(sortedNodes[i]) < getMaxResBWCore(sortedNodes[j])
		})
	case MmRB:
		// get the node with a core that the core has lowest residual bandwidth among all cores in all nodes
		sort.Slice(sortedNodes, func(i, j int) bool {
			return getMinResBWCore(sortedNodes[i]) > getMinResBWCore(sortedNodes[j])
		})

	}
	// Extract sorted NodeNames
	sortedNodeNames := make([]NodeName, len(sortedNodes))
	for i, node := range sortedNodes {
		sortedNodeNames[i] = node.NodeName
	}
	log.Info("sorted nodes: ", sortedNodeNames)
	return sortedNodeNames, nil
}

func (o *Orchestrator) sortNodesNoFilter(nodes Nodes, nodeSelection config.Heuristic) ([]NodeName, error) {
	// sort nodes according to the heuristic
	sortedNodes := []Node{}
	for _, node := range nodes {
		// must filter out the nodes that do not pass admission test
		sortedNodes = append(sortedNodes, *node)
		log.Info("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageConsumedBandwidth, "total residual bandwidth: ", node.TotalConsumedBandwidth)
	}

	switch nodeSelection {

	case MinMin:
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageConsumedBandwidth > sortedNodes[j].AverageConsumedBandwidth
			}
			return len(sortedNodes[i].Cores) < len(sortedNodes[j].Cores)
		})

	// 	sort.Slice(sortedNodes, func(i, j int) bool {
	// 		return sortedNodes[i].AverageConsumedBandwidth < sortedNodes[j].AverageConsumedBandwidth
	// 	})
	case MaxMax:
		// 	sort.Slice(sortedNodes, func(i, j int) bool {
		// 		return sortedNodes[i].AverageConsumedBandwidth > sortedNodes[j].AverageConsumedBandwidth
		// 	})
		// } //changed to maxmax strategy
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageConsumedBandwidth > sortedNodes[j].AverageConsumedBandwidth
			}
			return len(sortedNodes[i].Cores) > len(sortedNodes[j].Cores)
		})
	}
	// Extract sorted NodeNames
	sortedNodeNames := make([]NodeName, len(sortedNodes))
	for i, node := range sortedNodes {
		sortedNodeNames[i] = node.NodeName
	}
	return sortedNodeNames, nil
}

func (o *Orchestrator) sortServicesBW(services Services) []ServiceID {
	sortedEventIDs := make([]ServiceID, 0, len(services))
	for id := range services {
		sortedEventIDs = append(sortedEventIDs, id)
	}
	sort.Slice(sortedEventIDs, func(i, j int) bool {
		return services[sortedEventIDs[i]].StandardMode.bandwidthEdge > services[sortedEventIDs[j]].StandardMode.bandwidthEdge
	})
	return sortedEventIDs
}

func (o *Orchestrator) sortServicesForUpgrade(services Services) []ServiceID {
	sortedEventIDs := make([]ServiceID, 0, len(services))
	for id := range services {
		sortedEventIDs = append(sortedEventIDs, id)
	}
	sort.Slice(sortedEventIDs, func(i, j int) bool {
		return float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS) > float64(services[sortedEventIDs[i]].StandardQoS-services[sortedEventIDs[i]].ReducedQoS)
	})
	return sortedEventIDs
}

func getMinResBWCore(node Node) float64 {
	maxBW := 0.0
	for _, core := range node.Cores {
		if core.ConsumedBandwidth > maxBW {
			maxBW = core.ConsumedBandwidth
		}
	}
	return (100.0 - maxBW)
}

func getMaxResBWCore(node Node) float64 {
	minBW := 100.0
	for _, core := range node.Cores {
		if core.ConsumedBandwidth < minBW {
			minBW = core.ConsumedBandwidth
		}
	}
	return (100.0 - minBW)
}
