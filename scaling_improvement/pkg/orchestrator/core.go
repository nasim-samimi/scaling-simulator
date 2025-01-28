package orchestrator

import "fmt"

type CoreID string
type Core struct {
	ID                CoreID
	ConsumedBandwidth float64
}

type Cores map[CoreID]*Core

func NewCore(id CoreID) *Core {
	return &Core{
		ID:                id,
		ConsumedBandwidth: 0,
	}
}

func CreateNodeCores(numCores int) Cores {
	cores := make(Cores)
	for i := 0; i < numCores; i++ {
		coreID := CoreID(fmt.Sprintf("core-%d", i))
		cores[coreID] = NewCore(coreID)
	}
	return cores
}
