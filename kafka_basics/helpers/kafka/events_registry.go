package kafka_helpers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const (
	EVENT_STATUS_UNKNOWN           uint = 0
	EVENT_STATUS_QUEUED            uint = 1
	EVENT_STATUS_PROCESSING        uint = 2
	EVENT_STATUS_PROCESSED         uint = 3
	EVENT_STATUS_PROCESSING_FAILED uint = 4
)

type EventsRegistry struct {
	sync.RWMutex
	registry map[uuid.UUID]uint
}

func NewEventsRegistry() *EventsRegistry {
	return &EventsRegistry{
		registry: map[uuid.UUID]uint{},
	}
}

func (er *EventsRegistry) Set(eventID uuid.UUID, status uint) {
	er.Lock()
	er.registry[eventID] = status
	defer er.Unlock()
}

func (er *EventsRegistry) Add(eventID uuid.UUID) {
	er.Lock()
	er.registry[eventID] = EVENT_STATUS_QUEUED
	defer er.Unlock()
}

func (er *EventsRegistry) Get(eventID uuid.UUID) (uint, error) {
	er.RLock()
	if retVal, ok := er.registry[eventID]; ok {
		defer er.RUnlock()
		return retVal, nil
	}
	defer er.RUnlock()
	return EVENT_STATUS_UNKNOWN, errors.New("no event found within the registry")
}

func (er *EventsRegistry) Exists(eventID uuid.UUID) bool {
	er.RLock()
	_, ok := er.registry[eventID]
	defer er.RUnlock()
	return ok
}

func (er *EventsRegistry) String() string {
	er.RLock()
	retVal := ""
	for k, v := range er.registry {
		retVal += k.String() + ":" + fmt.Sprint(v) + ", "
	}
	defer er.RUnlock()
	return retVal
}
