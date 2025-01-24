package scaling

import "fmt"

type reducedSched struct {
	bandwidthEdge  float64
	cpusEdge       uint64
	bandwidthCloud float64
	cpusCloud      uint64
}

type standardSched struct {
	bandwidthEdge float64
	cpusEdge      uint64
}

func (s *standardSched) ServiceAllocate(service *Service, node *Node, eventID ServiceID, cpuThreshold float64) (bool, *Service, error) {
	allocated := false
	fmt.Println("Allocating service: ", service)
	allocatedCores, err := node.NodeAllocate(s.cpusEdge, s.bandwidthEdge, service, eventID, cpuThreshold)
	if err != nil {
		fmt.Println("Error allocating service: ", err)
		return allocated, service, err
	}
	fmt.Println("node name:", node.NodeName)
	allocated = true

	newSvc := &Service{
		StandardMode:             service.StandardMode,
		ReducedMode:              service.ReducedMode,
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                service.serviceID,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedCoresCloud:      service.AllocatedCoresCloud,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedNodeCloud:       service.AllocatedNodeCloud,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           StandardMode,
		AverageResidualBandwidth: s.bandwidthEdge,
		TotalResidualBandwidth:   s.bandwidthEdge * float64(s.cpusEdge),
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
	}
	node.AllocatedServices[eventID] = &Service{
		StandardMode:             service.StandardMode,
		ReducedMode:              service.ReducedMode,
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                service.serviceID,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedCoresCloud:      service.AllocatedCoresCloud,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedNodeCloud:       service.AllocatedNodeCloud,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           StandardMode,
		AverageResidualBandwidth: s.bandwidthEdge,
		TotalResidualBandwidth:   s.bandwidthEdge * float64(s.cpusEdge),
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
	}
	fmt.Println("Allocating service after: ", service)
	return allocated, newSvc, err
}

func (s *standardSched) ServiceDeallocate(eventID ServiceID, node *Node) (ServiceID, error) {
	deallocated := node.NodeDeallocate(eventID)
	if deallocated {
		return eventID, nil
	}
	return eventID, fmt.Errorf("Service not deallocated")
}

func (r *reducedSched) ServiceAllocate(service *Service, node *Node, loc Location, eventID ServiceID, cpuThreshold float64) (bool, *Service, error) {
	allocated := false

	switch loc {
	case edgeLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusEdge, r.bandwidthEdge, service, eventID, cpuThreshold)

		if err != nil {
			return allocated, service, err
		}

		allocated = true

		newSvc := &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresEdge:       allocatedCores,
			AllocatedNodeEdge:        node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / float64(r.cpusEdge+r.cpusCloud),
			TotalResidualBandwidth:   r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud),
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		node.AllocatedServices[eventID] = &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresEdge:       allocatedCores,
			AllocatedNodeEdge:        node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / float64(r.cpusEdge+r.cpusCloud),
			TotalResidualBandwidth:   r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud),
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		return allocated, newSvc, nil

	case cloudLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, service, eventID, cpuThreshold)
		if err != nil {
			return false, service, err
		}
		allocated = true

		newSvc := &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresCloud:      allocatedCores,
			AllocatedNodeCloud:       node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / float64(r.cpusEdge+r.cpusCloud),
			TotalResidualBandwidth:   r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud),
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		node.AllocatedServices[eventID] = &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresCloud:      allocatedCores,
			AllocatedNodeCloud:       node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / float64(r.cpusEdge+r.cpusCloud),
			TotalResidualBandwidth:   r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud),
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		return allocated, newSvc, nil
	}

	return true, &Service{}, nil

}

func (r *reducedSched) ServiceDeallocate(eventID ServiceID, node *Node, location Location) (ServiceID, error) {
	var deallocated bool
	switch location {
	case edgeLoc:
		deallocated = node.NodeDeallocate(eventID)
	case cloudLoc:
		deallocated = node.CloudNodeDeallocate(eventID)
	}
	if !deallocated {
		return eventID, fmt.Errorf("Service not deallocated")
	}
	return eventID, nil
}

type ServiceID string
type ServiceMode string
type Services map[ServiceID]*Service

const (
	ReducedMode     ServiceMode = "Reduced"
	StandardMode    ServiceMode = "Standard"
	EdgeReducedMode ServiceMode = "EdgeReduced"
)

type Service struct {
	ImportanceFactor         float64
	serviceID                ServiceID
	ReducedMode              *reducedSched
	StandardMode             *standardSched
	AllocatedCoresEdge       []CoreID
	AllocatedCoresCloud      []CoreID
	AllocatedNodeEdge        NodeName
	AllocatedNodeCloud       NodeName
	AllocatedDomain          DomainID
	AllocationMode           ServiceMode
	AverageResidualBandwidth float64
	TotalResidualBandwidth   float64
	StandardQoS              QoS
	ReducedQoS               QoS
	EdgeReducedQoS           QoS
	// serviceModel           serviceModel
}

func NewService(importanceFactor float64, serviceID ServiceID, standardBandwidth float64, standardCores uint64, reducedEdgeBandwidth float64, reducedEdgeCores uint64, reducedCloudBandwidth float64, reducedCloudCores uint64) *Service {
	standard := &standardSched{
		bandwidthEdge: standardBandwidth,
		cpusEdge:      standardCores,
	}
	reduced := &reducedSched{
		bandwidthEdge:  reducedEdgeBandwidth,
		cpusEdge:       reducedEdgeCores,
		bandwidthCloud: reducedCloudBandwidth,
		cpusCloud:      reducedCloudCores,
	}
	service := &Service{
		ImportanceFactor: importanceFactor,
		serviceID:        serviceID,
		ReducedMode:      reduced,
		StandardMode:     standard,
		StandardQoS:      QoS(standard.bandwidthEdge * float64(standard.cpusEdge) * importanceFactor),
		ReducedQoS:       QoS((reduced.bandwidthEdge*float64(reduced.cpusEdge) + reduced.bandwidthCloud*float64(reduced.cpusCloud)) * importanceFactor),
		EdgeReducedQoS:   QoS(reduced.bandwidthCloud * float64(reduced.cpusCloud) * importanceFactor),
	}

	return service
}

func (r *reducedSched) ServiceUpgraded(event *Service, node *Node, location Location) error {
	switch location {
	case edgeLoc:
		node.Upgraded(event)
	case cloudLoc:
		node.Upgraded(event)
	}
	return nil
}

// func (r *standardSched) ServiceDegraded(event Service, node *Node) error {

// 	node.Upgraded(event)

// 	return nil
// }

func NewRunningService(service *Service, eventID ServiceID) *Service {
	return &Service{
		ImportanceFactor: service.ImportanceFactor,
		serviceID:        service.serviceID,
		ReducedMode:      service.ReducedMode,
		StandardMode:     service.StandardMode,
		StandardQoS:      service.StandardQoS,
		ReducedQoS:       service.ReducedQoS,
	}
}

func (r *reducedSched) EdgeServiceAllocate(service *Service, node *Node, eventID ServiceID, cpuThreshold float64) (bool, *Service, error) {
	allocated := false

	allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, service, eventID, cpuThreshold)

	if err != nil {
		return allocated, service, err
	}

	allocated = true

	newSvc := &Service{
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                service.serviceID,
		ReducedMode:              service.ReducedMode,
		StandardMode:             service.StandardMode,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           EdgeReducedMode,
		AverageResidualBandwidth: (r.bandwidthCloud * float64(r.cpusCloud)) / float64(r.cpusCloud),
		TotalResidualBandwidth:   r.bandwidthCloud * float64(r.cpusCloud),
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
		EdgeReducedQoS:           service.EdgeReducedQoS,
	}
	node.AllocatedServices[eventID] = &Service{
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                service.serviceID,
		ReducedMode:              service.ReducedMode,
		StandardMode:             service.StandardMode,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           EdgeReducedMode,
		AverageResidualBandwidth: (r.bandwidthCloud * float64(r.cpusCloud)) / float64(r.cpusCloud),
		TotalResidualBandwidth:   r.bandwidthCloud * float64(r.cpusCloud),
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
		EdgeReducedQoS:           service.EdgeReducedQoS,
	}
	return allocated, newSvc, nil

}

func (r *reducedSched) EdgeServiceDeallocate(eventID ServiceID, node *Node) (ServiceID, error) {

	node.CloudNodeDeallocate(eventID)

	return eventID, nil
}
