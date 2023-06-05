package stalker

import (
	"github.com/google/uuid"
)

type IStalkerOrchestrator[TConfig IStalkerConfig] interface {
	CreateStalker(TConfig) uuid.UUID
	DeleteStalker(uuid.UUID) error
	PauseStalker(uuid.UUID) error
	StartStalker(uuid.UUID) error
}
