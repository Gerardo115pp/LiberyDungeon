package server

import (
	"context"
	"fmt"
	"libery-dungeon-libs/categories_service_pb"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	"net"

	"github.com/Gerardo115pp/patriots_lib/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CategoriesServer struct {
	srv  *grpc.Server
	port string
	categories_service_pb.UnimplementedCategoriesServiceServer
}

func NewCategoriesServer(service_port string) (*CategoriesServer, error) {
	return &CategoriesServer{
		port: service_port,
	}, nil
}

func (s *CategoriesServer) CreateCategory(ctx context.Context, request *categories_service_pb.CreateCategoryRequest) (*categories_service_pb.CreateCategoryResponse, error) {
	response := new(categories_service_pb.CreateCategoryResponse)

	new_category, err := workflows.CreateNewCategory(ctx, request.Name, request.Parent, request.Cluster)
	if err != nil {
		return nil, err
	}

	response.Uuid = new_category.Uuid

	return response, nil
}

func (s *CategoriesServer) GetCategoriesCluster(ctx context.Context, request *categories_service_pb.GetCategoriesClusterRequest) (*categories_service_pb.GetCategoriesClusterResponse, error) {
	var response *categories_service_pb.GetCategoriesClusterResponse = new(categories_service_pb.GetCategoriesClusterResponse)

	cluster, err := repository.CategoriesClustersRepo.GetClusterByID(ctx, request.Uuid)
	if err != nil {
		return nil, err
	}

	var proto_cluster *categories_service_pb.CategoriesCluster = new(categories_service_pb.CategoriesCluster)

	proto_cluster.Uuid = cluster.Uuid
	proto_cluster.Name = cluster.Name
	proto_cluster.FsPath = cluster.FsPath
	proto_cluster.FilterCategory = cluster.FilterCategory
	proto_cluster.RootCategory = cluster.RootCategory

	response.Cluster = proto_cluster

	return response, nil
}

func (s *CategoriesServer) GetCategory(ctx context.Context, request *categories_service_pb.GetCategoryRequest) (*categories_service_pb.GetCategoryResponse, error) {
	response := new(categories_service_pb.GetCategoryResponse)

	category, err := repository.CategoriesRepo.GetCategory(ctx, request.Uuid)
	if err != nil {
		return nil, err
	}

	proto_category := new(categories_service_pb.Category)

	proto_category.Uuid = category.Uuid
	proto_category.Name = category.Name
	proto_category.Parent = category.Parent
	proto_category.Fullpath = category.Fullpath

	response.Category = proto_category

	return response, nil
}

func (s *CategoriesServer) Connect() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		echo.EchoFatal(err)
	}

	s.srv = grpc.NewServer()
	categories_service_pb.RegisterCategoriesServiceServer(s.srv, s)

	reflection.Register(s.srv)

	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting GRPC server on port: %s", s.port))
	if err := s.srv.Serve(listener); err != nil {
		echo.EchoFatal(err)
	}

	return nil
}

func (s *CategoriesServer) Shutdown() error {
	if s.srv == nil {
		return nil
	}

	s.srv.Stop()

	return nil
}
