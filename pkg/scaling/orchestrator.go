package scaling

// Node selection function
// largest residual bandwidth
// don't undestand
type Heuristic string
type NodeSelectionHeuristic Heuristic
type ReallocationHeuristic Heuristic
type Location string

const (
	cloudLoc Location = "cloud"
	edgeLoc  Location = "edge"
	bothLoc  Location = "both"
)

const (
	HBI   ReallocationHeuristic = "HBI"
	HCI   ReallocationHeuristic = "HCI"
	HBCI  ReallocationHeuristic = "HBCI"
	HBIcC ReallocationHeuristic = "HBIcC"
)

const (
	MinMin NodeSelectionHeuristic = "MinMin"
	MaxMax NodeSelectionHeuristic = "MaxMax"
)

type Orchestrator struct {
	NodeSelectionHeuristic NodeSelectionHeuristic
	ReallocationHeuristic  ReallocationHeuristic
	domains                Domains
	cloud                  *Cloud
	AllServices            []*Task
	RunningServices        []*Task // change name of task to service
}

func NewOrchestrator(nodeSelectionHeuristic NodeSelectionHeuristic, reallocationHeuristic ReallocationHeuristic) *Orchestrator {
	return &Orchestrator{
		NodeSelectionHeuristic: nodeSelectionHeuristic,
		ReallocationHeuristic:  reallocationHeuristic,
	}
}

func (o *Orchestrator) SelectNode(task *Task) *Node {
	var node *Node

	return node
}

func (o *Orchestrator) allocateEdge(task *Task, node *Node, domain *Domain) (bool, error) {
	allocated, err := task.standardMode.TaskAllocate(node)
	return allocated, err
}

func (o *Orchestrator) getReallocatedTask(node *Node, t *Task) (TaskID, error) {
	var selectedTaskID TaskID
	switch o.ReallocationHeuristic {
	case HBI:
		hbi := 0.0
		for _, task := range node.AllocatedTasks {
			if task.allocationMode == StandardMode {
				bi := float64(task.importanceFactor) * task.standardMode.bandwidthEdge
				if bi > hbi {
					selectedTaskID = task.taskID
					hbi = bi
				}
			}
		}
	case HCI:
		hci := uint64(0)

		for _, task := range node.AllocatedTasks {
			if task.allocationMode == StandardMode {
				ci := task.importanceFactor * task.standardMode.cpusEdge
				if ci > hci {
					selectedTaskID = task.taskID
					hci = ci
				}
			}
		}
	case HBCI:
		hbci := 0.0
		for _, task := range node.AllocatedTasks {
			if task.allocationMode == StandardMode {
				bci := float64(task.importanceFactor) * task.standardMode.bandwidthEdge * float64(task.standardMode.cpusEdge)
				if bci > hbci {
					selectedTaskID = task.taskID
					hbci = bci
				}
			}
		}
	case HBIcC:
		hbic := 0.0
		for _, task := range node.AllocatedTasks {
			if task.allocationMode == StandardMode {
				if task.standardMode.cpusEdge > t.standardMode.cpusEdge {
					bic := float64(task.importanceFactor) * task.standardMode.bandwidthEdge * float64(task.standardMode.cpusEdge)
					if bic > hbic {
						selectedTaskID = task.taskID
						hbic = bic
					}
				}
			}
		}

	default: // HBI
		hbi := 0.0
		for _, task := range node.AllocatedTasks {
			if task.allocationMode == StandardMode {
				bi := float64(task.importanceFactor) * task.standardMode.bandwidthEdge
				if bi > hbi {
					selectedTaskID = task.taskID
					hbi = bi
				}
			}
		}
	}
	return selectedTaskID, nil
}

func (o *Orchestrator) sortNodes(nodes Nodes) ([]NodeName, error) {
	// sort nodes according to the heuristic
	switch o.NodeSelectionHeuristic {
	case MinMin:
		for _, node := range nodes {

			for _, cores := range node.Cores {
				// sort cores
			}
		}
	case MaxMax:
		for _, node := range nodes {
			for _, cores := range node.Cores {
				// sort cores
			}
		}
	}
	return nil, nil
}

func (o *Orchestrator) intraNodeRealloc(task *Task, node *Node) (bool, error) {
	reallocatedTask := &Task{}
	reallocatedTaskID, err := o.getReallocatedTask(node, task)
	if err != nil {
		return false, err
	}

	reallocation, err := node.IntraNodeReallocateTest(task, reallocatedTaskID)
	if err != nil {
		return false, err
	}
	if reallocation {
		return true, nil
	}

	reallocatedTask, err = task.standardMode.TaskDeallocate(node.AllocatedTasks[reallocatedTaskID], node)
	if err != nil {
		return false, err
	}

	task.standardMode.TaskAllocate(node)
	reallocatedTask.standardMode.TaskAllocate(node)

	return true, nil
}

func (o *Orchestrator) intraDomainRealloc(task *Task, node *Node, domain *Domain) (bool, error) {
	searchingNodes := domain.OnNodes
	delete(searchingNodes, node.NodeName)
	sortedNodes, _ := o.sortNodes(searchingNodes)
	otherTaskID, _ := o.getReallocatedTask(node, task)
	reallocated := false

	otherTask := node.AllocatedTasks[otherTaskID]
	reallocation, _ := node.IntraDomainReallocateTest(task, otherTaskID)
	if reallocation {
		for _, nodeName := range sortedNodes {
			otherNode := domain.OnNodes[nodeName]

			allocatedCore, _ := otherNode.NodeAdmission.Admission(task.standardMode.cpusEdge, task.standardMode.bandwidthEdge, otherNode.Cores)
			if allocatedCore != nil {

				otherTask.standardMode.TaskDeallocate(task, node)
				task.standardMode.TaskAllocate(node)
				otherTask.standardMode.TaskAllocate(otherNode)
				reallocated = true
			}
			break
		}

	}
	return reallocated, nil
}

func (o *Orchestrator) SplitSched(task *Task, domain *Domain) (bool, bool, error) {
	// edge-cloud split (has qos degradation) -- there is no cloud only apparently
	sortedNodes, _ := o.sortNodes(domain.OnNodes)
	edgeAllocated := false
	cloudAllocated := false
	for _, edgeNodeName := range sortedNodes {
		node := domain.OnNodes[edgeNodeName]
		edgeAllocated, _ = task.reducedMode.TaskAllocate(node, edgeLoc)
		if edgeAllocated {
			break
		}
	}

	sortedNodes, _ = o.sortNodes(o.cloud.OnNodes)
	for _, cloudNodeName := range sortedNodes {
		node := o.cloud.OnNodes[cloudNodeName]
		cloudAllocated, _ = task.reducedMode.TaskAllocate(node, cloudLoc)
		if cloudAllocated {
			break
		}
	}
	return edgeAllocated, cloudAllocated, nil
}

func (o *Orchestrator) PowerOffNode() bool {
	return true
}

func (o *Orchestrator) edgePowerOnNode() bool {
	// add a node id and initialise it
	// NewNode()
	return true
}
func (o *Orchestrator) cloudPowerOnNode() bool {
	// add a node id and initialise it
	// NewNode()
	return true
}

func (o *Orchestrator) Allocate(domain *Domain, task *Task) bool {
	allocated := false
	// sort nodes once here
	sortedNodes, _ := o.sortNodes(domain.OnNodes)
	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ = o.allocateEdge(task, node, domain)
		if allocated {
			return allocated
		}
	}

	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ := o.intraNodeRealloc(task, node)
		if allocated {
			return allocated
		}
	}
	for _, nodeName := range sortedNodes {
		node := domain.OnNodes[nodeName]
		allocated, _ := o.intraDomainRealloc(task, node, domain)
		if allocated {
			return allocated
		}
	}
	edgeAllocated, cloudAllocated, _ := o.SplitSched(task, domain)
	if edgeAllocated && cloudAllocated {
		allocated = true
	} else {
		if !edgeAllocated {
			o.edgePowerOnNode()
		}
		if !cloudAllocated {
			o.cloudPowerOnNode()
		}
	}

	return allocated
}

func (o *Orchestrator) Deallocate() bool {
	return true
}
