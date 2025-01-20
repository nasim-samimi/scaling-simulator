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
	LBHC  ReallocationHeuristic = "LBHC"
	HBLC  ReallocationHeuristic = "HBLC"
	HBIcC ReallocationHeuristic = "HBIcC"
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
}

func NewOrchestrator(nodeSelectionHeuristic NodeSelectionHeuristic, reallocationHeuristic ReallocationHeuristic, partitionHeuristic Heuristic, cloud *Cloud, domains Domains, services Services) *Orchestrator {
	domainCost := Cost(0)
	for _, d := range domains {
		domainCost += Cost(len(d.ActiveNodes)) * EdgeNodeCost
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

func (o *Orchestrator) allocateEdge(service *Service, node *Node, eventID ServiceID) (bool, *Service, error) {
	fmt.Println("Allocating standard service: ", service.serviceID, " to node: ", node.NodeName)
	allocated, svc, err := service.StandardMode.ServiceAllocate(node, eventID)
	fmt.Println("service:", svc)
	fmt.Println("Allocated? ", allocated)
	return allocated, svc, err
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
		case HBLC:
			return (service.StandardMode.bandwidthEdge * float64(1/service.StandardMode.cpusEdge))
		case LBHC:
			return (1 / service.StandardMode.bandwidthEdge * float64(service.StandardMode.cpusEdge))
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

		sort.Slice(sortedNodes, func(i, j int) bool {
			return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
		})
	case MaxMax:
		sort.Slice(sortedNodes, func(i, j int) bool {
			return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
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

	switch o.NodeSelectionHeuristic {

	case MinMin:

		sort.Slice(sortedNodes, func(i, j int) bool {
			return sortedNodes[i].AverageResidualBandwidth < sortedNodes[j].AverageResidualBandwidth
		})
	case MaxMax:
		sort.Slice(sortedNodes, func(i, j int) bool {
			return sortedNodes[i].AverageResidualBandwidth > sortedNodes[j].AverageResidualBandwidth
		})
	}
	// Extract sorted NodeNames
	sortedNodeNames := make([]NodeName, len(sortedNodes))
	for i, node := range sortedNodes {
		sortedNodeNames[i] = node.NodeName
	}
	return sortedNodeNames, nil
}

func (o *Orchestrator) intraNodeRealloc(service *Service, node *Node, eventID ServiceID, reallocatedEventID ServiceID, NewCores Cores) (bool, error) {

	if reallocatedEventID == "" {
		return false, nil
	}
	oldBandwidth := o.RunningServices[reallocatedEventID].StandardMode.bandwidthEdge
	_, err := node.NodeAdmission.Admission(node.AllocatedServices[reallocatedEventID].StandardMode.cpusEdge, oldBandwidth, NewCores)

	if err != nil {
		return false, err
	}
	fmt.Println("show the status of the intra node reallocation:", true)

	reallocatedService := node.AllocatedServices[reallocatedEventID]
	_, err = service.StandardMode.ServiceDeallocate(reallocatedEventID, node)
	if err != nil {
		return false, err
	}
	_, newSvc, _ := service.StandardMode.ServiceAllocate(node, eventID)
	fmt.Println("in inra node reallocation, node average residual bandwidth after first allocation: ", node.AverageResidualBandwidth)
	_, oldSvc, _ := reallocatedService.StandardMode.ServiceAllocate(node, reallocatedEventID)
	fmt.Println("Allocated services in the end: ", node.AllocatedServices)
	o.RunningServices[reallocatedEventID] = oldSvc
	o.RunningServices[eventID] = newSvc
	fmt.Println("intra node reallocation completed")

	return true, nil
}

func (o *Orchestrator) intraDomainRealloc(service *Service, node *Node, domain *Domain, sortedNodes []NodeName, eventID ServiceID, otherEventID ServiceID) (bool, error) {
	if otherEventID == "" {
		return false, nil
	}

	reallocated := false

	otherService := node.AllocatedServices[otherEventID]

	for _, nodeName := range sortedNodes {
		if nodeName == node.NodeName {
			continue
		}
		otherNode := domain.ActiveNodes[nodeName]

		allocatedCore, _ := otherNode.NodeAdmission.Admission(service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge, otherNode.Cores)
		if allocatedCore != nil {
			fmt.Println("reallocation successful")

			otherService.StandardMode.ServiceDeallocate(otherEventID, node)
			_, newSvc, _ := service.StandardMode.ServiceAllocate(node, eventID)
			_, oldSvc, _ := otherService.StandardMode.ServiceAllocate(otherNode, otherEventID)
			o.RunningServices[otherEventID] = oldSvc
			o.RunningServices[eventID] = newSvc
			reallocated = true
			fmt.Println("intra domain reallocation completed")
		}
		return reallocated, nil
	}
	fmt.Println("intra domain reallocation failed")
	return reallocated, nil
}

func (o *Orchestrator) SplitSched(service *Service, domainID DomainID, eventID ServiceID) (bool, bool, *Service, error) {
	// edge-cloud split (has qos degradation) -- there is no cloud only apparently
	fmt.Println("inside split scheduling")
	sortedNodes, _ := o.sortNodes(o.Domains[domainID].ActiveNodes, service.ReducedMode.cpusEdge, service.ReducedMode.bandwidthEdge)
	edgeAllocated := false
	cloudAllocated := false
	var svcEdge, svcCloud *Service

	for _, edgeNodeName := range sortedNodes {
		edgeAllocated, svcEdge, _ = service.ReducedMode.ServiceAllocate(o.Domains[domainID].ActiveNodes[edgeNodeName], edgeLoc, eventID)
		if edgeAllocated {
			break
		}
	}

	sortedNodes, _ = o.sortNodes(o.Cloud.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	for _, cloudNodeName := range sortedNodes {
		cloudAllocated, svcCloud, _ = service.ReducedMode.ServiceAllocate(o.Cloud.ActiveNodes[cloudNodeName], cloudLoc, eventID)
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
		ReducedQoS:               svcEdge.ReducedQoS,
	}
	svcEdge = nil
	svcCloud = nil
	return edgeAllocated, cloudAllocated, newSvc, nil
}

func (o *Orchestrator) edgePowerOffNode(domainID DomainID, nodeName NodeName) bool {
	cores := CreateNodeCores(len(o.Domains[domainID].ActiveNodes[nodeName].Cores))
	o.Domains[domainID].InactiveNodes[nodeName] = NewNode(cores, o.Domains[domainID].ActiveNodes[nodeName].ReallocHeuristic, nodeName)
	o.Cost = o.Cost - EdgeNodeCost
	delete(o.Domains[domainID].ActiveNodes, nodeName)

	return true
}

func (o *Orchestrator) cloudPowerOffNode(nodeName NodeName) bool {
	o.Cloud.InactiveNodes[nodeName] = NewNode(o.Cloud.ActiveNodes[nodeName].Cores, o.Cloud.ActiveNodes[nodeName].ReallocHeuristic, nodeName)
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
		o.Domains[domainID].ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName)
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
		o.Cloud.ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName)
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
	o.RunningServices[eventID] = NewRunningService(service, eventID)
	fmt.Println("added running service: ", o.RunningServices[eventID])

	fmt.Println("orchestraator domains:", o.Domains)
	fmt.Println("dimain id:", domainID)
	fmt.Println("domain:", domain)
	fmt.Println("domain active nodes:", domain.AllNodes)

	sortedNodes, _ := o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		fmt.Println("node name after allocation", nodeName)
		allocated, svc, _ := o.allocateEdge(o.RunningServices[eventID], o.Domains[domainID].ActiveNodes[nodeName], eventID)
		if allocated {
			o.QoS = o.QoS + service.StandardQoS
			o.RunningServices[eventID] = svc
			fmt.Println("running services after allocation for eventID:", o.RunningServices[eventID], eventID)
			fmt.Println("node after allocation", o.Domains[domainID].ActiveNodes[nodeName])
			fmt.Println("node name after allocation", nodeName)
			fmt.Println("allocated node:", o.Domains[domainID].ActiveNodes[o.RunningServices[eventID].AllocatedNodeEdge])
			fmt.Println("allocated cores:", o.RunningServices[eventID].AllocatedCoresEdge)
			return allocated, nil
		}
	}

	fmt.Println("the service", service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	sortedNodesNoFilter, _ := o.sortNodesNoFilter(domain.ActiveNodes)
	// newBandwidth := service.StandardMode.bandwidthEdge
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
			allocated, err := o.intraNodeRealloc(o.RunningServices[eventID], node, eventID, reallocatedEventID, NewCores)
			if allocated {
				o.QoS = o.QoS + service.StandardQoS
				fmt.Println("allocated with intra node reallocation")
				fmt.Println("node average residual bandwidth after allocation:", o.Domains[domainID].ActiveNodes[nodeName].AverageResidualBandwidth, "total residual bandwidth:", o.Domains[domainID].ActiveNodes[nodeName].TotalResidualBandwidth)
				return allocated, nil
			}
			if err != nil {
				fmt.Println("Error in intra node reallocation: ", err)
			}

			allocated, err = o.intraDomainRealloc(o.RunningServices[eventID], node, domain, sortedNodesNoFilter, eventID, reallocatedEventID)
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

	edgeAllocated, cloudAllocated, svc, _ := o.SplitSched(o.RunningServices[eventID], domainID, eventID)
	if edgeAllocated && cloudAllocated {
		fmt.Println("show svc in split scheduling", svc)
		o.RunningServices[eventID] = svc
		return true, nil
	} else {
		fmt.Println("the split scheduling didn't work. powering on some nodes")
		if !edgeAllocated {
			success, nodeName := o.edgePowerOnNode(domainID)
			fmt.Println("Edge node powered on, node name: ", nodeName)
			if success {
				allocated, svc, _ := o.allocateEdge(o.RunningServices[eventID], o.Domains[domainID].ActiveNodes[nodeName], eventID)
				if allocated {
					o.QoS = o.QoS + service.StandardQoS
					o.RunningServices[eventID] = svc
					return allocated, nil
				}
			}
		}

		if !cloudAllocated {
			o.cloudPowerOnNode()
		}
	}
	edgeAllocated, cloudAllocated, svc, _ = o.SplitSched(o.RunningServices[eventID], domainID, eventID)
	if edgeAllocated && cloudAllocated {
		allocated = true
		o.QoS = o.QoS + service.ReducedQoS
		fmt.Println("show svc in split scheduling", svc)
		o.RunningServices[eventID] = svc
	} else {
		return false, nil
	}

	return allocated, nil
}

func (o *Orchestrator) Deallocate(domainID DomainID, serviceID ServiceID, eventID ServiceID) bool {
	domain := o.Domains[domainID]
	service := o.RunningServices[eventID]
	allocatedMode := o.RunningServices[eventID].AllocationMode
	fmt.Println("Deallocating service: ", eventID)
	fmt.Println("the running services before deallocation: ", o.RunningServices[eventID])
	fmt.Println("Allocated mode: ", allocatedMode)
	fmt.Println("domain id", domainID)
	if allocatedMode == StandardMode {
		nodeN := service.AllocatedNodeEdge
		fmt.Println("node name:", nodeN)
		fmt.Println("node before deallocation", domain.ActiveNodes)
		node := domain.ActiveNodes[nodeN]
		cores := o.RunningServices[eventID].AllocatedCoresEdge
		fmt.Println("node after deallocation", node)
		fmt.Println("node", node)
		fmt.Println("cores", cores)
		service.StandardMode.ServiceDeallocate(eventID, node)
		// o.QoS = o.QoS - service.StandardQoS
	}
	if allocatedMode == ReducedMode {
		edgeNode := domain.ActiveNodes[o.RunningServices[eventID].AllocatedNodeEdge]
		cloudNode := o.Cloud.ActiveNodes[o.RunningServices[eventID].AllocatedNodeCloud]
		service.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
		service.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
		// o.QoS = o.QoS - service.ReducedQoS
	}

	for nodeName, node := range domain.ActiveNodes {
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.edgePowerOffNode(domainID, nodeName)
		}
	}
	for nodeName, node := range o.Cloud.ActiveNodes {
		if node.AverageResidualBandwidth == 0 && node.TotalResidualBandwidth == 0 {
			o.cloudPowerOffNode(nodeName)
		}
	}

	delete(o.RunningServices, eventID)

	return true
}
