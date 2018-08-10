package sewan_go_sdk

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"reflect"
	"strconv"
	"strings"
)

type SchemaTooler struct {
	SchemaTools Schemaer
}
type Schemaer interface {
	DeleteTerraformResource(d *schema.ResourceData)
	UpdateLocalResourceState(resourceState map[string]interface{},
		d *schema.ResourceData, schemaTools *SchemaTooler) error
	UpdateVdcResourcesNames(d *schema.ResourceData) error
	ReadElement(key interface{}, value interface{}) (interface{}, error)
}
type SchemaSchemaer struct{}

func (schemaer SchemaSchemaer) DeleteTerraformResource(d *schema.ResourceData) {
	d.SetId("")
}

// Update of resource state in .tfstate file through schema update
func (schemaer SchemaSchemaer) UpdateLocalResourceState(resourceState map[string]interface{},
	d *schema.ResourceData, schemaTools *SchemaTooler) error {
	var (
		updateError error = nil
		readValue   interface{}
	)
	for key, value := range resourceState {
		readValue, updateError = schemaTools.SchemaTools.ReadElement(key, value)
		if key == IdField {
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

// Trim meaningless part of vdc resource name to store a shorter name locally
// exemple of a trim : "<enterprise name>-mono-ram" -> "ram"
func (schemaer SchemaSchemaer) UpdateVdcResourcesNames(d *schema.ResourceData) error {
	var (
		vdcResourcesList       []interface{} = d.Get(VdcResourceField).([]interface{})
		vdcResourcesListUpdate []interface{} = []interface{}{}
		enterpriseName         string        = d.Get(EnterpriseField).(string)
		resourceName           string
	)
	for _, resource := range vdcResourcesList {
		resourceName = resource.(map[string]interface{})[ResourceField].(string)
		resourceName = strings.Replace(resourceName,
			enterpriseName, "", 1)
		resourceName = strings.Replace(resourceName,
			monoField, "", 1)
		resource.(map[string]interface{})[ResourceField] = resourceName
		vdcResourcesListUpdate = append(vdcResourcesListUpdate, resource)
	}
	return d.Set(VdcResourceField, vdcResourcesListUpdate)
}

// Format Element(key,value) value type to a type accepted by terraform :
//
// * value type -> terraform accepted type
//
// * string -> string
//
// * bool -> bool
//
// * float64 -> int (rounded to nearest int)
//
// * int -> int
//
// * map -> map (recursive call of function for each map element)
//
// * slice -> slice (recursive call of function for each slice element)
//
// * other types -> return error
func (schemaer SchemaSchemaer) ReadElement(key interface{},
	value interface{}) (interface{}, error) {
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
				mapValue)
			readMapValue[mapKey] = mapItem
		}
		readValue = readMapValue
	case []interface{}:
		var readListValue []interface{}
		var listItem interface{}
		for listKey, listValue := range valueType {
			listItem,
				readError = schemaer.ReadElement(listKey,
				listValue)
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
