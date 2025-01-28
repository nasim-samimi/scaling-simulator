package util

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

func Initialise(config *cnfg.Config) *src.Orchestrator {

	CloudNodes, reservedCloudNodes := LoadCloudFromCSV("../data/cloud.csv", "cloud")
	cloud := src.NewCloud(CloudNodes, reservedCloudNodes)
	// read domain csv files in domain folder
	svcs := LoadSVCFromCSV("../data/services/services0.csv")
	// nodeHeuristic, reallocHeuristic, partitionHeuristic := util.LoadHeuristicFromCSV("../data/heuristics.csv")
	nodeHeuristic := config.NodeHeuristic
	partitionHeuristic := config.PartitionHeuristic

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
		domainNodes, reservedNodes := LoadDomainFromCSV(fileName, "domain", src.DomainID(id))
		// fmt.Println("domain ID:", id)
		// fmt.Println("name of the file:", fileName)
		domains[src.DomainID(id)] = src.NewDomain(domainNodes, reservedNodes, src.DomainID(id))
	}

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(config, cloud, domains, svcs)
	fmt.Println("the addition from config file:", config.Addition)
	return orchestrator
}
