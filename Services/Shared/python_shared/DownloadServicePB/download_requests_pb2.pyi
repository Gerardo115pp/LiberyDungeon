from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class DownloadImagesBatchRequest(_message.Message):
    __slots__ = ("image_urls", "category_uuid", "cluster_token", "download_uuid")
    IMAGE_URLS_FIELD_NUMBER: _ClassVar[int]
    CATEGORY_UUID_FIELD_NUMBER: _ClassVar[int]
    CLUSTER_TOKEN_FIELD_NUMBER: _ClassVar[int]
    DOWNLOAD_UUID_FIELD_NUMBER: _ClassVar[int]
    image_urls: _containers.RepeatedScalarFieldContainer[str]
    category_uuid: str
    cluster_token: str
    download_uuid: str
    def __init__(self, image_urls: _Optional[_Iterable[str]] = ..., category_uuid: _Optional[str] = ..., cluster_token: _Optional[str] = ..., download_uuid: _Optional[str] = ...) -> None: ...

class DownloadBatchResponse(_message.Message):
    __slots__ = ("download_uuid",)
    DOWNLOAD_UUID_FIELD_NUMBER: _ClassVar[int]
    download_uuid: str
    def __init__(self, download_uuid: _Optional[str] = ...) -> None: ...
