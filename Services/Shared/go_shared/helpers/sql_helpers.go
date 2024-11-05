package helpers

import "strings"

func GetPreparedListPlaceholders(count int) string {
	var placeholders string

	placeholders_builder := strings.Builder{}
	placeholders_builder.WriteString("?")

	for h := 1; h < count; h++ {
		placeholders_builder.WriteString(", ?")
	}

	placeholders = placeholders_builder.String()

	return placeholders
}
