package repository

import (
	"libery-dungeon-libs/communication"
	dungeon_models "libery-dungeon-libs/models"
)

type PlatformEventsRegistry interface {
	RegisterEvent(event *communication.PlatformEvent)
	GetPublicEvent(event_uuid string) (*communication.PlatformEvent, *dungeon_models.LabeledError)
	GetPrivateEvent(event_uuid string) (*communication.PlatformEvent, *dungeon_models.LabeledError)
}

var PlatformEvents PlatformEventsRegistry

func SetPlatformEventsRepo(impl PlatformEventsRegistry) {
	PlatformEvents = impl
}
