package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/yandex-cloud/terraform-provider-yandex/yandex"
)

func walkStruct(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)
		fieldName := field.Name

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		switch value.Kind() {
		case reflect.Struct:
			result[fieldName] = walkStruct(value)
		case reflect.Map:
			result[fieldName] = walkMap(value)
		case reflect.Func:
			// Skip functions
		default:
			result[fieldName] = value.Interface()
		}
	}

	return result
}

func walkMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)

		switch value.Kind() {
		case reflect.Struct:
			result[key.String()] = walkStruct(value)
		case reflect.Map:
			result[key.String()] = walkMap(value)
		case reflect.Func:
			// Skip functions
		default:
			result[key.String()] = value.Interface()
		}
	}

	return result
}

func main() {
	// Instantiate the plugin
	provider := yandex.Provider()

	// Extract provider information
	providerValue := reflect.ValueOf(provider).Elem()
	providerInfo := walkStruct(providerValue)

	// Convert provider information to JSON
	m, err := json.MarshalIndent(providerInfo, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(m))
}
