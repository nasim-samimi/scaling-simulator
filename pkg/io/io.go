package io

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
)

func loadDomainFromCSV(filename string) src.Nodes {
	// Load the CSV file
	// Create a map of nodes
	// For each row in the CSV file
	// Create a new node
	// Add the node to the map
	// Return the map

	return nil
}

func loadTaskFromCSV(filename string) src.Tasks {
	// Load the CSV file
	// Create a map of nodes
	// For each row in the CSV file
	// Create a new node
	// Add the node to the map
	// Return the map

	return nil
}

func loadCloudFromCSV(filePath string) []Node {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var nodes []Node
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity
		nodeID := record[0]
		residualBandwidth, _ := strconv.ParseFloat(record[1], 64)
		capacity, _ := strconv.Atoi(record[2])

		node := Node{
			NodeID:            nodeID,
			ResidualBandwidth: residualBandwidth,
			Capacity:          capacity,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// Load events from CSV file
func loadEventsFromCSV(filePath string) []Event {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s, %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse CSV file %s, %v", filePath, err)
	}

	var events []Event
	for _, record := range records {
		// Assuming CSV has columns: EventType, TargetNodeID, Details
		eventType := record[0]
		targetNodeID := record[1]
		details := record[2]

		event := Event{
			Type:         eventType,
			TargetNodeID: targetNodeID,
			Details:      details,
		}
		events = append(events, event)
	}
	return events
}
