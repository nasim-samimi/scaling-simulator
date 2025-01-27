package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	trimmedAddition := strings.TrimSpace(addition)
	return trimmedAddition
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
		domains[src.DomainID(id)] = src.NewDomain(domainNodes, reservedNodes, src.DomainID(id))
	}

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(src.NodeSelectionHeuristic(nodeHeuristic), src.ReallocationHeuristic(reallocHeuristic), src.Heuristic(partitionHeuristic), cloud, domains, svcs)
	return orchestrator
}

func processEvents(orchestrator *src.Orchestrator) error {
	addition := addition()
	events := util.LoadEventsFromCSV("../data/events/events_" + addition + ".csv")
	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
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
				qos = append(qos, float64(orchestrator.QoS))
				cost = append(cost, float64(orchestrator.Cost))
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
	name := "addition=" + addition + "/" + string(orchestrator.NodeSelectionHeuristic) + "/" + string(orchestrator.PartitionHeuristic) + "/"
	name2 := "baseline"
	//first check if directory exists
	if _, err := os.Stat("../experiments/results/baseline/runtimes/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/baseline/runtimes/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/baseline/qosPerCost/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/baseline/qosPerCost/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/baseline/qos/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/baseline/qos/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/baseline/cost/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/baseline/cost/"+name, os.ModePerm)
	}
	util.WriteToCsv("../experiments/results/baseline/qosPerCost/"+name+name2+".csv", qosPerCost)
	util.WriteToCsv("../experiments/results/baseline/runtimes/"+name+name2+".csv", durations)
	util.WriteToCsv("../experiments/results/baseline/qos/"+name+name2+".csv", qos)
	util.WriteToCsv("../experiments/results/baseline/cost/"+name+name2+".csv", cost)
	fmt.Println("Durations: ", durations)
	return nil
}

func main() {
	orchestrator := initialise()
	processEvents(orchestrator)
}
