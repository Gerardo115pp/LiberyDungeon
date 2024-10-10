package middleware

import (
	"fmt"
	"libery_medias_service/helpers"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

var need_parsing map[string][]string = map[string][]string{
	"/medias-fs":     {"GET"},
	"/thumbnails-fs": {"GET"},
}

var exceptions map[string]struct{} = map[string]struct{}{
	"/thumbnails-fs/libery-trashcan/": struct{}{},
}

func ParseResourcePath(next func(response http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		direct_parent_resource, _ := filepath.Split(request.URL.Path)
		if IsException(direct_parent_resource) {
			echo.Echo(echo.BlueBG, "exception for media resource: "+request.URL.Path)
			next(response, request)
			return
		}

		var needs_parsing bool = false
		var route_prefix string = ""

		for path, methods := range need_parsing {
			route_prefix = path
			matches := strings.HasPrefix(request.URL.Path, path)
			echo.EchoDebug(fmt.Sprintf("path: %s\nrequest.URL.Path: %s\nmatches: %t", path, request.URL.Path, matches))
			if matches {

				for _, method := range methods {
					if request.Method == method {
						needs_parsing = true
						break
					}
				}
				break
			}
		}

		if !needs_parsing {
			next(response, request)
			return
		}

		media_resource_path := request.URL.Path

		echo.EchoDebug("route_prefix: " + route_prefix)
		echo.EchoDebug("media_resource_path: " + media_resource_path)

		media_resource_path = helpers.RemoveRoutePrefix(media_resource_path, strings.Trim(route_prefix, "/"))

		request.URL.Path = media_resource_path

		echo.Echo(echo.BlueBG, "request for media resource: "+media_resource_path)

		next(response, request)
	}
}

func IsException(resource string) bool {
	_, ok := exceptions[resource]
	return ok
}
