package orchestrator

import cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"

func (o *Orchestrator) edgePowerOffNode(domainID DomainID, nodeName NodeName) bool {
	// for _, noden := range o.Domains[domainID].AlwaysActiveNodes {
	// 	if noden == nodeName {
	// 		return false
	// 	}
	// }
	if len(o.Domains[domainID].ActiveNodes) == o.Domains[domainID].MinActiveNodes {
		log.Info("cannot power off from the always active nodes of domain:", o.Domains[domainID].MinActiveNodes)
		log.Info("active nodes in domain:", o.Domains[domainID].ActiveNodes)
		return false
	}
	cores := CreateNodeCores(len(o.Domains[domainID].ActiveNodes[nodeName].Cores))
	o.Domains[domainID].InactiveNodes[nodeName] = NewNode(cores, o.Domains[domainID].ActiveNodes[nodeName].ReallocHeuristic, nodeName, domainID)
	o.Cost = o.Cost - o.Config.EdgeNodeCost*cnfg.Cost(len(cores))
	delete(o.Domains[domainID].ActiveNodes, nodeName)
	log.Info("cost is update to ", o.Cost, "with reducing edge cost of", o.Config.EdgeNodeCost*cnfg.Cost(len(cores)))
	return true
}

func (o *Orchestrator) cloudPowerOffNode(nodeName NodeName) bool {
	cores := CreateNodeCores(len(o.Cloud.ActiveNodes[nodeName].Cores))
	o.Cloud.InactiveNodes[nodeName] = NewNode(cores, o.Cloud.ActiveNodes[nodeName].ReallocHeuristic, nodeName, "")
	o.Cost = o.Cost - o.Config.CloudNodeCost*cnfg.Cost(len(cores))
	delete(o.Cloud.ActiveNodes, nodeName)
	log.Info("cost is update to ", o.Cost, "with reducing cloud cost of", o.Config.CloudNodeCost*cnfg.Cost(len(cores)))
	return true
}

func (o *Orchestrator) edgePowerOnNode(domainID DomainID) (bool, NodeName) {
	// log.Info("active nodes in domain:", o.Domains[domainID].ActiveNodes)
	// var nodeName NodeName
	for nodeName, node := range o.Domains[domainID].InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Domains[domainID].ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, domainID)
		o.Cost = o.Cost + o.Config.EdgeNodeCost*cnfg.Cost(len(cores))
		delete(o.Domains[domainID].InactiveNodes, nodeName)
		log.Info("cost is update to ", o.Cost, "with adding edge cost of", o.Config.EdgeNodeCost*cnfg.Cost(len(cores)))
		return true, nodeName
	}
	// log.Info("active nodes in domain after powering on:", o.Domains[domainID].ActiveNodes)
	return false, ""
}
func (o *Orchestrator) cloudPowerOnNode() bool {
	// log.Info("active nodes in cloud:", o.Cloud.ActiveNodes)
	for nodeName, node := range o.Cloud.InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Cloud.ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, "")
		o.Cost = o.Cost + o.Config.CloudNodeCost*cnfg.Cost(len(cores))
		delete(o.Cloud.InactiveNodes, nodeName)
		log.Info("cost is update to ", o.Cost, "with adding cloud cost of", o.Config.CloudNodeCost*cnfg.Cost(len(cores)))
		break
	}
	// log.Info("active nodes in cloud after powering on:", o.Cloud.ActiveNodes)

	return true
}
