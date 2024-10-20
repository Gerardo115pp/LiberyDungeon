package helpers

func GetPreparedListPlaceholders(count int) string {
	var placeholders string = "?"

	for h := 1; h < count; h++ {
		placeholders += ", ?"
	}

	return placeholders
}
