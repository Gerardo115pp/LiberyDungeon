from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class ServiceOnlineNotification(_message.Message):
    __slots__ = ("service_name", "service_route", "service_port")
    SERVICE_NAME_FIELD_NUMBER: _ClassVar[int]
    SERVICE_ROUTE_FIELD_NUMBER: _ClassVar[int]
    SERVICE_PORT_FIELD_NUMBER: _ClassVar[int]
    service_name: str
    service_route: str
    service_port: str
    def __init__(self, service_name: _Optional[str] = ..., service_route: _Optional[str] = ..., service_port: _Optional[str] = ...) -> None: ...

class ServiceOfflineNotification(_message.Message):
    __slots__ = ("service_name",)
    SERVICE_NAME_FIELD_NUMBER: _ClassVar[int]
    service_name: str
    def __init__(self, service_name: _Optional[str] = ...) -> None: ...

class PlatformEvent(_message.Message):
    __slots__ = ("uuid", "event_type", "event_message", "event_payload")
    UUID_FIELD_NUMBER: _ClassVar[int]
    EVENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    EVENT_MESSAGE_FIELD_NUMBER: _ClassVar[int]
    EVENT_PAYLOAD_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    event_type: str
    event_message: str
    event_payload: str
    def __init__(self, uuid: _Optional[str] = ..., event_type: _Optional[str] = ..., event_message: _Optional[str] = ..., event_payload: _Optional[str] = ...) -> None: ...
