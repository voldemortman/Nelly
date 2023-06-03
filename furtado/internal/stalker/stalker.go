package stalker

type IStalker interface {
	StartStalking() error
	StopStalking() error
}
