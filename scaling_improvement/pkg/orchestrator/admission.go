package orchestrator

import (
	"fmt"
	"sort"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

type AdmissionTest struct {
	Cores     Cores
	Heuristic cnfg.Heuristic
}

func NewAdmissionTest(core Cores, heuristic cnfg.Heuristic) *AdmissionTest {
	return &AdmissionTest{
		Cores:     core,
		Heuristic: heuristic,
	}
}

func (at *AdmissionTest) QuickFilter(reqCpus uint64, reqBandwidth float64, cores Cores) (bool, error) {
	i := 0
	fmt.Println("average added bandwidth for quick filter:", reqBandwidth*float64(reqCpus)/float64(len(cores)))
	fmt.Println("bandwidth and cpus requested for quick filter:", reqBandwidth, reqCpus)
	for _, cpuinfo := range cores {
		if cpuinfo.ConsumedBandwidth+reqBandwidth <= 100.0 {
			i++
		}
		if reqCpus == uint64(i) {
			return true, nil
		}
	}
	return false, fmt.Errorf("not enough cpus to allocate")
}

func (at *AdmissionTest) Admission(reqCpus uint64, reqBandwidth float64, cores Cores, cpuThreshold float64) ([]CoreID, error) {
	// const cpuThreshold = 100
	type scoredCpu struct {
		cpu   CoreID
		score float64
	}
	var scoredCpus []scoredCpu
	for _, cpuinfo := range cores {
		score := cpuThreshold - cpuinfo.ConsumedBandwidth - reqBandwidth
		if score > 0 {
			scoredCpus = append(scoredCpus, scoredCpu{
				cpu:   cpuinfo.ID,
				score: score,
			})
		}
	}

	if uint64(len(scoredCpus)) < reqCpus {
		return nil, fmt.Errorf("not enough cpus to allocate")
	}
	switch at.Heuristic {
	case "worstFit":
		sort.SliceStable(scoredCpus, func(i, j int) bool {
			if scoredCpus[i].score > scoredCpus[j].score {
				return true
			}
			return false
		})
	case "bestFit":
		sort.SliceStable(scoredCpus, func(i, j int) bool {
			if scoredCpus[i].score < scoredCpus[j].score {
				return true
			}
			return false
		})
	default:
		sort.SliceStable(scoredCpus, func(i, j int) bool {
			if scoredCpus[i].score > scoredCpus[j].score {
				return true
			}
			return false
		})
	}

	var fittingCpus []CoreID
	for i := uint64(0); i < reqCpus; i++ {
		fittingCpus = append(fittingCpus, scoredCpus[i].cpu)
	}

	return fittingCpus, nil
}
