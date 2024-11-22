package scaling

type reducedSched struct {
	bandwidthEdge  float64
	cpusEdge       uint64
	bandwidthCloud float64
	cpusCloud      uint64
	task           *Task
}

type standardSched struct {
	bandwidthEdge float64
	cpusEdge      uint64
	task          *Task
}

func (s *standardSched) TaskAllocate(node *Node) (bool, error) {
	allocated := false

	allocatedCores, err := node.NodeAllocate(s.cpusEdge, s.bandwidthEdge, s.task)
	if err != nil {
		return allocated, err
	}
	s.task.allocatedNodeEdge = node.NodeName
	allocated = true
	s.task.allocatedCoresEdge = allocatedCores
	s.task.allocationMode = StandardMode

	s.task.AverageResidualBandwidth = s.bandwidthEdge
	s.task.TotalResidualBandwidth = s.bandwidthEdge * float64(s.cpusEdge)
	return allocated, err
}

func (s *standardSched) TaskDeallocate(task *Task, node *Node) (*Task, error) {
	node.NodeDeallocate(task.taskID)
	return task, nil
}

func (r *reducedSched) TaskAllocate(node *Node, loc Location) (bool, error) {
	allocated := false

	switch loc {
	case edgeLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusEdge, r.bandwidthEdge, r.task)

		if err != nil {
			return allocated, err
		}
		r.task.allocatedNodeEdge = node.NodeName
		allocated = true
		r.task.allocatedCoresEdge = allocatedCores
		r.task.allocationMode = ReducedMode

	case cloudLoc:
		allocatedCores, err := node.NodeAllocate(r.cpusCloud, r.bandwidthCloud, r.task)
		if err != nil {
			return false, err
		}
		r.task.allocatedNodeCloud = node.NodeName
		r.task.allocatedCoresCloud = allocatedCores
		allocated = true
	}
	r.task.AverageResidualBandwidth = (r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)) / (float64(r.cpusEdge) + float64(r.cpusCloud))
	r.task.TotalResidualBandwidth = r.bandwidthEdge*float64(r.cpusEdge) + r.bandwidthCloud*float64(r.cpusCloud)

	return true, nil

}

func (r *reducedSched) TaskDeallocate(task *Task, node *Node, location Location) (*Task, error) {
	switch location {
	case edgeLoc:
		node.NodeDeallocate(task.taskID)
	case cloudLoc:
		node.NodeDeallocate(task.taskID)
	}
	return task, nil
}

type TaskID string
type TaskMode string

const (
	ReducedMode  TaskMode = "Reduced"
	StandardMode TaskMode = "Standard"
)

type Task struct {
	importanceFactor         uint64
	taskID                   TaskID
	reducedMode              reducedSched
	standardMode             standardSched
	allocatedCoresEdge       []CoreID
	allocatedCoresCloud      []CoreID
	allocatedNodeEdge        NodeName
	allocatedNodeCloud       NodeName
	allocationMode           TaskMode
	AverageResidualBandwidth float64
	TotalResidualBandwidth   float64
	// taskModel           TaskModel
}

func NewTask(importanceFactor uint64, taskID TaskID, reduced reducedSched, standard standardSched) *Task {
	task := &Task{
		importanceFactor: importanceFactor,
		taskID:           taskID,
		reducedMode:      reduced,
		standardMode:     standard,
		// taskModel:        taskModel,
	}
	task.reducedMode.task = task
	task.standardMode.task = task
	return task
}

func (t *Task) TaskReallocate(task *Task, node *Node, domain *Domain, location Location) bool {
	return true
}
