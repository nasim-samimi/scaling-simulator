package main

import (
	"log"
	"path/filepath"

	src "github.com/nasim-samimi/scaling-simulator/pkg/scaling"
)

func main() {
	initialise()
}

//here must initialise the orchestrator
// turn on some nodes and initialise them
// initialise domain
// initialise cloud

// initialise the task as it arrives

// initialise the node as it turns on
func initialise() {
	CloudNodes := io.loadCloudFromCSV("data/cloud.csv")
	cloud := src.NewCloud(CloudNodes)
	// read domain csv files in domain folder
	domainFilesNames, err := filepath.Glob("data/domain/*.csv")
	if err != nil {
		log.Fatal(err)
	}
	var domains []*src.Domain
	for _, fileName := range domainFilesNames {
		domainNodes := io.loadDomainFromCSV(fileName)
		domains = append(domains, src.NewDomain(domainNodes))
	}
	// read task csv files in task folder
	svcFilesNames, err := filepath.Glob("data/task/*.csv")
	if err != nil {
		log.Fatal(err)
	}
	var svcs []*src.Task
	for _, fileName := range svcFilesNames {
		svcs = append(svcs, io.loadTaskFromCSV(fileName))
	}

}
