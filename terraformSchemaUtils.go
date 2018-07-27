package sewan_go_sdk

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"reflect"
	"strconv"
)

type SchemaTooler struct {
	SchemaTools Schemaer
}
type Schemaer interface {
	DeleteTerraformResource(d *schema.ResourceData)
	UpdateLocalResourceState(resource_state map[string]interface{},
		d *schema.ResourceData, schemaTools *SchemaTooler) error
	ReadElement(key interface{}, value interface{},
		logger *log.Logger) (interface{}, error)
}
type Schema_Schemaer struct{}

func (schemaer Schema_Schemaer) DeleteTerraformResource(d *schema.ResourceData) {
	d.SetId("")
}

func (schemaer Schema_Schemaer) UpdateLocalResourceState(resource_state map[string]interface{},
	d *schema.ResourceData, schemaTools *SchemaTooler) error {

	var (
		updateError error = nil
		readValue   interface{}
	)
	logger := LoggerCreate("update_local_resource_state_" +
		d.Get(NAME_FIELD).(string) + ".log")
	for key, value := range resource_state {
		readValue,
			updateError = schemaTools.SchemaTools.ReadElement(key,
			value,
			logger)
		logger.Println("Set \"", key, "\" to \"", readValue, "\"")
		if key == ID_FIELD {
			var s_id string = ""
			switch {
			case reflect.TypeOf(value).Kind() == reflect.Float64:
				s_id = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			case reflect.TypeOf(value).Kind() == reflect.Int:
				s_id = strconv.Itoa(value.(int))
			case reflect.TypeOf(value).Kind() == reflect.String:
				if value == nil {
					s_id = ""
				} else {
					s_id = value.(string)
				}
			default:
				updateError = errors.New("Format of " + key + "(" +
					reflect.TypeOf(value).Kind().String() + ") not handled.")
			}
			d.SetId(s_id)
		} else {
			updateError = d.Set(key, readValue)
		}
		readValue = nil
	}
	return updateError
}

func (schemaer Schema_Schemaer) ReadElement(key interface{}, value interface{},
	logger *log.Logger) (interface{}, error) {

	var (
		readError error = nil
		readValue interface{}
	)
	switch valueType := value.(type) {
	case string:
		readValue = value.(string)
	case bool:
		readValue = value.(bool)
	case float64:
		readValue = int(value.(float64))
	case int:
		readValue = value.(int)
	case map[string]interface{}:
		var readMapValue map[string]interface{}
		readMapValue = make(map[string]interface{})
		var mapItem interface{}
		for mapKey, mapValue := range valueType {
			mapItem,
				readError = schemaer.ReadElement(mapKey,
				mapValue,
				logger)
			readMapValue[mapKey] = mapItem
		}
		readValue = readMapValue
	case []interface{}:
		var readListValue []interface{}
		var listItem interface{}
		for listKey, listValue := range valueType {
			listItem,
				readError = schemaer.ReadElement(listKey,
				listValue,
				logger)
			readListValue = append(readListValue, listItem)
		}
		readValue = readListValue
	default:
		if value == nil {
			readValue = nil
		} else {
			readError = errors.New("Format " +
				reflect.TypeOf(valueType).Kind().String() + " not handled.")
		}
	}
	return readValue, readError
}
