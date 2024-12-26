package service_clients

import (
	"context"
	"errors"
	"fmt"
	"libery-dungeon-libs/metadata_service_pb"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetadataServiceClient struct {
	BaseServiceClient
}

func (metadata_client MetadataServiceClient) Alive() (bool, error) {
	var metadata_endpoint string

	metadata_endpoint = metadata_client.getHttpsEndpoint()

	metadata_endpoint += "/alive"

	request, err := http.NewRequest("GET", metadata_endpoint, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Transport: metadata_client.HttpTransport,
	}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return response.StatusCode >= 200 && response.StatusCode < 300, nil
}

func (metadata_client MetadataServiceClient) getHttpsEndpoint() string {
	return fmt.Sprintf("https://%s%s", metadata_client.BaseDomain, metadata_client.HttpAddress)
}

func (metadata_client MetadataServiceClient) CheckClusterPrivate(cluster_uuid string) (bool, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.IsClusterPrivate{
		ClusterUuid: cluster_uuid,
	}

	boolean_response, err := metadata_grpc_client.CheckClusterPrivate(ctx, &message)

	return boolean_response.Response, err
}

func (metadata_client MetadataServiceClient) CopyEntityTagsToEntities(source_entity, cluster_domain, entities_type string, entities_uuids []string) (bool, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.CopyEntityTags{
		SourceEntity:  source_entity,
		ClusterDomain: cluster_domain,
		EntitiesType:  entities_type,
		Entities:      entities_uuids,
	}

	response, err := metadata_grpc_client.CopyEntityTagsToEntityList(ctx, &message)
	if err != nil {
		return false, errors.Join(err, fmt.Errorf("In Communication/MetadataService.CopyEntityTagsToEntities, while calling metadata_grpc_client.CopyEntityTagsToEntityList"))
	}

	return response.Response, nil
}

func (metadata_client MetadataServiceClient) GetAllPrivateClusters() ([]string, error) {
	var private_clusters []string = make([]string, 0)

	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return private_clusters, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	clusters_response, err := metadata_grpc_client.GetAllPrivateClusters(ctx, &emptypb.Empty{})
	if err != nil {
		return private_clusters, err
	}

	private_clusters = clusters_response.PrivateClusters

	return private_clusters, nil
}

func (metadata_client MetadataServiceClient) GetEntitiesWithTaggings(tag_list []int) (entities_by_type map[string][]string, err error) {
	entities_by_type = make(map[string][]string)

	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tag_list32 := make([]int32, len(tag_list))
	for h, v := range tag_list {
		tag_list32[h] = int32(v)
	}

	message := metadata_service_pb.TagList{
		TagId: tag_list32,
	}

	entities_by_type_response, err := metadata_grpc_client.GetEntitiesWithTaggings(ctx, &message)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("In Communication/MetadataService.GetEntitiesWithTaggings, while calling metadata_grpc_client.GetEntitiesWithTaggings"))
	}

	for entity_type, entities_list := range entities_by_type_response.EntitiesByType {
		entities_by_type[entity_type] = entities_list.EntitiesUuids
	}

	return
}

func (metadata_client MetadataServiceClient) GetEntityTags(entity_uuid, cluster_domain string) ([]int, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.Entity{
		EntityUuid:    entity_uuid,
		ClusterDomain: cluster_domain,
	}

	response, err := metadata_grpc_client.GetEntityTags(ctx, &message)
	if err != nil {
		return nil, err
	}

	tag_list := make([]int, len(response.TagId))
	for h, tag_id := range response.TagId {
		tag_list[h] = int(tag_id)
	}

	return tag_list, err
}

func (metadata_client MetadataServiceClient) TagEntities(tag_id int, entities_uuids []string, entity_type string) (bool, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.TaggableEntities{
		TagId:         int32(tag_id),
		EntitiesUuids: entities_uuids,
		EntityType:    entity_type,
	}

	boolean_response, err := metadata_grpc_client.TagEntities(ctx, &message)

	return boolean_response.Response, err
}

func (metadata_client MetadataServiceClient) UntagEntities(tag_id int, entities_uuids []string) (bool, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.TaggableEntities{
		TagId:         int32(tag_id),
		EntitiesUuids: entities_uuids,
	}

	boolean_response, err := metadata_grpc_client.UntagEntities(ctx, &message)

	return boolean_response.Response, err
}

func (metadata_client MetadataServiceClient) RemoveAllTaggingsForEntities(entities []string) (bool, error) {
	conn, err := grpc.Dial(metadata_client.GrpcAddress, grpc.WithTransportCredentials(metadata_client.GrpcTransport))
	if err != nil {
		return false, err
	}
	defer conn.Close()

	metadata_grpc_client := metadata_service_pb.NewMetadataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := metadata_service_pb.EntityList{
		EntitiesUuids: entities,
	}

	boolean_response, err := metadata_grpc_client.DeleteEntitiesTaggings(ctx, &message)

	return boolean_response.Response, err
}
