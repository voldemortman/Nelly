package furtado

type IStalker interface {
	StartStalking() error
	StopStalking() error
}
