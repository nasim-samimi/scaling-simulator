package orchestrator

import (
	"fmt"
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
		fmt.Println("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageResidualBandwidth, "total residual bandwidth: ", node.TotalResidualBandwidth)
	}

	fmt.Println("inside switch case", o.Config.NodeHeuristic)
	switch o.Config.NodeHeuristic {

	case MinMin:
		// Sort by number of cores (descending) first, then by average residual bandwidth (ascending)
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
			}
			return len(sortedNodes[i].Cores) > len(sortedNodes[j].Cores)
		})

	case MaxMax:
		// Sort by number of cores (descending) first, then by average residual bandwidth (descending)
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
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

func (o *Orchestrator) sortNodesNoFilter(nodes Nodes, nodeSelection config.Heuristic) ([]NodeName, error) {
	// sort nodes according to the heuristic
	sortedNodes := []Node{}
	for _, node := range nodes {
		// must filter out the nodes that do not pass admission test
		sortedNodes = append(sortedNodes, *node)
		fmt.Println("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageResidualBandwidth, "total residual bandwidth: ", node.TotalResidualBandwidth)
	}

	switch nodeSelection {

	case MinMin:
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
			}
			return len(sortedNodes[i].Cores) < len(sortedNodes[j].Cores)
		})

	// 	sort.Slice(sortedNodes, func(i, j int) bool {
	// 		return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
	// 	})
	case MaxMax:
		// 	sort.Slice(sortedNodes, func(i, j int) bool {
		// 		return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
		// 	})
		// } //changed to maxmax strategy
		sort.Slice(sortedNodes, func(i, j int) bool {
			if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
				return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
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
