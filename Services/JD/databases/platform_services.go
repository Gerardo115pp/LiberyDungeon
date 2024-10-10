package databases

import (
	"fmt"
	dungeon_models "libery-dungeon-libs/models"
	"libery-dungeon-libs/models/platform_services"
	service_models "libery_JD_service/models"
)

type PlatformServicesDatastore struct {
	onlineServices map[platform_services.PlatformServiceName]*service_models.PlatformService
}

func NewPlatformServicesDatastore() *PlatformServicesDatastore {
	var new_datastore *PlatformServicesDatastore = new(PlatformServicesDatastore)

	new_datastore.onlineServices = make(map[platform_services.PlatformServiceName]*service_models.PlatformService)
	new_datastore.fillOnlineServicesTable()

	return new_datastore
}

func (psd *PlatformServicesDatastore) fillOnlineServicesTable() {
	for _, service := range platform_services.ExistingServices {
		var new_service *service_models.PlatformService = service_models.NewPlatformService(service, "", "")

		if service == platform_services.JD_SERVICE {
			new_service.IsOnline = true
		}

		psd.onlineServices[service] = new_service
	}
}

func (psd PlatformServicesDatastore) IsServiceOnline(service_name string) bool {
	var service_platform_name platform_services.PlatformServiceName = platform_services.PlatformServiceName(service_name)

	service, ok := psd.onlineServices[service_platform_name]
	if !ok {
		return false
	}

	return service.IsOnline
}

// Sets the service online state, if there is an error, it will be because the service does not exist.
func (psd *PlatformServicesDatastore) SetServiceOnline(service_name string, online bool) (lerr *dungeon_models.LabeledError) {
	if !platform_services.ServiceNameValid(service_name) {
		lerr = dungeon_models.NewLabeledError(fmt.Errorf("service '%s' is not valid", service_name), "In NewPLatformServicesDatastore.SetServiceOnline while checking if service is valid", service_models.ErrService_NotSuchService)
		return
	}

	var service_platform_name platform_services.PlatformServiceName = platform_services.PlatformServiceName(service_name)

	_, ok := psd.onlineServices[service_platform_name]
	if !ok {
		lerr = dungeon_models.NewLabeledError(fmt.Errorf("service '%s' does not exist", service_name), "In NewPLatformServicesDatastore.SetServiceOnline while checking if service exists", service_models.ErrService_NotSuchService)
		return
	}

	psd.onlineServices[service_platform_name].IsOnline = online

	return nil
}

func (psd *PlatformServicesDatastore) SetServiceMetadata(new_service service_models.PlatformService) (lerr *dungeon_models.LabeledError) {
	service, exists := psd.onlineServices[new_service.ServiceName]
	if !exists {
		lerr = dungeon_models.NewLabeledError(fmt.Errorf("service '%s' does not exist", new_service.ServiceName), "In NewPLatformServicesDatastore.SetServiceMetadata while checking if service exists", service_models.ErrService_NotSuchService)
		return
	}

	service.ServiceRoute = new_service.ServiceRoute
	service.ServicePort = new_service.ServicePort
	service.IsOnline = new_service.IsOnline

	return nil
}
