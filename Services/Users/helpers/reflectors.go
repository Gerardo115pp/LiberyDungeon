package helpers

import "reflect"

func StructToMap(obj interface{}) map[string]interface{} {
	interface_map := make(map[string]interface{})
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	value_type := value.Type()

	for h := 0; h < value.NumField(); h++ {
		field := value_type.Field(h)
		field_value := value.Field(h).Interface()
		tag_name := field.Tag.Get("json")

		if tag_name == "" {
			tag_name = field.Name
		}

		interface_map[tag_name] = field_value
	}

	return interface_map
}
