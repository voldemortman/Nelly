package stalker

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type IStalkerOrchestrator[TConfig IStalkerConfig] interface {
	CreateStalker(TConfig) uuid.UUID
	DeleteStalker(uuid.UUID) error
	PauseStalker(uuid.UUID) error
	StartStalker(uuid.UUID) error
}

type RemoteStalkerOrchestrator struct {
	stalkerMap *map[uuid.UUID]remoteStalker
}

func (orchestrator *RemoteStalkerOrchestrator) CreateStalker(config RemoteStalkerConfig) uuid.UUID {
	quitChannel := make(chan struct{})
	stalker := &remoteStalker{config, quitChannel}
	id := uuid.New()
	(*orchestrator.stalkerMap)[id] = *stalker
	return id
}

func (orchestrator *RemoteStalkerOrchestrator) DeleteStalker(id uuid.UUID) error {
	stalker, ok := (*orchestrator.stalkerMap)[id]
	if !ok {
		return errors.New(fmt.Sprint("Stalker with id: ", id, " does not exist"))
	}
	if stalker.IsRunning() {
		stalker.StopStalking()
	}
	delete(*orchestrator.stalkerMap, id)
	return nil
}

func (orchestrator *RemoteStalkerOrchestrator) PauseStalker(id uuid.UUID) error {
	stalker, ok := (*orchestrator.stalkerMap)[id]
	if !ok {
		return errors.New(fmt.Sprint("Stalker with id: ", id, " does not exist"))
	}
	if stalker.IsRunning() {
		stalker.StopStalking()
	}
	return nil
}

func (orchestrator *RemoteStalkerOrchestrator) StartStalker(id uuid.UUID) error {
	stalker, ok := (*orchestrator.stalkerMap)[id]
	if !ok {
		return errors.New(fmt.Sprint("Stalker with id: ", id, " does not exist"))
	}
	if !stalker.IsRunning() {
		stalker.StartStalking()
	}
	return nil
}
