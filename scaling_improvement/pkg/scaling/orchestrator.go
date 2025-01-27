package scaling

import (
	"fmt"
	"sort"
)

type Heuristic string
type NodeSelectionHeuristic Heuristic
type ReallocationHeuristic Heuristic
type Location string

const (
	cloudLoc Location = "cloud"
	edgeLoc  Location = "edge"
	bothLoc  Location = "both"
)

const (
	HBI   ReallocationHeuristic = "HBI"
	HCI   ReallocationHeuristic = "HCI"
	HBCI  ReallocationHeuristic = "HBCI"
	HB    ReallocationHeuristic = "HB"
	HC    ReallocationHeuristic = "HC"
	HBC   ReallocationHeuristic = "HBC"
	LB    ReallocationHeuristic = "LB"
	LC    ReallocationHeuristic = "LC"
	LBC   ReallocationHeuristic = "LBC"
	HCLI  ReallocationHeuristic = "HCLI"
	HBLI  ReallocationHeuristic = "HBLI"
	HBIcC ReallocationHeuristic = "HBIcC"
	LBCI  ReallocationHeuristic = "LBCI"
)

const (
	MinMin NodeSelectionHeuristic = "MinMin"
	MaxMax NodeSelectionHeuristic = "MaxMax"
)

type QoS int
type Cost int
type EventID string

const (
	CloudNodeCost Cost = 3
	EdgeNodeCost  Cost = 1
)

type Orchestrator struct {
	NodeSelectionHeuristic NodeSelectionHeuristic
	ReallocationHeuristic  ReallocationHeuristic
	PartitionHeuristic     Heuristic
	Domains                Domains
	Cloud                  *Cloud
	AllServices            Services
	RunningServices        Services // change name of service to service
	Cost                   Cost
	QoS                    QoS
	sortedNodes            []NodeName
}

func NewOrchestrator(nodeSelectionHeuristic NodeSelectionHeuristic, reallocationHeuristic ReallocationHeuristic, partitionHeuristic Heuristic, cloud *Cloud, domains Domains, services Services) *Orchestrator {
	domainCost := Cost(0)
	for _, d := range domains {
		domainCost += Cost(len(d.ActiveNodes)) * EdgeNodeCost * 2
	}
	fmt.Println("Domain cost:", domainCost)
	cloudCost := Cost(len(cloud.ActiveNodes)) * CloudNodeCost

	cost := domainCost + cloudCost
	o := &Orchestrator{
		NodeSelectionHeuristic: nodeSelectionHeuristic,
		ReallocationHeuristic:  reallocationHeuristic,
		PartitionHeuristic:     partitionHeuristic,
		Domains:                domains,
		Cloud:                  cloud,
		AllServices:            services,
		Cost:                   cost,
		QoS:                    0,
		RunningServices:        make(Services),
	}

	o.cloudPowerOnNode()
	// for domainID := range o.Domains {
	// 	o.edgePowerOnNode(domainID)
	// }
	return o

}

func (o *Orchestrator) SelectNode(service *Service) *Node {
	var node *Node

	return node
}

func (o *Orchestrator) allocateEdge(service *Service, node *Node, eventID ServiceID) (bool, error) {
	const cpuThreshold = 100.0 //80.0
	fmt.Println("Allocating standard service: ", service.serviceID, " to node: ", node.NodeName)
	allocated, svc, err := service.StandardMode.ServiceAllocate(service, node, eventID, cpuThreshold)
	if allocated {
		o.RunningServices[eventID] = svc
		fmt.Println("service:", svc)
		fmt.Println("Allocated? ", allocated)
	}

	return allocated, err
}

func (o *Orchestrator) allocateEdgeReduced(service *Service, node *Node, eventID ServiceID) (bool, error) {
	const cpuThreshold = 100.0 //80.0

	allocated, svc, err := service.ReducedMode.EdgeServiceAllocate(service, node, eventID, cpuThreshold)
	if allocated {
		o.RunningServices[eventID] = svc
		fmt.Println("service:", svc)
		fmt.Println("Allocated? ", allocated)
	}

	return allocated, err
}

func (o *Orchestrator) getReallocatedService(node *Node, t *Service) (ServiceID, error) {
	var selectedEventID ServiceID
	var bestScore float64

	calculateScore := func(service *Service, heuristic ReallocationHeuristic) float64 {
		switch heuristic {
		case HB:
			return (service.StandardMode.bandwidthEdge)
		case HC:
			return float64(service.StandardMode.cpusEdge)
		case HBC:
			return 1 / (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
		case LB:
			return (1 / service.StandardMode.bandwidthEdge)
		case LC:
			return 1 / float64(service.StandardMode.cpusEdge)
		case LBC:
			return (service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
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

	for eventID, service := range node.AllocatedServices {
		if service.AllocationMode == StandardMode {
			score := calculateScore(service, o.ReallocationHeuristic)
			if score > bestScore {
				bestScore = score
				selectedEventID = eventID
			}
		}
	}

	if selectedEventID == "" {
		return "", fmt.Errorf("no suitable service found for reallocation using heuristic %s", o.ReallocationHeuristic)
	}

	return selectedEventID, nil
}

func (o *Orchestrator) sortNodes(nodes Nodes, serviceCpus uint64, serviceBandwidth float64) ([]NodeName, error) {
	// sort nodes according to the heuristic
	fmt.Println("inside sort nodes")
	fmt.Println("checking the active nodes for sorting:", nodes)
	sortedNodes := []Node{}
	for _, node := range nodes {
		// must filter out the nodes that do not pass admission test
		available, _ := node.NodeAdmission.QuickFilter(serviceCpus, serviceBandwidth, node.Cores)
		if available {
			sortedNodes = append(sortedNodes, *node)
		}
		for _, core := range node.Cores {
			fmt.Println("Node cores: ", core)
		}
		fmt.Println("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageResidualBandwidth, "total residual bandwidth: ", node.TotalResidualBandwidth)
	}

	fmt.Println("inside switch case", o.NodeSelectionHeuristic)
	switch o.NodeSelectionHeuristic {

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

func (o *Orchestrator) sortNodesNoFilter(nodes Nodes) ([]NodeName, error) {
	// sort nodes according to the heuristic
	sortedNodes := []Node{}
	for _, node := range nodes {
		// must filter out the nodes that do not pass admission test
		sortedNodes = append(sortedNodes, *node)
		fmt.Println("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageResidualBandwidth, "total residual bandwidth: ", node.TotalResidualBandwidth)
	}

	// switch o.NodeSelectionHeuristic {

	// case MinMin:

	// 	sort.Slice(sortedNodes, func(i, j int) bool {
	// 		return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
	// 	})
	// case MaxMax:
	// 	sort.Slice(sortedNodes, func(i, j int) bool {
	// 		return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
	// 	})
	// }
	sort.Slice(sortedNodes, func(i, j int) bool {
		if len(sortedNodes[i].Cores) == len(sortedNodes[j].Cores) {
			return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
		}
		return len(sortedNodes[i].Cores) > len(sortedNodes[j].Cores)
	})
	// Extract sorted NodeNames
	sortedNodeNames := make([]NodeName, len(sortedNodes))
	for i, node := range sortedNodes {
		sortedNodeNames[i] = node.NodeName
	}
	return sortedNodeNames, nil
}

func (o *Orchestrator) intraNodeRealloc(service *Service, node *Node, eventID ServiceID, reallocatedEventID ServiceID, NewCores Cores) (bool, error) {
	const cpuThreshold = 100.0
	if reallocatedEventID == "" {
		return false, nil
	}
	oldBandwidth := o.RunningServices[reallocatedEventID].StandardMode.bandwidthEdge
	_, err := node.NodeAdmission.Admission(node.AllocatedServices[reallocatedEventID].StandardMode.cpusEdge, oldBandwidth, NewCores, cpuThreshold)

	if err != nil {
		return false, err
	}
	fmt.Println("show the status of the intra node reallocation:", true)

	reallocatedService := node.AllocatedServices[reallocatedEventID]
	_, err = service.StandardMode.ServiceDeallocate(reallocatedEventID, node)
	if err != nil {
		return false, err
	}
	_, newSvc, _ := service.StandardMode.ServiceAllocate(service, node, eventID, cpuThreshold)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", node.AverageResidualBandwidth)
	_, oldSvc, _ := reallocatedService.StandardMode.ServiceAllocate(reallocatedService, node, reallocatedEventID, cpuThreshold)
	fmt.Println("Allocated services in the end: ", node.AllocatedServices)
	// delete(o.RunningServices, reallocatedEventID)
	o.RunningServices[reallocatedEventID] = oldSvc
	o.RunningServices[eventID] = newSvc
	fmt.Println("intra node reallocation completed")

	return true, nil
}

func (o *Orchestrator) intraDomainRealloc(service *Service, node *Node, domain *Domain, sortedNodes []NodeName, eventID ServiceID, otherEventID ServiceID) (bool, error) {
	const cpuThreshold = 100.0
	if otherEventID == "" {
		return false, nil
	}

	reallocated := false

	otherService := node.AllocatedServices[otherEventID]
	fmt.Println("inside the intra domain reallocation")
	for _, nodeName := range sortedNodes {
		if nodeName == node.NodeName {
			continue
		}
		otherNode := domain.ActiveNodes[nodeName]
		fmt.Println("other event id:", otherEventID)
		fmt.Println("other node:", otherNode.NodeName)
		fmt.Println("other service:", otherService)
		for _, core := range otherNode.Cores {
			fmt.Println("cores of the other node:", core)
		}

		allocatedCore, _ := otherNode.NodeAdmission.Admission(otherService.StandardMode.cpusEdge, otherService.StandardMode.bandwidthEdge, otherNode.Cores, cpuThreshold)

		fmt.Println("allocated core:", allocatedCore)
		if allocatedCore != nil {
			fmt.Println("reallocation successful")

			otherService.StandardMode.ServiceDeallocate(otherEventID, node)
			_, newSvc, _ := service.StandardMode.ServiceAllocate(service, node, eventID, cpuThreshold)
			_, oldSvc, _ := otherService.StandardMode.ServiceAllocate(otherService, otherNode, otherEventID, cpuThreshold)
			// delete(o.RunningServices, otherEventID)
			fmt.Println(" the services of the node:", node.AllocatedServices)
			fmt.Println(" the services of the other node:", otherNode.AllocatedServices)
			fmt.Println("the new service:", newSvc)
			fmt.Println("the old service:", oldSvc)
			o.RunningServices[otherEventID] = oldSvc
			o.RunningServices[eventID] = newSvc
			reallocated = true
			fmt.Println("intra domain reallocation completed")
			return reallocated, nil
		}
	}
	return reallocated, fmt.Errorf("intra domain reallocation failed")
}

func (o *Orchestrator) IntraNodeCloudRealloc(service *Service, node *Node, eventID ServiceID, reallocatedEventID ServiceID, NewCores Cores) (bool, error) {
	const cpuThreshold = 100.0
	if reallocatedEventID == "" {
		return false, nil
	}
	otherService := node.AllocatedServices[reallocatedEventID]
	oldBandwidth := o.RunningServices[reallocatedEventID].ReducedMode.bandwidthCloud
	oldCpus := o.RunningServices[reallocatedEventID].ReducedMode.cpusCloud
	selectedCPUs, err := node.NodeAdmission.Admission(oldCpus, oldBandwidth, NewCores, cpuThreshold)

	if err != nil || selectedCPUs == nil {
		return false, err
	}
	// cloudAllocated := false
	var svcEdge, svcCloud *Service
	var cloudNode NodeName
	sortedNodes, _ := o.sortNodes(o.Cloud.ActiveNodes, otherService.ReducedMode.cpusCloud, otherService.ReducedMode.bandwidthCloud)
	for _, cloudNodeName := range sortedNodes {
		selectedCPUs, err := node.NodeAdmission.Admission(otherService.ReducedMode.cpusCloud, otherService.ReducedMode.bandwidthCloud, o.Cloud.ActiveNodes[cloudNodeName].Cores, cpuThreshold)
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

	_, err = otherService.StandardMode.ServiceDeallocate(reallocatedEventID, node)
	if err != nil {
		return false, err
	}
	allocated, newSvc, _ := service.StandardMode.ServiceAllocate(service, node, eventID, cpuThreshold)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", node.AverageResidualBandwidth)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation")
	}
	allocated, svcEdge, _ = otherService.ReducedMode.ServiceAllocate(otherService, node, edgeLoc, reallocatedEventID, cpuThreshold)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation in edge")
	}
	allocated, svcCloud, _ = otherService.ReducedMode.ServiceAllocate(otherService, o.Cloud.ActiveNodes[cloudNode], cloudLoc, reallocatedEventID, cpuThreshold)
	if !allocated {
		return false, fmt.Errorf("service not allocated in intra node cloud reallocation in cloud")
	}

	fmt.Println("Allocated services in the end: ", node.AllocatedServices)
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
	node.AllocatedServices[reallocatedEventID] = oldSvc
	svcEdge = nil
	svcCloud = nil
	o.RunningServices[reallocatedEventID] = oldSvc
	o.RunningServices[eventID] = newSvc
	fmt.Println("intra node cloud reallocation completed")
	return true, nil
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

func (o *Orchestrator) edgePowerOffNode(domainID DomainID, nodeName NodeName) bool {
	for _, noden := range o.Domains[domainID].AlwaysActiveNodes {
		if noden == nodeName {
			return false
		}
	}
	cores := CreateNodeCores(len(o.Domains[domainID].ActiveNodes[nodeName].Cores))
	o.Domains[domainID].InactiveNodes[nodeName] = NewNode(cores, o.Domains[domainID].ActiveNodes[nodeName].ReallocHeuristic, nodeName, domainID)
	o.Cost = o.Cost - EdgeNodeCost
	delete(o.Domains[domainID].ActiveNodes, nodeName)

	return true
}

func (o *Orchestrator) cloudPowerOffNode(nodeName NodeName) bool {
	o.Cloud.InactiveNodes[nodeName] = NewNode(o.Cloud.ActiveNodes[nodeName].Cores, o.Cloud.ActiveNodes[nodeName].ReallocHeuristic, nodeName, "")
	o.Cost = o.Cost - CloudNodeCost
	delete(o.Cloud.ActiveNodes, nodeName)

	return true
}

func (o *Orchestrator) edgePowerOnNode(domainID DomainID) (bool, NodeName) {
	// fmt.Println("active nodes in domain:", o.Domains[domainID].ActiveNodes)
	// var nodeName NodeName
	for nodeName, node := range o.Domains[domainID].InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Domains[domainID].ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, domainID)
		o.Cost = o.Cost + EdgeNodeCost
		delete(o.Domains[domainID].InactiveNodes, nodeName)
		return true, nodeName
	}
	// fmt.Println("active nodes in domain after powering on:", o.Domains[domainID].ActiveNodes)
	return false, ""
}
func (o *Orchestrator) cloudPowerOnNode() bool {
	// fmt.Println("active nodes in cloud:", o.Cloud.ActiveNodes)
	for nodeName, node := range o.Cloud.InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Cloud.ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, "")
		o.Cost = o.Cost + CloudNodeCost
		delete(o.Cloud.InactiveNodes, nodeName)
		break
	}
	// fmt.Println("active nodes in cloud after powering on:", o.Cloud.ActiveNodes)

	return true
}

func (o *Orchestrator) Allocate(domainID DomainID, serviceID ServiceID, eventID ServiceID) (bool, error) {
	allocated := false
	domain := o.Domains[domainID]
	service := o.AllServices[serviceID]

	sortedNodes, _ := o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID)
		if allocated {
			fmt.Println("qos before standard allocation", o.QoS)
			fmt.Println("qos of the service", service.StandardQoS)
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
			fmt.Println("node name:", nodeName)
			fmt.Println("node allocated services:", node.AllocatedServices)
			fmt.Println("node BW:", node.TotalResidualBandwidth, node.AverageResidualBandwidth)
			continue
		}
		reallocationHelp, NewCores, _ := ReallocateTest(service, reallocatedEventID, *node)
		if reallocationHelp {
			allocated, err := o.intraNodeRealloc(service, node, eventID, reallocatedEventID, NewCores)
			if allocated {
				fmt.Println("qos before intra node reallocation", o.QoS)
				fmt.Println("qos of the service", service.StandardQoS)
				o.QoS = o.QoS + service.StandardQoS
				fmt.Println("allocated with intra node reallocation")
				fmt.Println("node average residual bandwidth after allocation:", o.Domains[domainID].ActiveNodes[nodeName].AverageResidualBandwidth, "total residual bandwidth:", o.Domains[domainID].ActiveNodes[nodeName].TotalResidualBandwidth)
				return allocated, nil
			}
			if err != nil {
				fmt.Println("Error in intra node reallocation: ", err)
			}
			fmt.Println("going to intra domain reallocation")
			if intraDomainHelp {
				allocated, err = o.intraDomainRealloc(service, node, domain, sortedNodesNoFilter, eventID, reallocatedEventID)
				if allocated {
					fmt.Println("qos before intra domain reallocation", o.QoS)
					fmt.Println("qos of the service", service.StandardQoS)
					o.QoS = o.QoS + service.StandardQoS
					fmt.Println("allocated with intra domain reallocation")
					return allocated, nil
				}
				if err != nil {
					fmt.Println("Error in intra domain reallocation:", err)
				}
			}
			// intraDomainHelp = false
			otherEvent := o.RunningServices[reallocatedEventID]

			fmt.Println("going to intra node cloud reallocation")
			allocated, _ = o.IntraNodeCloudRealloc(service, node, eventID, reallocatedEventID, NewCores)
			if allocated {
				o.QoS = o.QoS + service.StandardQoS - otherEvent.StandardQoS + otherEvent.ReducedQoS
				fmt.Println("allocated with intra node cloud reallocation")
				return allocated, nil
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
				allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID)
				if allocated {
					fmt.Println("qos before edge power on", o.QoS)
					fmt.Println("qos of the service", service.StandardQoS)
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
		fmt.Println("qos before split sched", o.QoS)
		fmt.Println("qos of the service", service.StandardQoS)
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
	fmt.Println("event id for deallocation:", eventID)
	fmt.Println(service.StandardMode)
	fmt.Println(service.ReducedMode)

	if allocatedMode == StandardMode {
		serviceQoS = service.StandardQoS
		nodeN := service.AllocatedNodeEdge
		node := domain.ActiveNodes[nodeN]
		_, err := service.StandardMode.ServiceDeallocate(eventID, node)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		} else {
			fmt.Println("qos before deallocate", o.QoS)
			fmt.Println("qos of the service", service.StandardQoS)
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
			fmt.Println("qos before reduced deallocation", o.QoS)
			fmt.Println("qos of the service", service.ReducedQoS)
			fmt.Println("qos of the service", service.StandardQoS)
			o.QoS = o.QoS - service.ReducedQoS
		}

	}
	if allocatedMode == EdgeReducedMode {
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		_, err := service.ReducedMode.EdgeServiceDeallocate(eventID, edgeNode)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		} else {
			fmt.Println("qos before reduced deallocation", o.QoS)
			fmt.Println("qos of the service", service.ReducedQoS)
			fmt.Println("qos of the service", service.StandardQoS)
			fmt.Println("qos of the service", service.EdgeReducedQoS)
			fmt.Println("event id for edge reduced deallocation:", eventID)
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
	sortedNodes, _ := o.sortNodesNoFilter(domain.ActiveNodes)
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
					for eventID, service := range allocatedService {
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
	for eventID, event := range o.RunningServices {
		domain := o.Domains[event.AllocatedDomain]
		if event.AllocationMode == ReducedMode {
			fmt.Println("service to upgrade:", event)
			fmt.Println("domain id:", event.AllocatedDomain)
			fmt.Println("domain active nodes:", domain.ActiveNodes)
			fmt.Println("bandwidth edge:", event.StandardMode.bandwidthEdge)
			fmt.Println("cpus edge:", event.StandardMode.cpusEdge)
			sortedNodes, _ := o.sortNodes(domain.ActiveNodes, event.StandardMode.cpusEdge, event.StandardMode.bandwidthEdge)
			edgeNode := domain.ActiveNodes[event.AllocatedNodeEdge]
			fmt.Println("edge node in upgrade service:", edgeNode)
			cloudNode := o.Cloud.ActiveNodes[event.AllocatedNodeCloud]
			fmt.Println("cloud node in upgrade service:", cloudNode)
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
				// o.RunningServices[eventID] = svc
				// o.NodeReclaim(domain.DomainID)
				oldEvent = nil
				fmt.Println("upgrade successful")
				return nil

			}
		}
	}
	return nil
}
