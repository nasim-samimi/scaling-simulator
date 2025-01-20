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
		Domains:                domains,
		PartitionHeuristic:     partitionHeuristic,
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
	// newSvc := &Service{
	// 	StandardMode:             svc.StandardMode,
	// 	ReducedMode:              svc.ReducedMode,
	// 	ImportanceFactor:         svc.ImportanceFactor,
	// 	serviceID:                svc.serviceID,
	// 	AllocatedCoresEdge:       svc.AllocatedCoresEdge,
	// 	AllocatedCoresCloud:      svc.AllocatedCoresCloud,
	// 	AllocatedNodeEdge:        svc.AllocatedNodeEdge,
	// 	AllocatedNodeCloud:       svc.AllocatedNodeCloud,
	// 	AllocatedDomain:          svc.AllocatedDomain,
	// 	AllocationMode:           svc.AllocationMode,
	// 	AverageResidualBandwidth: svc.AverageResidualBandwidth,
	// 	TotalResidualBandwidth:   svc.TotalResidualBandwidth,
	// 	StandardQoS:              svc.StandardQoS,
	// }
	// node.AllocatedServices[eventID] = newSvc
	return allocated, svc, err
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
		fmt.Println("Node: ", node.NodeName, " Average Residual Bandwidth: ", node.AverageResidualBandwidth, "total residual bandwidth: ", node.TotalResidualBandwidth)
	}
	// nodeNames := make([]NodeName, 0, len(sortedNodes))
	// for nodeN := range sortedNodes {
	// 	nodeNames = append(nodeNames, nodeN)
	// }
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
	// for i, node := range sortedNodes {
	// 	sortedNodeNames[i] = node.NodeName
	// }
	return sortedNodeNames, nil
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
	for nodeName, node := range o.Cloud.InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Cloud.ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName)
		o.Cost = o.Cost + CloudNodeCost
		delete(o.Cloud.InactiveNodes, nodeName)
		break
	}

	return true
}

func (o *Orchestrator) Allocate(domainID DomainID, serviceID ServiceID, eventID ServiceID) (bool, error) {
	allocated := false
	domain := o.Domains[domainID]
	service := o.AllServices[serviceID]
	o.RunningServices[eventID] = NewRunningService(service, eventID)
	fmt.Println("added running service: ", o.RunningServices[eventID])

	sortedNodes, _ := o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		fmt.Println("node name after allocation", nodeName)
		allocated, svc, _ := o.allocateEdge(o.RunningServices[eventID], o.Domains[domainID].ActiveNodes[nodeName], eventID)
		if allocated {
			o.QoS = o.QoS + service.StandardQoS
			o.RunningServices[eventID] = svc
			return allocated, nil
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
	if allocatedMode == StandardMode {
		node := domain.ActiveNodes[o.RunningServices[eventID].AllocatedNodeEdge]
		service.StandardMode.ServiceDeallocate(eventID, node)
	}
	if allocatedMode == ReducedMode {
		edgeNode := domain.ActiveNodes[o.RunningServices[eventID].AllocatedNodeEdge]
		cloudNode := o.Cloud.ActiveNodes[o.RunningServices[eventID].AllocatedNodeCloud]
		service.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
		service.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
	}
	delete(o.RunningServices, eventID)

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

	return true
}
