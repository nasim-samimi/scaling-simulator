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

func (s *standardSched) ServiceAllocate(service *Service, node *Node, eventID ServiceID) (bool, *Service, error) {
	allocated := false
	fmt.Println("Allocating service: ", service)
	allocatedCores, err := node.NodeAllocate(s.cpusEdge, s.bandwidthEdge, service, eventID)
	if err != nil {
		fmt.Println("Error allocating service: ", err)
		return allocated, service, err
	}
	fmt.Println("node name:", node.NodeName)
	service.AllocatedNodeEdge = node.NodeName
	allocated = true
	service.AllocatedCoresEdge = allocatedCores
	service.AllocationMode = StandardMode

	service.AverageResidualBandwidth = s.bandwidthEdge
	service.TotalResidualBandwidth = s.bandwidthEdge * float64(s.cpusEdge)
	service.AllocatedDomain = node.DomainID
	newSvc := &Service{
		StandardMode:             service.StandardMode,
		ReducedMode:              service.ReducedMode,
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                service.serviceID,
		AllocatedCoresEdge:       service.AllocatedCoresEdge,
		AllocatedCoresCloud:      service.AllocatedCoresCloud,
		AllocatedNodeEdge:        service.AllocatedNodeEdge,
		AllocatedNodeCloud:       service.AllocatedNodeCloud,
		AllocatedDomain:          service.AllocatedDomain,
		AllocationMode:           service.AllocationMode,
		AverageResidualBandwidth: service.AverageResidualBandwidth,
		TotalResidualBandwidth:   service.TotalResidualBandwidth,
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
	}
	node.AllocatedServices[eventID] = &Service{
		ImportanceFactor:         service.ImportanceFactor,
		serviceID:                eventID,
		ReducedMode:              service.ReducedMode,
		StandardMode:             service.StandardMode,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           StandardMode,
		AverageResidualBandwidth: service.AverageResidualBandwidth,
		TotalResidualBandwidth:   service.TotalResidualBandwidth,
		StandardQoS:              service.StandardQoS,
		ReducedQoS:               service.ReducedQoS,
	}
	fmt.Println("Allocating service after: ", service)
	return allocated, newSvc, err
}

func (s *standardSched) ServiceDeallocate(eventID ServiceID, node *Node) (ServiceID, error) {
	node.NodeDeallocate(eventID)
	return eventID, nil
}

func (r *reducedSched) ServiceAllocate(service *Service, node *Node, loc Location, eventID ServiceID) (bool, *Service, error) {
	allocated := false

	switch loc {
	case edgeLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusEdge, r.bandwidthEdge, service, eventID)

		if err != nil {
			return allocated, service, err
		}
		service.AllocatedNodeEdge = node.NodeName
		allocated = true
		service.AllocatedCoresEdge = allocatedCores
		service.AllocationMode = ReducedMode
		service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)
		service.AverageResidualBandwidth = service.TotalResidualBandwidth / float64(r.cpusEdge+r.cpusCloud)
		newSvc := &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresEdge:       allocatedCores,
			AllocatedNodeEdge:        node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: service.AverageResidualBandwidth,
			TotalResidualBandwidth:   service.TotalResidualBandwidth,
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
			AverageResidualBandwidth: service.AverageResidualBandwidth,
			TotalResidualBandwidth:   service.TotalResidualBandwidth,
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		return allocated, newSvc, nil

	case cloudLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, service, eventID)
		if err != nil {
			return false, service, err
		}
		service.AllocatedNodeCloud = node.NodeName
		service.AllocatedCoresCloud = allocatedCores
		allocated = true
		service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)
		service.AverageResidualBandwidth = service.TotalResidualBandwidth / float64(r.cpusEdge+r.cpusCloud)
		newSvc := &Service{
			ImportanceFactor:         service.ImportanceFactor,
			serviceID:                service.serviceID,
			ReducedMode:              service.ReducedMode,
			StandardMode:             service.StandardMode,
			AllocatedCoresCloud:      allocatedCores,
			AllocatedNodeCloud:       node.NodeName,
			AllocatedDomain:          service.AllocatedDomain,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: service.AverageResidualBandwidth,
			TotalResidualBandwidth:   service.TotalResidualBandwidth,
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
			AllocatedDomain:          service.AllocatedDomain,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: service.AverageResidualBandwidth,
			TotalResidualBandwidth:   service.TotalResidualBandwidth,
			StandardQoS:              service.StandardQoS,
			ReducedQoS:               service.ReducedQoS,
		}
		return allocated, newSvc, nil
	}
	// r.service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)
	// r.service.AverageResidualBandwidth = r.service.TotalResidualBandwidth / float64(r.cpusEdge+r.cpusCloud)

	// newSvc := &Service{
	// 	StandardMode:             r.service.StandardMode,
	// 	ReducedMode:              r.service.ReducedMode,
	// 	ImportanceFactor:         r.service.ImportanceFactor,
	// 	serviceID:                r.service.serviceID,
	// 	AllocatedCoresEdge:       r.service.AllocatedCoresEdge,
	// 	AllocatedCoresCloud:      r.service.AllocatedCoresCloud,
	// 	AllocatedNodeEdge:        r.service.AllocatedNodeEdge,
	// 	AllocatedNodeCloud:       r.service.AllocatedNodeCloud,
	// 	AllocatedDomain:          r.service.AllocatedDomain,
	// 	AllocationMode:           r.service.AllocationMode,
	// 	AverageResidualBandwidth: r.service.AverageResidualBandwidth,
	// 	TotalResidualBandwidth:   r.service.TotalResidualBandwidth,
	// 	StandardQoS:              r.service.StandardQoS,
	// }

	return true, &Service{}, nil

}

func (r *reducedSched) ServiceDeallocate(eventID ServiceID, node *Node, location Location) (ServiceID, error) {
	switch location {
	case edgeLoc:
		node.NodeDeallocate(eventID)
	case cloudLoc:
		node.NodeDeallocate(eventID)
	}
	return eventID, nil
}

type ServiceID string
type ServiceMode string
type Services map[ServiceID]*Service

const (
	ReducedMode  ServiceMode = "Reduced"
	StandardMode ServiceMode = "Standard"
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
	}

	return service
}

func (t *Service) ServiceReallocate(service *Service, node *Node, domain *Domain, location Location) bool {
	return true
}

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
