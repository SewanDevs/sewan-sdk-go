package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type TemplatesTooler struct {
	TemplatesTools Templater
}
type Templater interface {
	FetchTemplateFromList(templateName string,
		templateList []interface{}) (map[string]interface{}, error)
	ValidateTemplate(template map[string]interface{}) error
	UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
		template map[string]interface{}) error
	CreateTemplateOverrideConfig(d *schema.ResourceData,
		template map[string]interface{}) (error, string)
}

type Template_Templater struct{}

type DiskModifiableFields struct {
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Storage_class string `json:"storage_class"`
}

type NicModifiableFields struct {
	Vlan      string `json:"vlan"`
	Connected bool   `json:"connected"`
}

type TemplateCreatedVmOverride struct {
	Name       string        `json:"name"`
	OS         string        `json:"os"`
	RAM        int           `json:"ram"`
	CPU        int           `json:"cpu"`
	Disks      []interface{} `json:"disks,omitempty"`
	Nics       []interface{} `json:"nics,omitempty"`
	Vdc        string        `json:"vdc"`
	Boot       string        `json:"boot"`
	Backup     string        `json:"backup"`
	Disk_image string        `json:"disk_image"`
}

func (templater Template_Templater) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {

	var (
		template          map[string]interface{} = nil
		templateListError error                  = nil
	)
	for i := 0; i < len(templateList); i++ {
		switch reflect.TypeOf(templateList[i]).Kind() {
		case reflect.Map:
			var listTemplateName string = templateList[i].(map[string]interface{})[NAME_FIELD].(string)
			if listTemplateName == templateName {
				template = templateList[i].(map[string]interface{})
				break
			}
		default:
			templateListError = errors.New("One of the fetch template " +
				"has a wrong format." +
				"\ngot : " + reflect.TypeOf(templateList[i]).Kind().String() +
				"\nwant : " + reflect.Map.String())
			break
		}
	}
	if template == nil && templateListError == nil {
		templateListError = errors.New("Template \"" + templateName +
			"\" does not exists, please validate it's name.")
	}
	return template, templateListError
}

func (templater Template_Templater) ValidateTemplate(template map[string]interface{}) error {
	var (
		templateError              error
		templateRequiredFieldSlice []string = []string{NAME_FIELD, OS_FIELD, RAM_FIELD,
			CPU_FIELD, ENTERPRISE_FIELD, DISKS_FIELD}
		missingFieldsList strings.Builder
	)
	for _, elem := range templateRequiredFieldSlice {
		if _, ok := template[elem]; !ok {
			missingFieldsList.WriteString("\"")
			missingFieldsList.WriteString(elem)
			missingFieldsList.WriteString("\" ")
		}
	}
	if missingFieldsList.String() != "" {
		templateError = errors.New("Template missing fields : " +
			missingFieldsList.String())
	} else {
		if _, ok := template[NICS_FIELD]; ok {
			if reflect.TypeOf(template[NICS_FIELD]).Kind() != reflect.Slice {
				templateError = errors.New("Template " + NICS_FIELD +
					" is not a list as required but a " +
					reflect.TypeOf(template[NICS_FIELD]).Kind().String())
			}
		}
	}
	return templateError
}

func (templater Template_Templater) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	var templateHandleError error = nil
	if d.Id() == "" {
		for key, value := range template {
			if reflect.ValueOf(key).IsValid() && reflect.ValueOf(value).IsValid() {
				var (
					templateParamName     string      = reflect.ValueOf(key).String()
					interfaceTemplateName interface{} = reflect.ValueOf(value).Interface()
					templateParamValue    string      = reflect.ValueOf(value).String()
				)
				switch reflect.TypeOf(value).Kind() {
				case reflect.String:
					switch {
					case templateParamName == ID_FIELD:
					case templateParamName == OS_FIELD:
					case templateParamName == NAME_FIELD:
					case templateParamName == DATACENTER_FIELD:
					default:
						if d.Get(templateParamName) == "" {
							d.Set(templateParamName,
								templateParamValue)
						}
					}
				case reflect.Int:
					switch {
					case templateParamName == ID_FIELD:
					default:
						if d.Get(templateParamName).(int) == 0 {
							d.Set(templateParamName,
								int(interfaceTemplateName.(int)))
						}
					}
				case reflect.Slice:
					switch {
					case key == DISKS_FIELD:
					case key == NICS_FIELD:
						var (
							nicMap          map[string]interface{}
							schemaNicsSlice []interface{}
						)
						for _, nic := range value.([]interface{}) {
							nicMap = map[string]interface{}{}
							for nicParamName, nicParamValue := range nic.(map[string]interface{}) {
								switch nicParamName {
								case VLAN_NAME_FIELD:
									nicMap[nicParamName] = nicParamValue
								case CONNECTED_FIELD:
									nicMap[nicParamName] = nicParamValue
								default:
								}
							}
							schemaNicsSlice = append(schemaNicsSlice, nicMap)
						}
						for _, nic := range d.Get(templateParamName).([]interface{}) {
							schemaNicsSlice = append(schemaNicsSlice,
								nic.(map[string]interface{}))
						}
						d.Set(templateParamName, schemaNicsSlice)
					default:
					}
				default:
				}
			}
		}
	} else {
		templateHandleError = errors.New("Template field should not be set on " +
			"an existing resource, please review the configuration field." +
			"\n : The resource schema has not been updated.")
	}
	return templateHandleError
}

func (templater Template_Templater) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	vm := TemplateCreatedVmOverride{
		RAM:        d.Get(RAM_FIELD).(int),
		CPU:        d.Get(CPU_FIELD).(int),
		Vdc:        d.Get(VDC_FIELD).(string),
		Boot:       d.Get(BOOT_FIELD).(string),
		Backup:     d.Get(BACKUP_FIELD).(string),
		Disk_image: d.Get(DISK_IMAGE_FIELD).(string),
	}
	var (
		schemaer               Schema_Schemaer
		writeOverrideFileError error = nil
		readListValue          []interface{}
		listItem               interface{}
		overrideFile           strings.Builder
		vmName                 strings.Builder
		isSet                  bool
	)
	switch {
	case d.Get(TEMPLATE_FIELD) == "":
		writeOverrideFileError = errors.New("Schema \"Template\" field is empty, " +
			"can not create a template override configuration.")
	default:
		logger := LoggerCreate("CreateTemplateOverrideConfig_" +
			d.Get(TEMPLATE_FIELD).(string) + "_.log")
		overrideFile.WriteString(d.Get(TEMPLATE_FIELD).(string))
		overrideFile.WriteString("Template_override.tf.json")
		vmName.WriteString(d.Get(NAME_FIELD).(string))
		_, isSet = d.GetOk(INSTANCE_NUMBER_FIELD)
		if isSet {
			vmName.WriteString(RESOURCE_NAME_COUNT_SEPARATOR)
			vmName.WriteString(RESOURCE_INSTANCE_NUMBER)
		}
		vm.OS = template[OS_FIELD].(string)
		vm.Name = vmName.String()
		if _, err := os.Stat(overrideFile.String()); os.IsNotExist(err) {
			for listKey, listValue := range template[DISKS_FIELD].([]interface{}) {
				listItem, _ = schemaer.ReadElement(listKey,
					listValue,
					logger)
				disk := DiskModifiableFields{
					Name:          listItem.(map[string]interface{})[NAME_FIELD].(string),
					Size:          listItem.(map[string]interface{})[SIZE_FIELD].(int),
					Storage_class: listItem.(map[string]interface{})[STORAGE_CLASS_FIELD].(string),
				}
				readListValue = append(readListValue, disk)
			}
			vm.Disks = readListValue
			readListValue = []interface{}{}
			for listKey, listValue := range d.Get(NICS_FIELD).([]interface{}) {
				listItem, _ = schemaer.ReadElement(listKey,
					listValue,
					logger)
				nic := NicModifiableFields{
					Vlan:      listItem.(map[string]interface{})[VLAN_NAME_FIELD].(string),
					Connected: listItem.(map[string]interface{})[CONNECTED_FIELD].(bool),
				}
				readListValue = append(readListValue, nic)
			}
			vm.Nics = readListValue
			vm_fields_map := map[string]interface{}{d.Get(NAME_FIELD).(string): vm}
			vm_map := map[string]interface{}{"sewan_clouddc_vm": vm_fields_map}
			resources_map := map[string]interface{}{"resource": vm_map}
			vmJson, _ := json.Marshal(resources_map)
			writeOverrideFileError = ioutil.WriteFile(overrideFile.String(),
				vmJson, 0644)
		}
	}
	return writeOverrideFileError, overrideFile.String()
}
