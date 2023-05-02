package main

import (
	"encoding/json"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs"
	"os"
	"reflect"
)

func processValue(value reflect.Value) (interface{}, bool) {
	flag := false
	for value.Kind() == reflect.Ptr {
		if !value.IsValid() || value.IsNil() {
			flag = true
			break
		}
		value = value.Elem()
	}
	if flag || valueIsValid(value) {
		return nil, false
	}

	switch value.Kind() {
	case reflect.Struct:
		return walkStruct(value), true
	case reflect.Map:
		return walkMap(value), true
	case reflect.Slice:
		return walkSlice(value), true
	case reflect.Interface:
		return processValue(value.Elem()) // Handle the interface by processing its underlying value
	case reflect.Func:
		return nil, false
	default:
		return value.Interface(), true
	}
}

func walkStruct(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		if field.PkgPath != "" {
			continue
		}

		if processedValue, ok := processValue(value); ok {
			result[field.Name] = processedValue
		}
	}

	return result
}

func walkSlice(v reflect.Value) []interface{} {
	result := make([]interface{}, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		value := v.Index(i)

		if processedValue, ok := processValue(value); ok {
			result = append(result, processedValue)
		}
	}

	return result
}

func walkMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)

		if processedValue, ok := processValue(value); ok {
			result[key.String()] = processedValue
		}
	}

	return result
}

func valueIsValid(value reflect.Value) bool {
	return !value.IsValid() ||
		value.Type().String() == "schema.SchemaDiffSuppressFunc" ||
		value.Type().String() == "schema.SchemaStateFunc" ||
		value.Type().String() == "schema.SchemaValidateFunc" ||
		value.Type().String() == "schema.DeleteContextFunc" ||
		value.Type().String() == "schema.ReadContextFunc" ||
		value.Type().String() == "schema.UpdateContextFunc" ||
		value.Type().String() == "schema.CustomizeDiffFunc" ||
		value.Type().String() == "schema.StateFunc" ||
		value.Type().String() == "schema.StateContextFunc" ||
		value.Type().String() == "schema.SchemaSetFunc" ||
		value.Type().String() == "schema.SchemaValidateDiagFunc" ||
		value.Type().String() == "schema.SchemaDefaultFunc" ||
		value.Type().String() == "schema.StateMigrateFunc" ||
		value.Type().String() == "schema.CreateFunc" ||
		value.Type().String() == "schema.ReadFunc" ||
		value.Type().String() == "schema.UpdateFunc" ||
		value.Type().String() == "schema.DeleteFunc" ||
		value.Type().String() == "schema.ExistsFunc" ||
		value.Type().String() == "schema.CreateContextFunc" ||
		value.Type().String() == "schema.StateMigrateFunc" ||
		value.Type().String() == "schema.ConfigureFunc" ||
		value.Type().String() == "schema.ConfigureContextFunc"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Instantiate the plugin
	provider := vkcs.Provider()

	// Extract provider information
	providerValue := reflect.ValueOf(provider).Elem()
	providerInfo := walkStruct(providerValue)

	//fmt.Println("providerInfo: ", providerInfo)

	// Convert provider information to JSON
	m, err := json.MarshalIndent(providerInfo, "", "  ")
	check(err)

	// create file if not exists
	_, err = os.Stat("./vk.json")
	if os.IsNotExist(err) {
		_, err = os.Create("./vk.json")
		check(err)
	}
	err = os.WriteFile("./vk.json", m, 0644)
	check(err)
	// save to file

}
