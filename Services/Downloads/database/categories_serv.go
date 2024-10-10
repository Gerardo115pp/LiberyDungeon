package database

import (
	"context"
	"fmt"
	"libery-dungeon-libs/categories_service_pb"
	app_config "libery_downloads_service/Config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CategoriesServiceConn struct {
	CategoriesAddress string
	CertFile          string
}

func NewCategoriesServiceConn() *CategoriesServiceConn {
	return &CategoriesServiceConn{
		CategoriesAddress: app_config.GRPC_SERVER,
		CertFile:          app_config.SSL_CA_PATH,
	}
}

func (csc *CategoriesServiceConn) getCredentials() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(csc.CertFile, "")
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func (csc *CategoriesServiceConn) CreateCategory(ctx context.Context, name string, parent_uuid string) (category_uuid string, err error) {
	creds, err := csc.getCredentials()
	if err != nil {
		return "", err
	}

	conn, err := grpc.Dial(csc.CategoriesAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("Error connecting to Categories Service: %s", err.Error())
	}

	defer conn.Close()

	client := categories_service_pb.NewCategoriesServiceClient(conn)

	request_data := new(categories_service_pb.CreateCategoryRequest)

	request_data.Name = name
	request_data.Parent = parent_uuid

	request, err := client.CreateCategory(ctx, request_data)
	if err != nil {
		return "", fmt.Errorf("Error creating category: %s", err.Error())
	}

	return request.Uuid, nil
}
