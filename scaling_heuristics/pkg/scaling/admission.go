package scaling

import (
	"fmt"
	"sort"
)

type AdmissionTest struct {
	Cores     Cores
	Heuristic Heuristic
}

func NewAdmissionTest(core Cores, heuristic Heuristic) *AdmissionTest {
	return &AdmissionTest{
		Cores:     core,
		Heuristic: heuristic,
	}
}

func (at *AdmissionTest) Admission(reqCpus uint64, reqBandwidth float64, cores Cores) ([]CoreID, error) {
	const cpuThreshold = 0.95
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