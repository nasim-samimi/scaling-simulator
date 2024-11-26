package main

import (
	"fmt"
	"log"
	"path/filepath"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
	util "github.com/nasim-samimi/scaling-simulator/pkg/util"
)

func main() {
	orchestrator := initialise()
	processEvents(orchestrator)
}

func initialise() *src.Orchestrator {
	CloudNodes := util.LoadCloudFromCSV("data/cloud.csv")
	cloud := src.NewCloud(CloudNodes)
	// read domain csv files in domain folder
	domainFilesNames, err := filepath.Glob("data/domain/*.csv")
	if err != nil {
		log.Fatal(err)
	}
	domains := make(src.Domains)
	for _, fileName := range domainFilesNames {
		domainNodes := util.LoadDomainFromCSV(fileName)
		domains[src.DomainID(fileName)] = src.NewDomain(domainNodes, src.DomainID(fileName))
	}

	svcs := util.LoadSVCFromCSV("data/svcs.csv")
	nodeHeuristic, reallocHeuristic := util.LoadHeuristicFromCSV("data/heuristics.csv")

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(src.NodeSelectionHeuristic(nodeHeuristic), src.ReallocationHeuristic(reallocHeuristic), cloud, domains, svcs)
	return orchestrator
}

func processEvents(orchestrator *src.Orchestrator) error {
	events := util.LoadEventsFromCSV("data/events.csv")
	qosPerCost := make([]int, 0)
	for _, event := range events {
		if event.EventType == "allocate" {
			allocated, qos, cost, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID)
			if err != nil {
				fmt.Println(err)
			}
			if allocated {
				qosPerCost = append(qosPerCost, int(qos)/int(cost))
			}
		}
		if event.EventType == "deallocate" {
			fmt.Println("Deallocate")
		}
	}
	return nil
}
