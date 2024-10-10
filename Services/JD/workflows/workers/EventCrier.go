package workers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	dungeon_models "libery-dungeon-libs/models"
	"libery-dungeon-libs/models/collections"
	"libery_JD_service/repository"
	"libery_JD_service/workflows/jobs"
	"net/http"
	"time"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/gorilla/websocket"
)

type EventCrier struct {
	announcement_queue   *collections.Queue[string]
	announcement_arrived chan bool
	listeners_upgrader   *websocket.Upgrader
	event_listeners      []*websocket.Conn
}

func (ec *EventCrier) GetUpgrader() *websocket.Upgrader {
	return ec.listeners_upgrader
}

func (ec *EventCrier) Close() {
	echo.EchoDebug("Closing EventCrier")
	ec.announcement_arrived <- false

	for _, listener := range ec.event_listeners {
		echo.EchoDebug(fmt.Sprintf("Closing listener: %s", listener.RemoteAddr().String()))
		listener.Close()
	}
}

func (ec *EventCrier) cleanUpConn(conn *websocket.Conn) {
	echo.EchoDebug(fmt.Sprintf("Cleaning up connection: %s. Current listeners: %d", conn.RemoteAddr().String(), len(ec.event_listeners)))

	for h, listener := range ec.event_listeners {
		if listener == conn {
			ec.event_listeners = append(ec.event_listeners[:h], ec.event_listeners[h+1:]...)
			break
		}

		conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
	}

	echo.EchoDebug(fmt.Sprintf("Connection cleaned up. Current listeners: %d", len(ec.event_listeners)))
}

func (ec *EventCrier) monitorAnnouncementQueue() {
	var new_event_uuid_received bool

	for {
		if ec.announcement_queue.IsEmpty() {
			select {
			case new_event_uuid_received = <-ec.announcement_arrived:
				if !new_event_uuid_received {
					return
				}
			}
		}

		new_event_uuid, err := ec.announcement_queue.Dequeue()
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error dequeuing announcement: %s", err.Error()))
			continue
		}

		ec.oyez(*new_event_uuid) // announce the event
	}
}

// In English-speaking countries, a town crier carried a handbell to attract people's attention,
// as they shouted the words "Oyez, Oyez, Oyez!" before making their announcements.
func (ec *EventCrier) oyez(event_uuid string) {
	if len(ec.event_listeners) == 0 {
		return
	}

	var announcement_event *communication.PlatformEvent
	var lerr *dungeon_models.LabeledError

	announcement_event, lerr = repository.PlatformEvents.GetPublicEvent(event_uuid)
	if lerr != nil {
		lerr.AppendContext("EventCrier.oyez while getting public event")
		echo.Echo(echo.RedFG, lerr.Error())
		return
	}

	announcement_message, err := json.Marshal(announcement_event)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error marshalling announcement event: %s", err.Error()))
		return
	}

	for _, listener := range ec.event_listeners {
		echo.EchoDebug(fmt.Sprintf("Announcing event to listener: %s", listener.RemoteAddr().String()))

		err = listener.WriteMessage(websocket.TextMessage, announcement_message)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("Error writing announcement message to listener: %s", err.Error()))
			defer ec.cleanUpConn(listener)
			continue
		}
	}

	echo.Echo(echo.GreenFG, fmt.Sprintf("Announced event to %d listeners", len(ec.event_listeners)))
}

func (ec *EventCrier) RegisterListener(new_listener *websocket.Conn) {
	// new_listener.SetWriteDeadline(time.Now().Add(10 * time.Second))

	ec.event_listeners = append(ec.event_listeners, new_listener)
}

func (ec *EventCrier) RegisterEvent(new_event communication.PlatformEvent) error {
	// Event crier only announces public events
	if !communication.IsPublicEvent(new_event.EventType) {
		return fmt.Errorf("Event crier only announces public events")
	}

	repository.PlatformEvents.RegisterEvent(&new_event)

	ec.announcement_queue.Enqueue(&new_event.Uuid)

	ec.announcement_arrived <- true

	return nil
}

func (ec *EventCrier) GetCloser() jobs.Closer {
	return ec.Close
}

func NewEventCrier() *EventCrier {
	var new_event_crier *EventCrier = new(EventCrier)

	new_event_crier.announcement_queue = &collections.Queue[string]{}
	new_event_crier.announcement_arrived = make(chan bool)
	new_event_crier.event_listeners = make([]*websocket.Conn, 0)
	new_event_crier.listeners_upgrader = &websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	go new_event_crier.monitorAnnouncementQueue()

	return new_event_crier
}
