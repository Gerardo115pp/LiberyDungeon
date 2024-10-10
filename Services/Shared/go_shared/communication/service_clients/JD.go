package service_clients

import (
	"context"
	"fmt"
	JD_service_pb "libery-dungeon-libs/jd_service_pb"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type JD_Client struct {
	BaseServiceClient
}

func (jd_conf JD_Client) Alive() (bool, error) {
	var jd_endpoint string

	jd_endpoint = jd_conf.getHttpsEndpoint()

	jd_endpoint += "/alive"

	request, err := http.NewRequest("GET", jd_endpoint, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{
		Transport: jd_conf.HttpTransport,
	}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return response.StatusCode >= 200 && response.StatusCode < 300, nil
}

func (jd_conf JD_Client) getHttpsEndpoint() string {
	return fmt.Sprintf("https://%s%s", jd_conf.BaseDomain, jd_conf.HttpAddress)
}

func (jd_conf JD_Client) EmitPlatformEvent(event_uuid, event_type, event_message, event_payload string) error {
	conn, err := grpc.Dial(jd_conf.GrpcAddress, grpc.WithTransportCredentials(jd_conf.GrpcTransport))
	if err != nil {
		return err
	}
	defer conn.Close()

	jd_grpc_client := JD_service_pb.NewJDServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := JD_service_pb.PlatformEvent{
		Uuid:         event_uuid,
		EventType:    event_type,
		EventMessage: event_message,
		EventPayload: event_payload,
	}

	_, err = jd_grpc_client.EmitPlatformEvent(ctx, &message)

	return err
}

// TODO: sign in messages should be signed.

func (jd_conf JD_Client) NotifyServiceOnline(service_name, service_route, service_port string) error {
	conn, err := grpc.Dial(jd_conf.GrpcAddress, grpc.WithTransportCredentials(jd_conf.GrpcTransport))
	if err != nil {
		return err
	}
	defer conn.Close()

	jd_grpc_client := JD_service_pb.NewJDServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := JD_service_pb.ServiceOnlineNotification{
		ServiceName:  service_name,
		ServiceRoute: service_route,
		ServicePort:  service_port,
	}

	_, err = jd_grpc_client.NotifyServiceOnline(ctx, &message)

	return err
}

func (jd_conf JD_Client) NotifyServiceOffline(service_name string) error {
	conn, err := grpc.Dial(jd_conf.GrpcAddress, grpc.WithTransportCredentials(jd_conf.GrpcTransport))
	if err != nil {
		return err
	}
	defer conn.Close()

	jd_grpc_client := JD_service_pb.NewJDServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	message := JD_service_pb.ServiceOfflineNotification{
		ServiceName: service_name,
	}

	_, err = jd_grpc_client.NotifyServiceOffline(ctx, &message)

	return err
}
