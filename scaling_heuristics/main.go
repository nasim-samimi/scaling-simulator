package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
	util "github.com/nasim-samimi/scaling-simulator/pkg/util"
)

func initialise() *src.Orchestrator {
	CloudNodes := util.LoadCloudFromCSV("../data/cloud.csv")
	cloud := src.NewCloud(CloudNodes)
	// read domain csv files in domain folder
	domainFilesNames, err := filepath.Glob("../data/domain/*.csv")
	if err != nil {
		log.Fatal(err)
	}
	domains := make(src.Domains)
	i := 1
	for _, fileName := range domainFilesNames {
		id := strconv.Itoa(i)
		i++
		domainNodes := util.LoadDomainFromCSV(fileName)
		domains[src.DomainID(id)] = src.NewDomain(domainNodes, src.DomainID(id))
	}

	svcs := util.LoadSVCFromCSV("../data/svcs.csv")
	nodeHeuristic, reallocHeuristic := util.LoadHeuristicFromCSV("../data/heuristics.csv")
	fmt.Println("Node Heuristic:", nodeHeuristic)
	fmt.Println("Realloc Heuristic:", reallocHeuristic)

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(src.NodeSelectionHeuristic(nodeHeuristic), src.ReallocationHeuristic(reallocHeuristic), cloud, domains, svcs)
	return orchestrator
}

func processEvents(orchestrator *src.Orchestrator) ([]int, error) {
	events := util.LoadEventsFromCSV("../data/events.csv")
	qosPerCost := make([]int, 0)
	test := 0
	for _, event := range events {
		fmt.Println("event:", event)
		if event.EventType == "allocate" {
			allocated, qos, cost, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID)
			fmt.Println("Allocate:", allocated, qos, cost)
			if err != nil {
				fmt.Println(err)
			}
			if allocated {
				qosPerCost = append(qosPerCost, int(qos)/int(cost))
			}
			test++
			// if test == 2 {
			// 	break
			// }
		}
		if event.EventType == "deallocate" {
			fmt.Println("Deallocate")
		}
	}
	fmt.Println("QoS per Cost: ", qosPerCost)
	return qosPerCost, nil
}

func main() {
	orchestrator := initialise()
	processEvents(orchestrator)
}
