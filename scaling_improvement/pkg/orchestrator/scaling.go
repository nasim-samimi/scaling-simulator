package orchestrator

import cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"

func (o *Orchestrator) edgePowerOffNode(domainID DomainID, nodeName NodeName) bool {
	for _, noden := range o.Domains[domainID].AlwaysActiveNodes {
		if noden == nodeName {
			return false
		}
	}
	cores := CreateNodeCores(len(o.Domains[domainID].ActiveNodes[nodeName].Cores))
	o.Domains[domainID].InactiveNodes[nodeName] = NewNode(cores, o.Domains[domainID].ActiveNodes[nodeName].ReallocHeuristic, nodeName, domainID)
	o.Cost = o.Cost - o.Config.EdgeNodeCost*cnfg.Cost(len(cores))
	delete(o.Domains[domainID].ActiveNodes, nodeName)

	return true
}

func (o *Orchestrator) cloudPowerOffNode(nodeName NodeName) bool {
	o.Cloud.InactiveNodes[nodeName] = NewNode(o.Cloud.ActiveNodes[nodeName].Cores, o.Cloud.ActiveNodes[nodeName].ReallocHeuristic, nodeName, "")
	o.Cost = o.Cost - o.Config.CloudNodeCost*cnfg.Cost(len(o.Cloud.ActiveNodes[nodeName].Cores))
	delete(o.Cloud.ActiveNodes, nodeName)

	return true
}

func (o *Orchestrator) edgePowerOnNode(domainID DomainID) (bool, NodeName) {
	// fmt.Println("active nodes in domain:", o.Domains[domainID].ActiveNodes)
	// var nodeName NodeName
	for nodeName, node := range o.Domains[domainID].InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Domains[domainID].ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, domainID)
		o.Cost = o.Cost + o.Config.EdgeNodeCost*cnfg.Cost(len(cores))
		delete(o.Domains[domainID].InactiveNodes, nodeName)
		return true, nodeName
	}
	// fmt.Println("active nodes in domain after powering on:", o.Domains[domainID].ActiveNodes)
	return false, ""
}
func (o *Orchestrator) cloudPowerOnNode() bool {
	// fmt.Println("active nodes in cloud:", o.Cloud.ActiveNodes)
	for nodeName, node := range o.Cloud.InactiveNodes {
		node.Status = Active
		cores := CreateNodeCores(len(node.Cores))
		o.Cloud.ActiveNodes[nodeName] = NewNode(cores, node.ReallocHeuristic, nodeName, "")
		o.Cost = o.Cost + o.Config.CloudNodeCost*cnfg.Cost(len(cores))
		delete(o.Cloud.InactiveNodes, nodeName)
		break
	}
	// fmt.Println("active nodes in cloud after powering on:", o.Cloud.ActiveNodes)

	return true
}
