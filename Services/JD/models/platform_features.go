package models

import "libery-dungeon-libs/models/platform_services"

type PlatformService struct {
	ServiceName  platform_services.PlatformServiceName `json:"service_name"`
	ServiceRoute string                                `json:"service_route"`
	ServicePort  string                                `json:"service_port"`
	IsOnline     bool                                  `json:"is_online"`
}

func NewPlatformService(service_name platform_services.PlatformServiceName, service_route, service_port string) *PlatformService {
	return &PlatformService{
		ServiceName:  service_name,
		ServiceRoute: service_route,
		ServicePort:  service_port,
		IsOnline:     false,
	}
}
