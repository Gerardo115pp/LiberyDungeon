package helpers

import "strings"

// take a string of the form "uuid1,uuid2,uuid3" and returns a lookup
// map where the keys are the UUIDs. empty strings will be excluded.
func SplitCommaSeparatedUUIDs(input string) map[string]struct{} {
	var uuids []string = strings.Split(input, ",")
	var lookup_map map[string]struct{} = make(map[string]struct{})

	for _, uuid := range uuids {
		uuid = strings.TrimSpace(uuid)
		if uuid != "" {
			lookup_map[uuid] = struct{}{}
		}
	}

	return lookup_map
}
