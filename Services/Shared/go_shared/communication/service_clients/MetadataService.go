package service_clients

import (
	"context"
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
