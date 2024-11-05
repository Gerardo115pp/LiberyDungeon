from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class IsClusterPrivate(_message.Message):
    __slots__ = ("cluster_uuid",)
    CLUSTER_UUID_FIELD_NUMBER: _ClassVar[int]
    cluster_uuid: str
    def __init__(self, cluster_uuid: _Optional[str] = ...) -> None: ...

class TaggableEntities(_message.Message):
    __slots__ = ("entities_uuids", "tag_id", "entity_type")
    ENTITIES_UUIDS_FIELD_NUMBER: _ClassVar[int]
    TAG_ID_FIELD_NUMBER: _ClassVar[int]
    ENTITY_TYPE_FIELD_NUMBER: _ClassVar[int]
    entities_uuids: _containers.RepeatedScalarFieldContainer[str]
    tag_id: int
    entity_type: str
    def __init__(self, entities_uuids: _Optional[_Iterable[str]] = ..., tag_id: _Optional[int] = ..., entity_type: _Optional[str] = ...) -> None: ...

class BooleanResponse(_message.Message):
    __slots__ = ("response",)
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    response: bool
    def __init__(self, response: bool = ...) -> None: ...

class AllPrivateClustersResponse(_message.Message):
    __slots__ = ("private_clusters",)
    PRIVATE_CLUSTERS_FIELD_NUMBER: _ClassVar[int]
    private_clusters: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, private_clusters: _Optional[_Iterable[str]] = ...) -> None: ...
