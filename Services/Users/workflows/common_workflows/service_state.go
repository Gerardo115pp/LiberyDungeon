package common_workflows

import (
	"context"
	service_models "libery_users_service/models"
	"libery_users_service/repository"
)

func InitialSetupDone(ctx context.Context) (bool, error) {
	var all_users []*service_models.User
	var err error

	all_users, err = repository.UsersRepo.GetAllUsersCTX(ctx)
	if err != nil {
		return false, err
	}

	return len(all_users) != 0, nil
}
