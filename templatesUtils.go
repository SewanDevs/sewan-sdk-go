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
		template map[string]interface{}) (string, error)
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
			var (
				listTemplateName string = templateList[i].(map[string]interface{})[NameField].(string)
			)
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
		templateRequiredFieldSlice []string = []string{NameField, OsField, RamField,
			CpuField, EnterpriseField, DisksField}
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
		_, ok := template[NicsField]
		if ok && (reflect.TypeOf(template[NicsField]).Kind() != reflect.Slice) {
			templateError = errors.New("Template " + NicsField +
				" is not a list as required but a " +
				reflect.TypeOf(template[NicsField]).Kind().String())
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
				updateSchemaFieldOnResourceCreation(d, key, value)
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
	template map[string]interface{}) (string, error) {
	vm := TemplateCreatedVmOverride{
		RAM:        d.Get(RamField).(int),
		CPU:        d.Get(CpuField).(int),
		Vdc:        d.Get(VdcField).(string),
		Boot:       d.Get(BootField).(string),
		Backup:     d.Get(BackupField).(string),
		Disk_image: d.Get(DiskImageField).(string),
	}
	var (
		schemaer     Schema_Schemaer
		err          error
		listItem     interface{}
		overrideFile strings.Builder
		vmName       strings.Builder
	)
	logger := LoggerCreate("CreateTemplateOverrideConfig_" +
		d.Get(TemplateField).(string) + ".log")
	switch {
	case d.Get(TemplateField) == "":
		err = errors.New("Schema \"Template\" field is empty, " +
			"can not create a template override configuration.")
	default:
		overrideFile.WriteString(d.Get(TemplateField).(string))
		overrideFile.WriteString("_Template_override.tf.json")
		vmName.WriteString(d.Get(NameField).(string))
		_, isSet := d.GetOk(InstanceNumberField)
		if isSet {
			vmName.WriteString(resourceNameCountSeparator)
			vmName.WriteString(resourceDynamicInstanceNumber)
		}
		vm.OS = template[OsField].(string)
		vm.Name = vmName.String()
		if _, err := os.Stat(overrideFile.String()); os.IsNotExist(err) {
			readListValue := []interface{}{}
			for listKey, listValue := range template[DisksField].([]interface{}) {
				listItem, _ = schemaer.ReadElement(listKey,
					listValue,
					logger)
				disk := DiskModifiableFields{
					Name:          listItem.(map[string]interface{})[NameField].(string),
					Size:          listItem.(map[string]interface{})[SizeField].(int),
					Storage_class: listItem.(map[string]interface{})[StorageClassField].(string),
				}
				readListValue = append(readListValue, disk)
			}
			vm.Disks = readListValue
			readListValue = []interface{}{}
			for listKey, listValue := range d.Get(NicsField).([]interface{}) {
				listItem, _ = schemaer.ReadElement(listKey,
					listValue,
					logger)
				nic := NicModifiableFields{
					Vlan:      listItem.(map[string]interface{})[VlanNameField].(string),
					Connected: listItem.(map[string]interface{})[ConnectedField].(bool),
				}
				readListValue = append(readListValue, nic)
			}
			vm.Nics = readListValue
			vmFieldsMap := map[string]interface{}{d.Get(NameField).(string): vm}
			vmMap := map[string]interface{}{"sewan_clouddc_vm": vmFieldsMap}
			resourcesMap := map[string]interface{}{"resource": vmMap}
			vmJson, _ := json.Marshal(resourcesMap)
			err = ioutil.WriteFile(overrideFile.String(),
				vmJson, 0644)
		}
	}
	return overrideFile.String(), err
}

func conformizeNicsSliceOnResourceCreation(d *schema.ResourceData,
	templateParamName string,
	value []interface{}) []interface{} {
	var (
		nicMap          map[string]interface{}
		schemaNicsSlice []interface{}
	)
	for _, nic := range value {
		nicMap = map[string]interface{}{}
		for nicParamName, nicParamValue := range nic.(map[string]interface{}) {
			switch nicParamName {
			case VlanNameField:
				nicMap[nicParamName] = nicParamValue
			case ConnectedField:
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
	return schemaNicsSlice
}

func updateSchemaFieldOnResourceCreation(d *schema.ResourceData, key string, value interface{}) {
	var (
		templateParamName     string      = reflect.ValueOf(key).String()
		interfaceTemplateName interface{} = reflect.ValueOf(value).Interface()
		templateParamValue    string      = reflect.ValueOf(value).String()
	)
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		switch {
		case templateParamName == IdField:
		case templateParamName == OsField:
		case templateParamName == NameField:
		case templateParamName == DatacenterField:
		case d.Get(templateParamName) == "":
			d.Set(templateParamName, templateParamValue)
		default:
		}
	case reflect.Int:
		switch {
		case templateParamName == IdField:
		case d.Get(templateParamName).(int) == 0:
			d.Set(templateParamName, int(interfaceTemplateName.(int)))
		default:
		}
	case reflect.Slice:
		switch {
		case key == DisksField:
		case key == NicsField:
			schemaNicsSlice := conformizeNicsSliceOnResourceCreation(d,
				templateParamName,
				value.([]interface{}))
			d.Set(templateParamName, schemaNicsSlice)
		default:
		}
	default:
	}
}
