package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
)

func loadNodesFromCSV(filePath string) src.Nodes {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var nodes src.Nodes
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		nodeName := record[0]
		numCores, _ := strconv.Atoi(record[2])
		cores := src.CreateNodeCores(numCores)
		partitioningHeuristic, _ := strconv.Atoi(record[3])

		newNode := src.NewNode(cores, src.Heuristic(partitioningHeuristic), src.NodeName(nodeName))
		nodes[src.NodeName(nodeName)] = newNode
		// 	nodes = append(nodes, newNode)
	}
	return nodes
}

func LoadCloudFromCSV(filePath string) src.Nodes {
	return loadNodesFromCSV(filePath)
}

func LoadDomainFromCSV(filename string) src.Nodes {
	return loadNodesFromCSV(filename)

}

func LoadSVCFromCSV(filePath string) src.Services {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var svcs src.Services
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		svcID := record[0]
		standardBandwidth, _ := strconv.Atoi(record[2])
		standardCores, _ := strconv.Atoi(record[3])
		reducedEdgeBandwidth, _ := strconv.Atoi(record[4])
		reducedEdgeCores, _ := strconv.Atoi(record[5])
		reducedCloudBandwidth, _ := strconv.Atoi(record[6])
		reducedCloudCores, _ := strconv.Atoi(record[7])
		importanceFactor, _ := strconv.Atoi(record[8])

		newSVC := src.NewService(float64(importanceFactor), src.ServiceID(svcID), float64(standardBandwidth), uint64(standardCores), float64(reducedEdgeBandwidth), uint64(reducedEdgeCores), float64(reducedCloudBandwidth), uint64(reducedCloudCores))

		svcs[src.ServiceID(svcID)] = newSVC
		// 	nodes = append(nodes, newNode)
	}
	return svcs
}

func LoadHeuristicFromCSV(filePath string) (src.Heuristic, src.Heuristic) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var nodeHeu, reallocHeu string
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		nodeHeu = record[0]
		reallocHeu = record[1]
		break
	}
	return src.Heuristic(nodeHeu), src.Heuristic(reallocHeu)
}

// must define a parameter validation function specially for heuristics

// // Load events from CSV file
func LoadEventsFromCSV(filePath string) []src.Event {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var events []src.Event
	for _, record := range records {
		// Assuming CSV has columns: EventType, TargetNodeID, Details
		eventType := record[0]
		targetDomainID := record[1]
		TargetServiceID := record[2]

		event := src.Event{
			EventType:       eventType,
			TargetDomainID:  src.DomainID(targetDomainID),
			TargetServiceID: src.ServiceID(TargetServiceID),
		}
		events = append(events, event)
	}
	return events
}