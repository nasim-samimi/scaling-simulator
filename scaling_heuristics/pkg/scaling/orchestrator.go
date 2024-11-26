package scaling

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
type Orchestrator struct {
	NodeSelectionHeuristic NodeSelectionHeuristic
	ReallocationHeuristic  ReallocationHeuristic
	domains                Domains
	cloud                  *Cloud
	AllServices            Services
	RunningServices        Services // change name of service to service
	Cost                   int
	QoS                    int
}

func NewOrchestrator(nodeSelectionHeuristic NodeSelectionHeuristic, reallocationHeuristic ReallocationHeuristic, cloud *Cloud, domains Domains, services Services) *Orchestrator {
	return &Orchestrator{
		NodeSelectionHeuristic: nodeSelectionHeuristic,
		ReallocationHeuristic:  reallocationHeuristic,
		domains:                domains,
		cloud:                  cloud,
		AllServices:            services,
		Cost:                   0,
		QoS:                    0,
	}
}

func (o *Orchestrator) SelectNode(service *Service) *Node {
	var node *Node

	return node
}

func (o *Orchestrator) allocateEdge(service *Service, node *Node) (bool, error) {
	allocated, err := service.standardMode.ServiceAllocate(node)
	return allocated, err
}

func (o *Orchestrator) getReallocatedService(node *Node, t *Service) (ServiceID, error) {
	var selectedServiceID ServiceID
	switch o.ReallocationHeuristic {
	case HBI:
		hbi := 0.0
		for _, service := range node.AllocatedServices {
			if service.allocationMode == StandardMode {
				bi := float64(service.importanceFactor) * service.standardMode.bandwidthEdge
				if bi > hbi {
					selectedServiceID = service.serviceID
					hbi = bi
				}
			}
		}
	case HCI:
		hci := float64(0)

		for _, service := range node.AllocatedServices {
			if service.allocationMode == StandardMode {
				ci := service.importanceFactor * float64(service.standardMode.cpusEdge)
				if ci > hci {
					selectedServiceID = service.serviceID
					hci = ci
				}
			}
		}
	case HBCI:
		hbci := 0.0
		for _, service := range node.AllocatedServices {
			if service.allocationMode == StandardMode {
				bci := float64(service.importanceFactor) * service.standardMode.bandwidthEdge * float64(service.standardMode.cpusEdge)
				if bci > hbci {
					selectedServiceID = service.serviceID
					hbci = bci
				}
			}
		}
	case HBIcC:
		hbic := 0.0
		for _, service := range node.AllocatedServices {
			if service.allocationMode == StandardMode {
				if service.standardMode.cpusEdge > t.standardMode.cpusEdge {
					bic := float64(service.importanceFactor) * service.standardMode.bandwidthEdge * float64(service.standardMode.cpusEdge)
					if bic > hbic {
						selectedServiceID = service.serviceID
						hbic = bic
					}
				}
			}
		}

	default: // HBI
		hbi := 0.0
		for _, service := range node.AllocatedServices {
			if service.allocationMode == StandardMode {
				bi := float64(service.importanceFactor) * service.standardMode.bandwidthEdge
				if bi > hbi {
					selectedServiceID = service.serviceID
					hbi = bi
				}
			}
		}
	}
	return selectedServiceID, nil
}

func (o *Orchestrator) sortNodes(nodes Nodes) ([]NodeName, error) {
	// sort nodes according to the heuristic
	// switch o.NodeSelectionHeuristic {
	// case MinMin:
	// 	for _, node := range nodes {

	// 		for _, cores := range node.Cores {
	// 			// sort cores
	// 		}
	// 	}
	// case MaxMax:
	// 	for _, node := range nodes {
	// 		for _, cores := range node.Cores {
	// 			// sort cores
	// 		}
	// 	}
	// }
	return nil, nil
}

func (o *Orchestrator) intraNodeRealloc(service *Service, node *Node) (bool, error) {
	reallocatedService := &Service{}
	reallocatedServiceID, err := o.getReallocatedService(node, service)
	if err != nil {
		return false, err
	}

	reallocation, err := node.IntraNodeReallocateTest(service, reallocatedServiceID)
	if err != nil {
		return false, err
	}
	if reallocation {
		return true, nil
	}

	reallocatedService, err = service.standardMode.ServiceDeallocate(node.AllocatedServices[reallocatedServiceID], node)
	if err != nil {
		return false, err
	}

	service.standardMode.ServiceAllocate(node)
	reallocatedService.standardMode.ServiceAllocate(node)

	return true, nil
}

func (o *Orchestrator) intraDomainRealloc(service *Service, node *Node, domain *Domain) (bool, error) {
	searchingNodes := domain.OnNodes
	delete(searchingNodes, node.NodeName)
	sortedNodes, _ := o.sortNodes(searchingNodes)
	otherServiceID, _ := o.getReallocatedService(node, service)
	reallocated := false

	otherService := node.AllocatedServices[otherServiceID]
	reallocation, _ := node.IntraDomainReallocateTest(service, otherServiceID)
	if reallocation {
		for _, nodeName := range sortedNodes {
			otherNode := domain.OnNodes[nodeName]

			allocatedCore, _ := otherNode.NodeAdmission.Admission(service.standardMode.cpusEdge, service.standardMode.bandwidthEdge, otherNode.Cores)
			if allocatedCore != nil {

				otherService.standardMode.ServiceDeallocate(service, node)
				service.standardMode.ServiceAllocate(node)
				otherService.standardMode.ServiceAllocate(otherNode)
				reallocated = true
			}
			break
		}

	}
	return reallocated, nil
}

func (o *Orchestrator) SplitSched(service *Service, domain *Domain) (bool, bool, error) {
	// edge-cloud split (has qos degradation) -- there is no cloud only apparently
	sortedNodes, _ := o.sortNodes(domain.OnNodes)
	edgeAllocated := false
	cloudAllocated := false
	for _, edgeNodeName := range sortedNodes {
		node := domain.OnNodes[edgeNodeName]
		edgeAllocated, _ = service.reducedMode.ServiceAllocate(node, edgeLoc)
		if edgeAllocated {
			break
		}
	}

	sortedNodes, _ = o.sortNodes(o.cloud.OnNodes)
	for _, cloudNodeName := range sortedNodes {
		node := o.cloud.OnNodes[cloudNodeName]
		cloudAllocated, _ = service.reducedMode.ServiceAllocate(node, cloudLoc)
		if cloudAllocated {
			break
		}
	}
	return edgeAllocated, cloudAllocated, nil
}

func (o *Orchestrator) PowerOffNode() bool {
	return true
}

func (o *Orchestrator) edgePowerOnNode() bool {
	// add a node id and initialise it
	// NewNode()
	return true
}
func (o *Orchestrator) cloudPowerOnNode() bool {
	// add a node id and initialise it
	// NewNode()
	return true
}

func (o *Orchestrator) Allocate(domainID DomainID, serviceID ServiceID) (bool, QoS, Cost, error) {
	// TODO: calculation of qos per cost is not implemented
	allocated := false
	domain := o.domains[domainID]
	service := o.AllServices[serviceID]
	// sort nodes once here
	sortedNodes, _ := o.sortNodes(domain.OnNodes)
	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ = o.allocateEdge(service, node)
		if allocated {
			return allocated, 0, 0, nil
		}
	}

	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ := o.intraNodeRealloc(service, node)
		if allocated {
			return allocated, 0, 0, nil
		}
	}
	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ := o.intraDomainRealloc(service, node, domain)
		if allocated {
			return allocated, 0, 0, nil
		}
	}
	edgeAllocated, cloudAllocated, _ := o.SplitSched(service, domain)
	if edgeAllocated && cloudAllocated {
		allocated = true
	} else {
		if !edgeAllocated {
			o.edgePowerOnNode()
		}
		if !cloudAllocated {
			o.cloudPowerOnNode()
		}
	}

	return allocated, 0, 0, nil
}

func (o *Orchestrator) Deallocate() bool {
	return true
}
