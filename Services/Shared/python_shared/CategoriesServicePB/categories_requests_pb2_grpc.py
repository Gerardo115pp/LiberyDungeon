# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

import CategoriesServicePB.categories_requests_pb2 as categories__requests__pb2

GRPC_GENERATED_VERSION = '1.67.1'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in categories_requests_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class CategoriesServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateCategory = channel.unary_unary(
                '/categories_service.CategoriesService/CreateCategory',
                request_serializer=categories__requests__pb2.CreateCategoryRequest.SerializeToString,
                response_deserializer=categories__requests__pb2.CreateCategoryResponse.FromString,
                _registered_method=True)
        self.GetCategory = channel.unary_unary(
                '/categories_service.CategoriesService/GetCategory',
                request_serializer=categories__requests__pb2.GetCategoryRequest.SerializeToString,
                response_deserializer=categories__requests__pb2.GetCategoryResponse.FromString,
                _registered_method=True)
        self.GetCategoriesCluster = channel.unary_unary(
                '/categories_service.CategoriesService/GetCategoriesCluster',
                request_serializer=categories__requests__pb2.GetCategoriesClusterRequest.SerializeToString,
                response_deserializer=categories__requests__pb2.GetCategoriesClusterResponse.FromString,
                _registered_method=True)


class CategoriesServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CreateCategory(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetCategory(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetCategoriesCluster(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_CategoriesServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateCategory': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateCategory,
                    request_deserializer=categories__requests__pb2.CreateCategoryRequest.FromString,
                    response_serializer=categories__requests__pb2.CreateCategoryResponse.SerializeToString,
            ),
            'GetCategory': grpc.unary_unary_rpc_method_handler(
                    servicer.GetCategory,
                    request_deserializer=categories__requests__pb2.GetCategoryRequest.FromString,
                    response_serializer=categories__requests__pb2.GetCategoryResponse.SerializeToString,
            ),
            'GetCategoriesCluster': grpc.unary_unary_rpc_method_handler(
                    servicer.GetCategoriesCluster,
                    request_deserializer=categories__requests__pb2.GetCategoriesClusterRequest.FromString,
                    response_serializer=categories__requests__pb2.GetCategoriesClusterResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'categories_service.CategoriesService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('categories_service.CategoriesService', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class CategoriesService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CreateCategory(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/categories_service.CategoriesService/CreateCategory',
            categories__requests__pb2.CreateCategoryRequest.SerializeToString,
            categories__requests__pb2.CreateCategoryResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetCategory(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/categories_service.CategoriesService/GetCategory',
            categories__requests__pb2.GetCategoryRequest.SerializeToString,
            categories__requests__pb2.GetCategoryResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def GetCategoriesCluster(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/categories_service.CategoriesService/GetCategoriesCluster',
            categories__requests__pb2.GetCategoriesClusterRequest.SerializeToString,
            categories__requests__pb2.GetCategoriesClusterResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
