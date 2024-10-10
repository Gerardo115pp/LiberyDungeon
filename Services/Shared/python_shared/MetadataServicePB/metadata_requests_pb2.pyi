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
