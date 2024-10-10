package databases

import (
	"fmt"
	dungeon_models "libery-dungeon-libs/models"

	// "libery-dungeon-libs/communication"
	"libery-dungeon-libs/communication"
	service_models "libery_JD_service/models"
)

type PlatformEventsStore struct {
	PublicEvents  map[string]*communication.PlatformEvent
	privateEvents map[string]*communication.PlatformEvent
}

func NewPlatformEventsStore() *PlatformEventsStore {
	return &PlatformEventsStore{
		PublicEvents:  make(map[string]*communication.PlatformEvent),
		privateEvents: make(map[string]*communication.PlatformEvent),
	}
}

func (pes *PlatformEventsStore) RegisterEvent(event *communication.PlatformEvent) {
	if communication.IsPublicEvent(event.EventType) {
		pes.PublicEvents[event.Uuid] = event
		return
	}

	pes.privateEvents[event.Uuid] = event
}

func (pes PlatformEventsStore) GetPublicEvent(event_uuid string) (*communication.PlatformEvent, *dungeon_models.LabeledError) {
	var labeled_error *dungeon_models.LabeledError

	event, exists := pes.PublicEvents[event_uuid]

	if !exists {
		labeled_error = dungeon_models.NewLabeledError(fmt.Errorf("Event with UUID %s not found", event_uuid), "In PlatformEventsStore.GetPublicEvent", service_models.ErrService_NoSuchEventRegistered)
	}

	return event, labeled_error
}

func (pes PlatformEventsStore) GetPrivateEvent(event_uuid string) (*communication.PlatformEvent, *dungeon_models.LabeledError) {
	var labeled_error *dungeon_models.LabeledError

	event, exists := pes.privateEvents[event_uuid]

	if !exists {
		labeled_error = dungeon_models.NewLabeledError(fmt.Errorf("Event with UUID %s not found", event_uuid), "In PlatformEventsStore.GetPrivateEvent", service_models.ErrService_NoSuchEventRegistered)
	}

	return event, labeled_error
}
