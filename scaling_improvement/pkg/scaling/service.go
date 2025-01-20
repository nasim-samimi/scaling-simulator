package scaling

import "fmt"

type reducedSched struct {
	bandwidthEdge  float64
	cpusEdge       uint64
	bandwidthCloud float64
	cpusCloud      uint64
	service        *Service
}

type standardSched struct {
	bandwidthEdge float64
	cpusEdge      uint64
	service       *Service
}

func (s *standardSched) ServiceAllocate(node *Node, eventID ServiceID) (bool, *Service, error) {
	allocated := false
	fmt.Println("Allocating service: ", s.service)
	allocatedCores, err := node.NodeAllocate(s.cpusEdge, s.bandwidthEdge, s.service, eventID)
	if err != nil {
		fmt.Println("Error allocating service: ", err)
		return allocated, s.service, err
	}
	fmt.Println("node name:", node.NodeName)
	s.service.AllocatedNodeEdge = node.NodeName
	allocated = true
	s.service.AllocatedCoresEdge = allocatedCores
	s.service.AllocationMode = StandardMode

	s.service.AverageResidualBandwidth = s.bandwidthEdge
	s.service.TotalResidualBandwidth = s.bandwidthEdge * float64(s.cpusEdge)
	s.service.AllocatedDomain = node.DomainID
	newSvc := &Service{
		StandardMode:             s.service.StandardMode,
		ReducedMode:              s.service.ReducedMode,
		ImportanceFactor:         s.service.ImportanceFactor,
		serviceID:                s.service.serviceID,
		AllocatedCoresEdge:       s.service.AllocatedCoresEdge,
		AllocatedCoresCloud:      s.service.AllocatedCoresCloud,
		AllocatedNodeEdge:        s.service.AllocatedNodeEdge,
		AllocatedNodeCloud:       s.service.AllocatedNodeCloud,
		AllocatedDomain:          s.service.AllocatedDomain,
		AllocationMode:           s.service.AllocationMode,
		AverageResidualBandwidth: s.service.AverageResidualBandwidth,
		TotalResidualBandwidth:   s.service.TotalResidualBandwidth,
		StandardQoS:              s.service.StandardQoS,
	}
	node.AllocatedServices[eventID] = &Service{
		ImportanceFactor:         s.service.ImportanceFactor,
		serviceID:                eventID,
		ReducedMode:              s.service.ReducedMode,
		StandardMode:             s.service.StandardMode,
		AllocatedCoresEdge:       allocatedCores,
		AllocatedNodeEdge:        node.NodeName,
		AllocatedDomain:          node.DomainID,
		AllocationMode:           StandardMode,
		AverageResidualBandwidth: s.service.AverageResidualBandwidth,
		TotalResidualBandwidth:   s.service.TotalResidualBandwidth,
		StandardQoS:              s.service.StandardQoS,
		ReducedQoS:               s.service.ReducedQoS,
	}
	fmt.Println("Allocating service after: ", s.service)
	return allocated, newSvc, err
}

func (s *standardSched) ServiceDeallocate(eventID ServiceID, node *Node) (ServiceID, error) {
	node.NodeDeallocate(eventID)
	return eventID, nil
}

func (r *reducedSched) ServiceAllocate(node *Node, loc Location, eventID ServiceID) (bool, *Service, error) {
	allocated := false

	switch loc {
	case edgeLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusEdge, r.bandwidthEdge, r.service, eventID)

		if err != nil {
			return allocated, r.service, err
		}
		r.service.AllocatedNodeEdge = node.NodeName
		allocated = true
		r.service.AllocatedCoresEdge = allocatedCores
		r.service.AllocationMode = ReducedMode
		r.service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)
		r.service.AverageResidualBandwidth = r.service.TotalResidualBandwidth / float64(r.cpusEdge+r.cpusCloud)
		newSvc := &Service{
			ImportanceFactor:         r.service.ImportanceFactor,
			serviceID:                r.service.serviceID,
			ReducedMode:              r.service.ReducedMode,
			StandardMode:             r.service.StandardMode,
			AllocatedCoresEdge:       allocatedCores,
			AllocatedNodeEdge:        node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: r.service.AverageResidualBandwidth,
			TotalResidualBandwidth:   r.service.TotalResidualBandwidth,
			StandardQoS:              r.service.StandardQoS,
			ReducedQoS:               r.service.ReducedQoS,
		}
		node.AllocatedServices[eventID] = &Service{
			ImportanceFactor:         r.service.ImportanceFactor,
			serviceID:                r.service.serviceID,
			ReducedMode:              r.service.ReducedMode,
			StandardMode:             r.service.StandardMode,
			AllocatedCoresEdge:       allocatedCores,
			AllocatedNodeEdge:        node.NodeName,
			AllocatedDomain:          node.DomainID,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: r.service.AverageResidualBandwidth,
			TotalResidualBandwidth:   r.service.TotalResidualBandwidth,
			StandardQoS:              r.service.StandardQoS,
			ReducedQoS:               r.service.ReducedQoS,
		}
		return allocated, newSvc, nil

	case cloudLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, r.service, eventID)
		if err != nil {
			return false, r.service, err
		}
		r.service.AllocatedNodeCloud = node.NodeName
		r.service.AllocatedCoresCloud = allocatedCores
		allocated = true
		r.service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)
		r.service.AverageResidualBandwidth = r.service.TotalResidualBandwidth / float64(r.cpusEdge+r.cpusCloud)
		newSvc := &Service{
			ImportanceFactor:         r.service.ImportanceFactor,
			serviceID:                r.service.serviceID,
			ReducedMode:              r.service.ReducedMode,
			StandardMode:             r.service.StandardMode,
			AllocatedCoresCloud:      allocatedCores,
			AllocatedNodeCloud:       node.NodeName,
			AllocatedDomain:          r.service.AllocatedDomain,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: r.service.AverageResidualBandwidth,
			TotalResidualBandwidth:   r.service.TotalResidualBandwidth,
			StandardQoS:              r.service.StandardQoS,
			ReducedQoS:               r.service.ReducedQoS,
		}
		node.AllocatedServices[eventID] = &Service{
			ImportanceFactor:         r.service.ImportanceFactor,
			serviceID:                r.service.serviceID,
			ReducedMode:              r.service.ReducedMode,
			StandardMode:             r.service.StandardMode,
			AllocatedCoresCloud:      allocatedCores,
			AllocatedNodeCloud:       node.NodeName,
			AllocatedDomain:          r.service.AllocatedDomain,
			AllocationMode:           ReducedMode,
			AverageResidualBandwidth: r.service.AverageResidualBandwidth,
			TotalResidualBandwidth:   r.service.TotalResidualBandwidth,
			StandardQoS:              r.service.StandardQoS,
			ReducedQoS:               r.service.ReducedQoS,
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
		// serviceModel:        serviceModel,
	}
	standardQoS := standard.bandwidthEdge * float64(standard.cpusEdge)
	reducedQoS := reduced.bandwidthEdge*float64(reduced.cpusEdge) + reduced.bandwidthCloud*float64(reduced.cpusCloud)
	service.StandardQoS = QoS(standardQoS)
	service.ReducedQoS = QoS(reducedQoS)
	service.ReducedMode.service = service
	service.StandardMode.service = service
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
	}
}
