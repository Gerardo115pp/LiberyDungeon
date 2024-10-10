package helpers

import "strings"

func RemoveRoutePrefix(route string, prefix string) string {
	var clear_route string = ""
	segments := strings.Split(route, prefix)
	if len(segments) == 2 {
		clear_route = segments[1]
	}

	return clear_route
}
