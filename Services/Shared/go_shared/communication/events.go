package communication

import (
	"fmt"
	dungeon_models "libery-dungeon-libs/models"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// This is a list of well-known platform events
// upon which we take specific actions. Other names are also valid
// but JD will only emit them and store them, no further action will be taken.
const (
	PlatformEvent_Private_SystemLog        = "system_log"
	PlatformEvent_Private_FileSystemChange = "file_system_change"
	PlatformEvent_Public_ClusterFSChange   = "cluster_fs_change"
	PlatformEvent_Public_MediaDeleted      = "media_deleted"
	PlatformEvent_Public_MediaAdded        = "media_added"
)

var public_events = [...]string{
	PlatformEvent_Public_ClusterFSChange,
	PlatformEvent_Public_MediaDeleted,
	PlatformEvent_Public_MediaAdded,
}

var private_events = [...]string{
	PlatformEvent_Private_SystemLog,
	PlatformEvent_Private_FileSystemChange,
}

type PlatformEvent struct {
	Uuid         string `json:"uuid"`          // Some endpoints require async processing, in which case they will return a UUID to track the event's status
	EventType    string `json:"event_type"`    // Identify the type of effect this event has
	EventMessage string `json:"event_message"` // Human readable message
	EventPayload string `json:"event_payload"` // JWT token
}

func NewPlatformEvent(event_uuid, event_type, event_message, event_payload string) *PlatformEvent {
	return &PlatformEvent{
		Uuid:         event_uuid,
		EventType:    event_type,
		EventMessage: event_message,
		EventPayload: event_payload,
	}
}

func (pe PlatformEvent) IsAuthenticated(jwt_sk string) bool {
	token, err := jwt.Parse(pe.EventPayload, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_sk), nil
	})
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In PlatformEvent.IsAuthenticated: While parsing JWT token: %s", err.Error()))
		return false
	}

	return token.Valid
}

func (pe PlatformEvent) Emit() (err error) {
	err = JD.EmitPlatformEvent(
		pe.Uuid,
		pe.EventType,
		pe.EventMessage,
		pe.EventPayload,
	)

	return
}

func IsWellKnownEvent(event string) bool {
	if is_public := IsPublicEvent(event); is_public {
		return true
	}

	if is_private := IsPrivateEvent(event); is_private {
		return true
	}

	return false
}

func IsPublicEvent(event_type string) bool {
	var is_public bool = false

	for _, public_event := range public_events {
		if public_event == event_type {
			is_public = true
			break
		}
	}

	return is_public
}

func IsPrivateEvent(event_type string) bool {
	var is_private bool = false

	for _, private_event := range private_events {
		if private_event == event_type {
			is_private = true
			break
		}
	}

	return is_private
}

func GenerateEventUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("In GenerateEventUUID: %s", err.Error())
	}

	return uuid.String(), nil
}

// ----------------- Well-known event builders -----------------

type ClusterFSChangePayload struct {
	ClusterUUID        string `json:"cluster_uuid"`
	MediasDeletedCount int    `json:"medias_deleted"`
	MediasAddedCount   int    `json:"medias_added"`
	MediasMovedCount   int    `json:"medias_moved"`
	jwt.StandardClaims
}

func (cfscp ClusterFSChangePayload) SignPayload(sk string) (string, error) {
	token := jwt.NewWithClaims(dungeon_models.JwtSigningMethod, cfscp)
	return token.SignedString([]byte(sk))
}

func NewClusterFSChangeEvent(sk, cluster_uuid string, medias_deleted, medias_added, medias_moved int) *PlatformEvent {
	payload := &ClusterFSChangePayload{
		ClusterUUID:        cluster_uuid,
		MediasDeletedCount: medias_deleted,
		MediasAddedCount:   medias_added,
		MediasMovedCount:   medias_moved,
	}

	singed_payload, err := payload.SignPayload(sk)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In NewClusterFSChangeEvent: %s", err.Error()))
		return nil
	}

	event_uuid, err := GenerateEventUUID()
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In NewClusterFSChangeEvent: %s", err.Error()))
		return nil
	}

	event_message := fmt.Sprintf("Cluster '%s'changed. Medias deleted: %d, added: %d, moved: %d", cluster_uuid, medias_deleted, medias_added, medias_moved)

	return NewPlatformEvent(event_uuid, PlatformEvent_Public_ClusterFSChange, event_message, singed_payload)
}
