package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

func loadNodesFromCSV(filePath string, loc string, domainID src.DomainID) (src.Nodes, src.Nodes) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read and ignore the first row (headers)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Unable to read header row from CSV file %s, %v", filePath, err)
	}
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	nodes := make(src.Nodes)
	reservedNodes := make(src.Nodes)
	// l := int(math.Ceil(float64(len(records)/10))) + 1
	l := 1
	fmt.Println("l:", l)
	i := 0
	// lr := len(records) - l
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		if i < l {
			nodeName := record[0]
			numCores, _ := strconv.Atoi(record[1])
			cores := src.CreateNodeCores(numCores)
			partitioningHeuristic := record[2]

			newNode := src.NewNode(cores, src.Heuristic(partitioningHeuristic), src.NodeName(nodeName), domainID)
			nodes[src.NodeName(nodeName)] = newNode

		} else {
			nodeName := record[0]
			numCores, _ := strconv.Atoi(record[1])
			cores := src.CreateNodeCores(numCores)
			partitioningHeuristic := record[2]

			newNode := src.NewNode(cores, src.Heuristic(partitioningHeuristic), src.NodeName(nodeName), domainID)
			reservedNodes[src.NodeName(nodeName)] = newNode
		}
		i++
		// 	nodes = append(nodes, newNode)
	}
	return nodes, reservedNodes
}

func LoadCloudFromCSV(filePath string, loc string) (src.Nodes, src.Nodes) {
	return loadNodesFromCSV(filePath, loc, "")
}

func LoadDomainFromCSV(filename string, loc string, domainID src.DomainID) (src.Nodes, src.Nodes) {
	return loadNodesFromCSV(filename, loc, domainID)

}

func LoadSVCFromCSV(filePath string) src.Services {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read and ignore the first row (headers)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Unable to read header row from CSV file %s, %v", filePath, err)
	}
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	svcs := make(src.Services)
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		svcID := record[0]
		importanceFactor, _ := strconv.Atoi(record[1])
		standardBandwidth, _ := strconv.Atoi(record[2])
		standardCores, _ := strconv.Atoi(record[3])
		reducedEdgeBandwidth, _ := strconv.ParseFloat(record[4], 64)
		reducedEdgeCores, _ := strconv.ParseFloat(record[5], 64)
		reducedCloudBandwidth, _ := strconv.ParseFloat(record[6], 64)
		reducedCloudCores, _ := strconv.ParseFloat(record[7], 64)
		// print reduced parameters:
		fmt.Println(" reduced parameters:", reducedEdgeBandwidth, reducedEdgeCores, reducedCloudBandwidth, reducedCloudCores)

		newSVC := src.NewService(float64(importanceFactor), src.ServiceID(svcID), float64(standardBandwidth), uint64(standardCores), float64(reducedEdgeBandwidth), uint64(reducedEdgeCores), float64(reducedCloudBandwidth), uint64(reducedCloudCores))

		svcs[src.ServiceID(svcID)] = newSVC
		// 	nodes = append(nodes, newNode)
	}
	return svcs
}

func LoadHeuristicFromCSV(filePath string) (src.Heuristic, src.Heuristic, src.Heuristic) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read and ignore the first row (headers)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Unable to read header row from CSV file %s, %v", filePath, err)
	}
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var nodeHeu, reallocHeu, partitHeu string
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		nodeHeu = record[1]
		reallocHeu = record[0]
		partitHeu = record[2]
		break
	}
	return src.Heuristic(nodeHeu), src.Heuristic(reallocHeu), src.Heuristic(partitHeu)
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
	// Read and ignore the first row (headers)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Unable to read header row from CSV file %s, %v", filePath, err)
	}
	records, err := reader.ReadAll()
	// fmt.Println(records)
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var events []src.Event
	i := 0
	for _, record := range records {
		// Assuming CSV has columns: EventType, TargetNodeID, Details
		eventTime, _ := strconv.Atoi(record[0])
		eventType := record[1]
		TargetServiceID := record[2]
		targetDomainID := record[3]
		eventID := record[4]

		event := src.Event{
			EventType:       eventType,
			TargetDomainID:  src.DomainID(targetDomainID),
			TargetServiceID: src.ServiceID(TargetServiceID),
			EventTime:       eventTime,
			EventID:         src.ServiceID(eventID),
		}
		events = append(events, event)
		i++
	}
	return events
}
