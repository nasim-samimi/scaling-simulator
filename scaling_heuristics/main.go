package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
	util "github.com/nasim-samimi/scaling-simulator/pkg/util"
)

func addition() string {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided to show randomness!")
		fmt.Println("Usage: go run main.go <addition>")
		return "0"
	}
	fmt.Println("Received arguments:")
	// for i, arg := range os.Args[1:] {
	// 	fmt.Printf("Arg %d: %s\n", i+1, arg)
	// }
	addition := os.Args[1]
	return addition
}
func initialise() *src.Orchestrator {
	CloudNodes, reservedCloudNodes := util.LoadCloudFromCSV("../data/cloud.csv")
	cloud := src.NewCloud(CloudNodes, reservedCloudNodes)
	// read domain csv files in domain folder
	svcs := util.LoadSVCFromCSV("../data/services/services0.csv")
	nodeHeuristic, reallocHeuristic, partitionHeuristic := util.LoadHeuristicFromCSV("../data/heuristics.csv")
	fmt.Println("Node Heuristic:", nodeHeuristic)
	fmt.Println("Realloc Heuristic:", reallocHeuristic)

	domainFilesNames, err := filepath.Glob("../data/domainNodes/" + string(partitionHeuristic) + "/" + string(nodeHeuristic) + "/*.csv")
	fmt.Println("../data/domainNodes/" + string(nodeHeuristic) + "/" + string(partitionHeuristic))
	fmt.Println(domainFilesNames)
	if err != nil {
		log.Fatal(err)
	}
	domains := make(src.Domains)
	i := 0
	for _, fileName := range domainFilesNames {
		id := strconv.Itoa(i)
		i++
		domainNodes, reservedNodes := util.LoadDomainFromCSV(fileName)
		// fmt.Println("domain ID:", id)
		// fmt.Println("name of the file:", fileName)
		domains[src.DomainID(id)] = src.NewDomain(domainNodes, reservedNodes, src.DomainID(id))
	}

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(src.NodeSelectionHeuristic(nodeHeuristic), src.ReallocationHeuristic(reallocHeuristic), src.Heuristic(partitionHeuristic), cloud, domains, svcs)
	return orchestrator
}

func processEvents(orchestrator *src.Orchestrator) error {
	addition := addition()
	events := util.LoadEventsFromCSV("../data/events_" + addition + ".csv")
	qosPerCost := make([]float64, 0)
	durations := make([]float64, 0)
	test := 0
	for _, event := range events {
		eventID := event.EventID
		fmt.Println("event:", event)
		fmt.Println("service:", orchestrator.AllServices[event.TargetServiceID])
		if event.EventType == "allocate" {
			startTime := time.Now()
			allocated, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID, eventID)
			duration := time.Since(startTime)
			fmt.Println("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
			fmt.Println("Time to allocate:", duration)
			if err != nil {
				fmt.Println(err)
			}
			if allocated {
				qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
				durations = append(durations, float64(duration.Microseconds())/1000)
				fmt.Println("/////////////////////")
				fmt.Println("service is allocated ")
				fmt.Println("/////////////////////")
			} else {
				fmt.Println("/////////////////////")
				fmt.Println("service is rejected ")
				fmt.Println("/////////////////////")
			}
			test++
			// if test == 50 {
			// 	break
			// }
		}
		if event.EventType == "deallocate" {
			fmt.Println("/////////////////////")
			fmt.Println("Deallocate")
			orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)
			fmt.Println("/////////////////////")
		}
	}
	fmt.Println("QoS per Cost: ", qosPerCost)
	// save to csv file
	name := string(orchestrator.NodeSelectionHeuristic) + "_" + string(orchestrator.PartitionHeuristic) + "_" + string(orchestrator.ReallocationHeuristic) + "_" + addition
	util.WriteToCsv("../experiments/results/heuristics/qosPerCost_"+name+".csv", qosPerCost)
	util.WriteToCsv("../experiments/results/heuristics/runtimes_"+name+".csv", durations)
	fmt.Println("Durations: ", durations)
	return nil
}

func main() {
	orchestrator := initialise()
	processEvents(orchestrator)
}
