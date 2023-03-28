package utils

import "fmt"

func MapInterfaceToMapString(mapInterface interface{}) (map[string]string, error) {
	mapString := make(map[string]string)

	switch v := mapInterface.(type) {
	default:
		return nil, fmt.Errorf("unexpected type %T", mapInterface)
	case map[string]interface{}:
		for key, value := range v {
			mapString[key] = fmt.Sprintf("%v", value)
		}
	case map[interface{}]interface{}:
		for key, value := range v {
			mapString[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", value)
		}
	}

	return mapString, nil
}
