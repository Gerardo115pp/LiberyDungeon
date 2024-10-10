package server

import (
	"context"
	"fmt"
	"libery-dungeon-libs/metadata_service_pb"
	"libery-metadata-service/repository"
	"net"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetadataGrpcServer struct {
	srv  *grpc.Server
	port string
	metadata_service_pb.UnimplementedMetadataServiceServer
}

func NewMetadataGrpcServer(service_port string) (*MetadataGrpcServer, error) {
	return &MetadataGrpcServer{
		port: service_port,
	}, nil
}

func (ms *MetadataGrpcServer) CheckClusterPrivate(ctx context.Context, request *metadata_service_pb.IsClusterPrivate) (*metadata_service_pb.BooleanResponse, error) {
	var response *metadata_service_pb.BooleanResponse = new(metadata_service_pb.BooleanResponse)
	echo.Echo(echo.BlueFG, fmt.Sprintf("Checking if cluster is private: %s", request.ClusterUuid))

	var is_private bool = repository.ClusterMetadataRepo.IsPrivateCluster(request.ClusterUuid)

	response.Response = is_private

	return response, nil
}

func (ms *MetadataGrpcServer) GetAllPrivateClusters(ctx context.Context, duh *emptypb.Empty) (*metadata_service_pb.AllPrivateClustersResponse, error) {
	var response *metadata_service_pb.AllPrivateClustersResponse = new(metadata_service_pb.AllPrivateClustersResponse)
	echo.Echo(echo.BlueFG, "Getting all private clusters")

	var private_cluster_uuids []string = repository.ClusterMetadataRepo.GetPrivateClusters()

	response.PrivateClusters = private_cluster_uuids

	return response, nil
}

func (ms *MetadataGrpcServer) Connect() error {
	listener, err := net.Listen("tcp", ms.port)
	if err != nil {
		return err
	}

	ms.srv = grpc.NewServer()
	metadata_service_pb.RegisterMetadataServiceServer(ms.srv, ms)

	service_name := metadata_service_pb.MetadataService_ServiceDesc.ServiceName

	reflection.Register(ms.srv)

	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting Metadata GRPC Server  for service '%s' on port %s", service_name, ms.port))
	if err := ms.srv.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (ms *MetadataGrpcServer) Shutdown() error {
	if ms.srv == nil {
		return nil
	}

	ms.srv.GracefulStop()

	return nil
}
