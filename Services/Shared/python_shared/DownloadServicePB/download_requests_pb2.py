# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: download_requests.proto
# Protobuf Python Version: 5.28.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    3,
    '',
    'download_requests.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x17\x64ownload_requests.proto\x12\x10\x64ownload_service\"\x8c\x01\n\x1a\x44ownloadImagesBatchRequest\x12\x12\n\nimage_urls\x18\x01 \x03(\t\x12\x15\n\rcategory_uuid\x18\x02 \x01(\t\x12\x15\n\rcluster_token\x18\x03 \x01(\t\x12\x1a\n\rdownload_uuid\x18\x04 \x01(\tH\x00\x88\x01\x01\x42\x10\n\x0e_download_uuid\".\n\x15\x44ownloadBatchResponse\x12\x15\n\rdownload_uuid\x18\x01 \x01(\t2\x7f\n\x0f\x44ownloadService\x12l\n\x13\x44ownloadImagesBatch\x12,.download_service.DownloadImagesBatchRequest\x1a\'.download_service.DownloadBatchResponseBVZTgithub.com/Gerardo115pp/libery-dungeon/libery_downloads_service;downloads_service_pbb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'download_requests_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'ZTgithub.com/Gerardo115pp/libery-dungeon/libery_downloads_service;downloads_service_pb'
  _globals['_DOWNLOADIMAGESBATCHREQUEST']._serialized_start=46
  _globals['_DOWNLOADIMAGESBATCHREQUEST']._serialized_end=186
  _globals['_DOWNLOADBATCHRESPONSE']._serialized_start=188
  _globals['_DOWNLOADBATCHRESPONSE']._serialized_end=234
  _globals['_DOWNLOADSERVICE']._serialized_start=236
  _globals['_DOWNLOADSERVICE']._serialized_end=363
# @@protoc_insertion_point(module_scope)
