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

func (o *Orchestrator) allocateEdge(service *Service, node *Node, eventID ServiceID) (bool, error) {
	fmt.Println("Allocating standard service: ", service.serviceID, " to node: ", node.NodeName)
	allocated, svc, err := service.StandardMode.ServiceAllocate(service, node, eventID)
	o.RunningServices[eventID] = svc
	fmt.Println("service:", svc)
	fmt.Println("Allocated? ", allocated)
	return allocated, err
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
		edgeAllocated, svcEdge, _ = service.ReducedMode.ServiceAllocate(service, o.Domains[domainID].ActiveNodes[edgeNodeName], edgeLoc, eventID)
		if edgeAllocated {
			break
		}
	}

	sortedNodes, _ = o.sortNodes(o.Cloud.ActiveNodes, service.ReducedMode.cpusCloud, service.ReducedMode.bandwidthCloud)
	for _, cloudNodeName := range sortedNodes {
		cloudAllocated, svcCloud, _ = service.ReducedMode.ServiceAllocate(service, o.Cloud.ActiveNodes[cloudNodeName], cloudLoc, eventID)
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
	fmt.Println("added running service: ", o.RunningServices[eventID])

	sortedNodes, _ := o.sortNodes(domain.ActiveNodes, service.StandardMode.cpusEdge, service.StandardMode.bandwidthEdge)
	for _, nodeName := range sortedNodes {
		fmt.Println("node name after allocation", nodeName)
		allocated, _ := o.allocateEdge(service, o.Domains[domainID].ActiveNodes[nodeName], eventID)
		if allocated {
			fmt.Println("qos before allocation", o.QoS)
			fmt.Println("qos of the service", service.StandardQoS)
			o.QoS = o.QoS + service.StandardQoS
			// o.RunningServices[eventID] = svc
			fmt.Println("running services after allocation for eventID:", o.RunningServices[eventID], eventID)
			fmt.Println("node after allocation", o.Domains[domainID].ActiveNodes[nodeName])
			fmt.Println("node name after allocation", nodeName)
			fmt.Println("allocated node:", o.Domains[domainID].ActiveNodes[o.RunningServices[eventID].AllocatedNodeEdge])
			fmt.Println("allocated cores:", o.RunningServices[eventID].AllocatedCoresEdge)
			return allocated, nil
		}
	}

	edgeAllocated, cloudAllocated, svc, _ := o.SplitSched(service, domainID, eventID)
	if edgeAllocated && cloudAllocated {
		fmt.Println("show svc in split scheduling", svc)
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
		fmt.Println("qos before split scheduling", o.QoS)
		fmt.Println("qos of the service", service.ReducedQoS)
		o.QoS = o.QoS + service.ReducedQoS
		fmt.Println("show svc in split scheduling", svc)
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
	var serviceQoS QoS
	fmt.Println(service.StandardMode)
	fmt.Println(service.ReducedMode)

	if allocatedMode == StandardMode {
		serviceQoS = service.StandardQoS
		fmt.Println("qos before deallocation", o.QoS)
		fmt.Println("service standard qos", service.StandardQoS)
		nodeN := service.AllocatedNodeEdge
		fmt.Println("node name:", nodeN)
		fmt.Println("node before deallocation", domain.ActiveNodes)
		node := domain.ActiveNodes[nodeN]
		cores := service.AllocatedCoresEdge
		fmt.Println("node after deallocation", node)
		fmt.Println("node", node)
		fmt.Println("cores", cores)
		_, err := service.StandardMode.ServiceDeallocate(eventID, node)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
		}

	}
	if allocatedMode == ReducedMode {
		serviceQoS = service.ReducedQoS
		fmt.Println("qos before deallocation", o.QoS)
		fmt.Println("service standard qos", service.ReducedQoS)
		edgeNode := domain.ActiveNodes[service.AllocatedNodeEdge]
		cloudNode := o.Cloud.ActiveNodes[service.AllocatedNodeCloud]
		_, err := service.ReducedMode.ServiceDeallocate(eventID, edgeNode, edgeLoc)
		_, err = service.ReducedMode.ServiceDeallocate(eventID, cloudNode, cloudLoc)
		if err != nil {
			fmt.Println("Error in deallocation: ", err)
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
	o.QoS = o.QoS - serviceQoS

	return true
}
