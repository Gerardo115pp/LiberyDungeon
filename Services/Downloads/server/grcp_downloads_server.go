package server

import (
	"context"
	"fmt"
	downloads_service_pb "libery-dungeon-libs/download_service_pb"
	app_config "libery_downloads_service/Config"
	"libery_downloads_service/models"
	"libery_downloads_service/workflows"
	"libery_downloads_service/workflows/jobs"
	"net"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DownloadsServer struct {
	srv  *grpc.Server
	port string
	downloads_service_pb.UnimplementedDownloadServiceServer
}

func NewDownloadsServer(service_port string) (*DownloadsServer, error) {
	return &DownloadsServer{
		port: service_port,
	}, nil
}

func (s *DownloadsServer) DownloadImagesBatch(ctx context.Context, in *downloads_service_pb.DownloadImagesBatchRequest) (*downloads_service_pb.DownloadBatchResponse, error) {
	response := new(downloads_service_pb.DownloadBatchResponse)
	var cluster_jwt_token string

	cluster_jwt_token = in.ClusterToken

	recipient_cluster, err := models.GetCategoriesClusterFromToken(cluster_jwt_token, app_config.DOMAIN_SECRET)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In DownloadImagesBatch: error decoding cluster token because '%s'", err.Error()))
		return nil, err
	}

	echo.EchoDebug(fmt.Sprintf("Received download request with: %s", recipient_cluster.FsPath))

	unique_image_urls := workflows.ClearRepeatedDownloadFiles(*in.DownloadUuid, in.ImageUrls)

	if len(unique_image_urls) != len(in.ImageUrls) {
		echo.EchoDebug(fmt.Sprintf("Removed %d repeated images from download %s", len(in.ImageUrls)-len(unique_image_urls), *in.DownloadUuid))
	}

	new_download_uuid, err := jobs.DownloadWorker.DownloadImagesBatch(in.CategoryUuid, unique_image_urls, *in.DownloadUuid, recipient_cluster)
	if err != nil {
		return nil, err
	}

	response.DownloadUuid = new_download_uuid

	return response, nil
}

func (s *DownloadsServer) Connect() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		echo.EchoFatal(err)
	}

	s.srv = grpc.NewServer()
	downloads_service_pb.RegisterDownloadServiceServer(s.srv, s)

	reflection.Register(s.srv)

	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting GRPC server on port: %s", s.port))
	if err := s.srv.Serve(listener); err != nil {
		echo.EchoFatal(err)
	}

	return nil
}

func (s *DownloadsServer) Shutdown() error {
	echo.Echo(echo.GreenFG, "Shutting down GRPC server")
	s.srv.GracefulStop()
	return nil
}
