package jobs

import (
	"libery-dungeon-libs/communication"

	"github.com/gorilla/websocket"
)

type Closer func()

type EventTownCrier interface {
	RegisterEvent(new_event communication.PlatformEvent) error
	RegisterListener(ws *websocket.Conn)
	GetUpgrader() *websocket.Upgrader
	GetCloser() Closer
	Close()
}

var TownCrier EventTownCrier

func SetTownCrierImplementation(impl EventTownCrier) {
	TownCrier = impl
}
