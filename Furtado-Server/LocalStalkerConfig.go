package furtado

import (
	"github.com/google/uuid"
)

type LocalStalkerConfig struct {
	Interface string
	IsRunning bool
	BridgeID  uuid.UUID
}
