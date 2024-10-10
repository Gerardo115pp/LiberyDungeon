package common_workflows

import (
	"context"
	dungeon_models "libery-dungeon-libs/models"
	"libery_categories_service/repository"
)

func GetCategoryIdentityFromUUID(category_uuid string) (category_identity *dungeon_models.CategoryIdentity, lerr *dungeon_models.LabeledError) {
	category, err := repository.CategoriesRepo.GetCategory(context.Background(), category_uuid)
	if err != nil {
		lerr = dungeon_models.NewLabeledError(err, "In GetCategoryIdentityFromUUID", dungeon_models.ErrPlatform_NoSuchCategory)
		return
	}

	category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(context.Background(), category.Cluster)
	if err != nil {
		lerr = dungeon_models.NewLabeledError(err, "In GetCategoryIdentityFromUUID", dungeon_models.ErrPlatform_NoSuchCluster)
		return
	}

	category_identity = dungeon_models.CreateNewCategoryIdentity(&category, &category_cluster)

	return
}
