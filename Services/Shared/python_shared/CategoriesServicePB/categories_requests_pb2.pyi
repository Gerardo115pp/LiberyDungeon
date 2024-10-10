from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Category(_message.Message):
    __slots__ = ("uuid", "name", "fullpath", "parent")
    UUID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    FULLPATH_FIELD_NUMBER: _ClassVar[int]
    PARENT_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    name: str
    fullpath: str
    parent: str
    def __init__(self, uuid: _Optional[str] = ..., name: _Optional[str] = ..., fullpath: _Optional[str] = ..., parent: _Optional[str] = ...) -> None: ...

class CategoriesCluster(_message.Message):
    __slots__ = ("uuid", "name", "fs_path", "filter_category", "root_category")
    UUID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    FS_PATH_FIELD_NUMBER: _ClassVar[int]
    FILTER_CATEGORY_FIELD_NUMBER: _ClassVar[int]
    ROOT_CATEGORY_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    name: str
    fs_path: str
    filter_category: str
    root_category: str
    def __init__(self, uuid: _Optional[str] = ..., name: _Optional[str] = ..., fs_path: _Optional[str] = ..., filter_category: _Optional[str] = ..., root_category: _Optional[str] = ...) -> None: ...

class CreateCategoryRequest(_message.Message):
    __slots__ = ("name", "parent", "cluster")
    NAME_FIELD_NUMBER: _ClassVar[int]
    PARENT_FIELD_NUMBER: _ClassVar[int]
    CLUSTER_FIELD_NUMBER: _ClassVar[int]
    name: str
    parent: str
    cluster: str
    def __init__(self, name: _Optional[str] = ..., parent: _Optional[str] = ..., cluster: _Optional[str] = ...) -> None: ...

class CreateCategoryResponse(_message.Message):
    __slots__ = ("uuid",)
    UUID_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    def __init__(self, uuid: _Optional[str] = ...) -> None: ...

class GetCategoryRequest(_message.Message):
    __slots__ = ("uuid",)
    UUID_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    def __init__(self, uuid: _Optional[str] = ...) -> None: ...

class GetCategoriesClusterRequest(_message.Message):
    __slots__ = ("uuid",)
    UUID_FIELD_NUMBER: _ClassVar[int]
    uuid: str
    def __init__(self, uuid: _Optional[str] = ...) -> None: ...

class GetCategoryResponse(_message.Message):
    __slots__ = ("category",)
    CATEGORY_FIELD_NUMBER: _ClassVar[int]
    category: Category
    def __init__(self, category: _Optional[_Union[Category, _Mapping]] = ...) -> None: ...

class GetCategoriesClusterResponse(_message.Message):
    __slots__ = ("cluster",)
    CLUSTER_FIELD_NUMBER: _ClassVar[int]
    cluster: CategoriesCluster
    def __init__(self, cluster: _Optional[_Union[CategoriesCluster, _Mapping]] = ...) -> None: ...
