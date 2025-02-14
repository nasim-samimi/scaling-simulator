package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	eve "github.com/nasim-samimi/scaling-simulator/pkg/events"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

func loadNodesFromCSV(filePath string, domainID src.DomainID) src.Nodes {
	fmt.Println("loading nodes from csv file: ", filePath)
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

	i := 0
	// lr := len(records) - l
	for _, record := range records {
		// Assuming CSV has columns: NodeID, ResidualBandwidth, Capacity

		nodeName := record[0]
		numCores, _ := strconv.Atoi(record[1])
		cores := src.CreateNodeCores(numCores)
		partitioningHeuristic := record[2]

		newNode := src.NewNode(cores, cnfg.Heuristic(partitioningHeuristic), src.NodeName(nodeName), domainID)
		nodes[src.NodeName(nodeName)] = newNode

		i++
		// 	nodes = append(nodes, newNode)
	}
	return nodes
}

func LoadCloudFromCSV(filePath string) (src.Nodes, src.Nodes) {
	return nil, loadNodesFromCSV(filePath, "")
}

func LoadDomainFromCSV(filename string, domainID src.DomainID) (src.Nodes, src.Nodes) {
	return loadNodesFromCSV(filepath.Join(filename, "Active"), domainID), loadNodesFromCSV(filepath.Join(filename, "Reserved"), domainID)

}

func LoadDomains(baseFolder string) src.Domains {
	fmt.Println("loading domains from csv file: ", baseFolder)
	activeFolder := filepath.Join(baseFolder, "Active")
	reservedFolder := filepath.Join(baseFolder, "Reserved")

	// Get all active node files
	activeFiles, err := filepath.Glob(filepath.Join(activeFolder, "*.csv"))
	if err != nil {
		log.Fatal("Error reading active files:", err)
	}

	// Create Domains
	domains := make(src.Domains)
	i := 0

	for _, activeFile := range activeFiles {
		// Get the file name (without path)
		fileName := filepath.Base(activeFile)

		// Construct the corresponding reserved file path
		reservedFile := filepath.Join(reservedFolder, fileName)

		// Check if the reserved file exists
		if _, err := filepath.Glob(reservedFile); err != nil {
			fmt.Printf("Warning: No reserved file found for %s. Skipping...\n", fileName)
			continue
		}
		id := strconv.Itoa(i)
		i++
		activeNodes := loadNodesFromCSV(activeFile, src.DomainID(id))
		reservedNodes := loadNodesFromCSV(reservedFile, src.DomainID(id))
		domains[src.DomainID(id)] = src.NewDomain(activeNodes, reservedNodes, src.DomainID(id))
	}
	return domains
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

func LoadHeuristicFromCSV(filePath string) (cnfg.Heuristic, cnfg.Heuristic, cnfg.Heuristic) {
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
	return cnfg.Heuristic(nodeHeu), cnfg.Heuristic(reallocHeu), cnfg.Heuristic(partitHeu)
}

// must define a parameter validation function specially for heuristics

// // Load events from CSV file
func LoadEventsFromCSV(filePath string) []eve.Event {
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

	var events []eve.Event
	i := 0
	for _, record := range records {
		// Assuming CSV has columns: EventType, TargetNodeID, Details
		// eventTime, _ := strconv.Atoi(record[0])
		eventTime, _ := strconv.ParseFloat(record[0], 64)
		eventType := record[1]
		TargetServiceID := record[2]
		targetDomainID := record[3]
		eventID := record[4]
		totalUtil, _ := strconv.Atoi(record[5])

		event := eve.Event{
			EventType:       eventType,
			TargetDomainID:  src.DomainID(targetDomainID),
			TargetServiceID: src.ServiceID(TargetServiceID),
			EventTime:       eventTime,
			EventID:         src.ServiceID(eventID),
			TotalUtil:       totalUtil,
		}
		events = append(events, event)
		i++
	}
	return events
}
