from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

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

class TagList(_message.Message):
    __slots__ = ("tag_id",)
    TAG_ID_FIELD_NUMBER: _ClassVar[int]
    tag_id: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, tag_id: _Optional[_Iterable[int]] = ...) -> None: ...

class EntitiesByType(_message.Message):
    __slots__ = ("entities_by_type",)
    class EntitiesByTypeEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: EntityList
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[EntityList, _Mapping]] = ...) -> None: ...
    ENTITIES_BY_TYPE_FIELD_NUMBER: _ClassVar[int]
    entities_by_type: _containers.MessageMap[str, EntityList]
    def __init__(self, entities_by_type: _Optional[_Mapping[str, EntityList]] = ...) -> None: ...

class EntityList(_message.Message):
    __slots__ = ("entities_uuids",)
    ENTITIES_UUIDS_FIELD_NUMBER: _ClassVar[int]
    entities_uuids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, entities_uuids: _Optional[_Iterable[str]] = ...) -> None: ...

class Entity(_message.Message):
    __slots__ = ("entity_uuid", "cluster_domain")
    ENTITY_UUID_FIELD_NUMBER: _ClassVar[int]
    CLUSTER_DOMAIN_FIELD_NUMBER: _ClassVar[int]
    entity_uuid: str
    cluster_domain: str
    def __init__(self, entity_uuid: _Optional[str] = ..., cluster_domain: _Optional[str] = ...) -> None: ...

class CopyEntityTags(_message.Message):
    __slots__ = ("source_entity", "entities", "cluster_domain", "entities_type")
    SOURCE_ENTITY_FIELD_NUMBER: _ClassVar[int]
    ENTITIES_FIELD_NUMBER: _ClassVar[int]
    CLUSTER_DOMAIN_FIELD_NUMBER: _ClassVar[int]
    ENTITIES_TYPE_FIELD_NUMBER: _ClassVar[int]
    source_entity: str
    entities: _containers.RepeatedScalarFieldContainer[str]
    cluster_domain: str
    entities_type: str
    def __init__(self, source_entity: _Optional[str] = ..., entities: _Optional[_Iterable[str]] = ..., cluster_domain: _Optional[str] = ..., entities_type: _Optional[str] = ...) -> None: ...
