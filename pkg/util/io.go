package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
)

func loadCloudFromCSV(filePath string) []*src.Node {
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

	var nodes []*src.Node
	// for _, record := range records {
	// 	// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
	// 	nodeID := record[0]
	// 	cores, _ := strconv.Atoi(record[2])
	// 	heuristic,_ := strconv.Atoi(record[3])

	// 	newNode := src.NewNode(src.Cores{1: src.Core{ConsumedBandwidth: 0}}, heuristic, src.NodeName(nodeID))
	// 	nodes = append(nodes, newNode)
	// }
	return nodes
}

func loadDomainFromCSV(filename string) src.Nodes {
	// Load the CSV file
	// Create a map of nodes
	// For each row in the CSV file
	// Create a new node
	// Add the node to the map
	// Return the map

	return nil
}

// Load events from CSV file
func loadEventsFromCSV(filePath string) []src.Event {
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
	// for _, record := range records {
	// 	// Assuming CSV has columns: EventType, TargetNodeID, Details
	// 	eventType := record[0]
	// 	targetNodeID := record[1]
	// 	details := record[2]

	// event := src.Event{
	// 	Type:         eventType,
	// 	TargetNodeID: targetNodeID,
	// 	Details:      details,
	// }
	// 	events = append(events, event)
	// }
	return events
}
