package furtado

import (
	"github.com/google/uuid"
)

type IStalkerConfig interface {
	RemoteStalkerConfig | LocalStalkerConfig
}

type RemoteStalkerConfig struct {
	Interface string
	IsRunning bool
	BridgeIP  string
	Port      int
}

type LocalStalkerConfig struct {
	Interface string
	IsRunning bool
	BridgeID  uuid.UUID
}