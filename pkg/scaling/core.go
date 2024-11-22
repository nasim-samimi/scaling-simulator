package scaling

type CoreID string
type Core struct {
	ID                CoreID
	ConsumedBandwidth float64
}

type Cores map[CoreID]*Core
