package middleware

import (
	app_config "libery_categories_service/Config"
	"libery_categories_service/repository"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func CheckMainCategoryProxyID(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		category_id := request.URL.Query().Get("category_id")

		echo.EchoDebug("main_category_proxy_id: " + app_config.MAIN_CATEGORY_PROXY_ID)
		if category_id == app_config.MAIN_CATEGORY_PROXY_ID {
			cluster_id := request.URL.Query().Get("cluster_id")

			if cluster_id == "" {
				echo.Echo(echo.RedBG, "Missing cluster_id parameter, impossible to get main category without cluster_id")
				response.WriteHeader(400)
				return
			}

			category_cluster, err := repository.CategoriesClustersRepo.GetClusterByID(request.Context(), cluster_id)
			if err != nil {
				echo.Echo(echo.RedBG, "Error getting category cluster: "+err.Error())
				response.WriteHeader(404)
				return
			}

			category_id = category_cluster.RootCategory

			params := request.URL.Query()
			params.Set("category_id", category_id)
			request.URL.RawQuery = params.Encode()
		}

		if category_id != "" {
			echo.Echo(echo.GreenFG, "category_id: "+request.URL.Query().Get("category_id"))
		}

		next(response, request)
	}
}
