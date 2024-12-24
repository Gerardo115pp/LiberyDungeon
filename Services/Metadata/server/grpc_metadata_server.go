package server

import (
	"context"
	"errors"
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

func (ms *MetadataGrpcServer) GetEntitiesWithTaggings(ctx context.Context, tag_list *metadata_service_pb.TagList) (*metadata_service_pb.EntitiesByType, error) {
	tag_ids := make([]int, len(tag_list.TagId))
	for h, tag_id := range tag_list.TagId {
		tag_ids[h] = int(tag_id)
	}

	if len(tag_ids) == 0 {
		return nil, errors.New("No tags provided")
	}

	entities, err := repository.DungeonTagsRepo.GetEntitiesWithTaggingsCTX(ctx, tag_ids)
	if err != nil {
		return nil, err
	}

	var entities_by_type map[string][]string = make(map[string][]string)
	entities_by_type_message := &metadata_service_pb.EntitiesByType{
		EntitiesByType: make(map[string]*metadata_service_pb.EntityList),
	}

	for _, entity := range entities {
		entities_of_type, list_exist := entities_by_type[entity.EntityType]

		if !list_exist {
			entities_of_type = make([]string, 0)
		}

		entities_of_type = append(entities_of_type, entity.TaggedEntityUUID)
		entities_by_type[entity.EntityType] = entities_of_type
	}

	for entity_type, entities := range entities_by_type {
		var entity_list *metadata_service_pb.EntityList = new(metadata_service_pb.EntityList)

		entity_list.EntitiesUuids = entities

		entities_by_type_message.EntitiesByType[entity_type] = entity_list
	}

	return entities_by_type_message, nil
}

func (ms *MetadataGrpcServer) TagEntities(ctx context.Context, request *metadata_service_pb.TaggableEntities) (*metadata_service_pb.BooleanResponse, error) {
	var response *metadata_service_pb.BooleanResponse = new(metadata_service_pb.BooleanResponse)

	var tag_id int = int(request.TagId)
	var entities_uuids []string = request.EntitiesUuids
	var entity_type string = request.EntityType

	err := repository.DungeonTagsRepo.TagEntitiesCTX(ctx, tag_id, entities_uuids, entity_type)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error while tagging entities: %s", err))
	}

	response.Response = err == nil

	return response, nil
}

func (ms *MetadataGrpcServer) UntagEntities(ctx context.Context, request *metadata_service_pb.TaggableEntities) (*metadata_service_pb.BooleanResponse, error) {
	var response *metadata_service_pb.BooleanResponse = new(metadata_service_pb.BooleanResponse)

	var tag_id int = int(request.TagId)
	var entities_uuids []string = request.EntitiesUuids

	err := repository.DungeonTagsRepo.RemoveTagFromEntitiesCTX(ctx, tag_id, entities_uuids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error while untagging entities: %s", err))
	}

	response.Response = err == nil

	return response, nil
}

func (ms *MetadataGrpcServer) DeleteEntitiesTaggings(ctx context.Context, request *metadata_service_pb.EntityList) (*metadata_service_pb.BooleanResponse, error) {
	var response *metadata_service_pb.BooleanResponse = new(metadata_service_pb.BooleanResponse)
	entity_uuids := request.EntitiesUuids

	response.Response = true

	err := repository.DungeonTagsRepo.RemoveAllTaggingsForEntitiesCTX(ctx, entity_uuids)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error while deleting taggings for %d entities"))
		err = errors.Join(fmt.Errorf("In server/grpc_metadata_server.DeleteEntitiesTaggings: Error while calling repository.DungeonTagsRepo.RemoveAllTaggingsForEntitiesCTX"), err)
		echo.EchoErr(err)
		response.Response = false
	}

	return response, err
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
