package scaling

import (
	"fmt"
)

type NodeName string
type AllocatedTasks map[TaskID]*Task

type Node struct {
	Cores                    Cores
	ReallocHeuristic         Heuristic
	NodeName                 NodeName
	NodeAdmission            *AdmissionTest
	Location                 Location
	DomainID                 DomainID
	AllocatedTasks           AllocatedTasks
	AverageResidualBandwidth float64
	TotalResidualBandwidth   float64
}

func NewNode(cores Cores, heuristic Heuristic, nodeName NodeName) *Node {
	admissionTest := NewAdmissionTest(cores, heuristic)
	fmt.Println("Node Admission Test: ", admissionTest)
	return &Node{
		Cores:            cores,
		ReallocHeuristic: heuristic,
		NodeName:         nodeName,
		NodeAdmission:    admissionTest,
	}
}

func (n *Node) NodeAllocate(reqCpus uint64, reqBandwidth float64, task *Task) ([]CoreID, error) {
	selectedCpus, err := n.NodeAdmission.Admission(reqCpus, reqBandwidth, n.Cores)
	if err != nil {
		return selectedCpus, err
	}
	for _, coreID := range selectedCpus {
		core := n.Cores[coreID]
		core.ConsumedBandwidth += reqBandwidth
		n.Cores[coreID] = core
	}
	n.AllocatedTasks[task.taskID] = task
	return selectedCpus, nil
}

func (n *Node) IntraDomainReallocateTest(newTask *Task, oldTaskID TaskID) (bool, error) {
	NewCores := n.Cores
	oldTaskCores := n.AllocatedTasks[oldTaskID].allocatedCoresEdge
	bandwidth := n.AllocatedTasks[oldTaskID].standardMode.bandwidthEdge
	for _, coreID := range oldTaskCores {
		NewCores[coreID].ConsumedBandwidth -= bandwidth
	}

	_, err := n.NodeAdmission.Admission(newTask.standardMode.cpusEdge, newTask.standardMode.bandwidthEdge, NewCores)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (n *Node) IntraNodeReallocateTest(newTask *Task, oldTaskID TaskID) (bool, error) {
	NewCores := n.Cores
	oldTaskCores := n.AllocatedTasks[oldTaskID].allocatedCoresEdge
	oldBandwidth := n.AllocatedTasks[oldTaskID].standardMode.bandwidthEdge
	newBandwidth := newTask.standardMode.bandwidthEdge
	for _, coreID := range oldTaskCores {
		NewCores[coreID].ConsumedBandwidth -= oldBandwidth
	}

	possibleCores, err := n.NodeAdmission.Admission(newTask.standardMode.cpusEdge, newTask.standardMode.bandwidthEdge, NewCores)
	if err == nil {
		for _, coreID := range possibleCores {
			NewCores[coreID].ConsumedBandwidth += newBandwidth
		}
		_, err = n.NodeAdmission.Admission(n.AllocatedTasks[oldTaskID].standardMode.cpusEdge, oldBandwidth, NewCores)
		if err == nil {
			return true, nil
		}
	}
	return false, err
}

func (n *Node) NodeDeallocate(taskID TaskID) bool {
	cores := n.AllocatedTasks[taskID].allocatedCoresEdge
	mode := n.AllocatedTasks[taskID].allocationMode
	for _, core := range cores {
		switch mode {
		case StandardMode:
			n.Cores[core].ConsumedBandwidth -= n.AllocatedTasks[taskID].standardMode.bandwidthEdge
		case ReducedMode:
			n.Cores[core].ConsumedBandwidth -= n.AllocatedTasks[taskID].reducedMode.bandwidthEdge
		}
	}
	delete(n.AllocatedTasks, taskID)
	return true
}
