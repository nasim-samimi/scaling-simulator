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

func (s *standardSched) ServiceAllocate(node *Node) (bool, error) {
	allocated := false

	allocatedCores, err := node.NodeAllocate(s.cpusEdge, s.bandwidthEdge, s.service)
	if err != nil {
		fmt.Println("Error allocating service: ", err)
		return allocated, err
	}
	s.service.allocatedNodeEdge = node.NodeName
	allocated = true
	s.service.allocatedCoresEdge = allocatedCores
	s.service.allocationMode = StandardMode

	s.service.AverageResidualBandwidth = s.bandwidthEdge
	s.service.TotalResidualBandwidth = s.bandwidthEdge * float64(s.cpusEdge)
	s.service.allocatedDomain = node.DomainID
	return allocated, err
}

func (s *standardSched) ServiceDeallocate(service *Service, node *Node) (*Service, error) {
	node.NodeDeallocate(service.serviceID)
	return service, nil
}

func (r *reducedSched) ServiceAllocate(node *Node, loc Location) (bool, error) {
	allocated := false

	switch loc {
	case edgeLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusEdge, r.bandwidthEdge, r.service)

		if err != nil {
			return allocated, err
		}
		r.service.allocatedNodeEdge = node.NodeName
		allocated = true
		r.service.allocatedCoresEdge = allocatedCores
		r.service.allocationMode = ReducedMode

	case cloudLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, r.service)
		if err != nil {
			return false, err
		}
		r.service.allocatedNodeCloud = node.NodeName
		r.service.allocatedCoresCloud = allocatedCores
		allocated = true
	}
	r.service.AverageResidualBandwidth = (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / (float64(r.cpusEdge) + float64(r.cpusCloud))
	r.service.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)

	return true, nil

}

func (r *reducedSched) ServiceDeallocate(service *Service, node *Node, location Location) (*Service, error) {
	switch location {
	case edgeLoc:
		node.NodeDeallocate(service.serviceID)
	case cloudLoc:
		node.NodeDeallocate(service.serviceID)
	}
	return service, nil
}

type ServiceID string
type ServiceMode string
type Services map[ServiceID]*Service

const (
	ReducedMode  ServiceMode = "Reduced"
	StandardMode ServiceMode = "Standard"
)

type Service struct {
	importanceFactor         float64
	serviceID                ServiceID
	reducedMode              *reducedSched
	standardMode             *standardSched
	allocatedCoresEdge       []CoreID
	allocatedCoresCloud      []CoreID
	allocatedNodeEdge        NodeName
	allocatedNodeCloud       NodeName
	allocatedDomain          DomainID
	allocationMode           ServiceMode
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
		importanceFactor: importanceFactor,
		serviceID:        serviceID,
		reducedMode:      reduced,
		standardMode:     standard,
		// serviceModel:        serviceModel,
	}
	standardQoS := standard.bandwidthEdge * float64(standard.cpusEdge)
	reducedQoS := reduced.bandwidthEdge*float64(reduced.cpusEdge) + reduced.bandwidthCloud*float64(reduced.cpusCloud)
	service.StandardQoS = QoS(standardQoS)
	service.ReducedQoS = QoS(reducedQoS)
	service.reducedMode.service = service
	service.standardMode.service = service
	return service
}

func (t *Service) ServiceReallocate(service *Service, node *Node, domain *Domain, location Location) bool {
	return true
}
