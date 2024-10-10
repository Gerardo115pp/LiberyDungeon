package server

import (
	"context"
	"fmt"
	"libery-dungeon-libs/communication"
	JD_service_pb "libery-dungeon-libs/jd_service_pb"
	"libery-dungeon-libs/models/platform_services"
	app_config "libery_JD_service/Config"
	service_models "libery_JD_service/models"
	"libery_JD_service/repository"
	"libery_JD_service/workflows/jobs"
	"net"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type JDGrpcServer struct {
	srv  *grpc.Server
	port string
	JD_service_pb.UnimplementedJDServiceServer
}

func NewJDGrpcServer(service_port string) (*JDGrpcServer, error) {
	return &JDGrpcServer{
		port: service_port,
	}, nil
}

func (s *JDGrpcServer) EmitPlatformEvent(ctx context.Context, new_event *JD_service_pb.PlatformEvent) (*emptypb.Empty, error) {
	var platform_event *communication.PlatformEvent = communication.NewPlatformEvent(new_event.Uuid, new_event.EventType, new_event.EventMessage, new_event.EventPayload)

	if !platform_event.IsAuthenticated(app_config.JWT_SECRET) {
		echo.Echo(echo.RedFG, fmt.Sprintf("Event %s is not authenticated", new_event.Uuid))
		return nil, fmt.Errorf("Invalid request")
	}

	err := jobs.TownCrier.RegisterEvent(*platform_event)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("Error registering event: %s", err.Error()))
		return nil, fmt.Errorf("Error registering event")
	}

	echo.EchoDebug(fmt.Sprintf("Event '%s' registered", new_event.Uuid))

	return new(emptypb.Empty), nil
}

func (s *JDGrpcServer) NotifyServiceOnline(ctx context.Context, request *JD_service_pb.ServiceOnlineNotification) (result *emptypb.Empty, err error) {
	result = new(emptypb.Empty)

	var service_name platform_services.PlatformServiceName
	if !platform_services.ServiceNameValid(request.ServiceName) {
		echo.Echo(echo.RedFG, fmt.Sprintf("Service name '%s' does not match any existing service", request.ServiceName))
		return nil, fmt.Errorf("Invalid request")
	}
	service_name = platform_services.PlatformServiceName(request.ServiceName)

	new_service := service_models.NewPlatformService(service_name, request.ServiceRoute, request.ServicePort)
	new_service.IsOnline = true

	err = repository.PlatformFeaturesRepo.SetServiceMetadata(*new_service)

	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Service %s is online", request.ServiceName))

	return result, nil
}

func (s *JDGrpcServer) NotifyServiceOffline(ctx context.Context, request *JD_service_pb.ServiceOfflineNotification) (result *emptypb.Empty, err error) {
	result = new(emptypb.Empty)

	if !platform_services.ServiceNameValid(request.ServiceName) {
		echo.Echo(echo.RedFG, fmt.Sprintf("Service name '%s' does not match any existing service", request.ServiceName))
		return nil, fmt.Errorf("Invalid request")
	}

	err = repository.PlatformFeaturesRepo.SetServiceOnline(request.ServiceName, false)

	echo.Echo(echo.SkyBlueFG, fmt.Sprintf("Service %s is offline", request.ServiceName))

	return result, nil
}

func (s *JDGrpcServer) Connect() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	JD_service_pb.RegisterJDServiceServer(s.srv, s)

	service_name := JD_service_pb.JDService_ServiceDesc.ServiceName

	reflection.Register(s.srv)

	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting JD GRPC Server for service '%s' on port %s", service_name, s.port))
	if err := s.srv.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *JDGrpcServer) Shutdown() error {
	if s.srv == nil {
		return nil
	}

	s.srv.GracefulStop()

	return nil
}
