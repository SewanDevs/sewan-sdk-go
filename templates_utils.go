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
	FetchTemplateFromList(template_name string,
		templateList []interface{}) (map[string]interface{}, error)
	ValidateTemplate(template map[string]interface{}) error
	UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
		template map[string]interface{}) error
	CreateTemplateOverrideConfig(d *schema.ResourceData,
		template map[string]interface{}) (error, string)
}

type Template_Templater struct{}

type Disk_modifiable_fields struct {
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Storage_class string `json:"storage_class"`
}

type Nic_modifiable_fields struct {
	Vlan      string `json:"vlan"`
	Connected bool   `json:"connected"`
}

type Template_created_VM_override struct {
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

func (templater Template_Templater) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	var (
		template          map[string]interface{} = nil
		template_list_err error                  = nil
	)
	for i := 0; i < len(templateList); i++ {
		switch reflect.TypeOf(templateList[i]).Kind() {
		case reflect.Map:
			var list_template_name string = templateList[i].(map[string]interface{})[NAME_FIELD].(string)
			if list_template_name == template_name {
				template = templateList[i].(map[string]interface{})
				break
			}
		default:
			template_list_err = errors.New("One of the fetch template " +
				"has a wrong format." +
				"\ngot : " + reflect.TypeOf(templateList[i]).Kind().String() +
				"\nwant : " + reflect.Map.String())
			break
		}
	}
	if template == nil && template_list_err == nil {
		template_list_err = errors.New("Template \"" + template_name +
			"\" does not exists, please validate it's name.")
	}
	return template, template_list_err
}

func (templater Template_Templater) ValidateTemplate(template map[string]interface{}) error {
	var (
		template_error                error
		template_required_field_slice []string = []string{NAME_FIELD, OS_FIELD, RAM_FIELD,
			CPU_FIELD, ENTERPRISE_FIELD, DISKS_FIELD}
		missing_fields_list strings.Builder
	)
	for _, elem := range template_required_field_slice {
		if _, ok := template[elem]; !ok {
			missing_fields_list.WriteString("\"")
			missing_fields_list.WriteString(elem)
			missing_fields_list.WriteString("\" ")
		}
	}
	if missing_fields_list.String() != "" {
		template_error = errors.New("Template missing fields : " +
			missing_fields_list.String())
	} else {
		if _, ok := template[NICS_FIELD]; ok {
			if reflect.TypeOf(template[NICS_FIELD]).Kind() != reflect.Slice {
				template_error = errors.New("Template " + NICS_FIELD +
					" is not a list as required but a " +
					reflect.TypeOf(template[NICS_FIELD]).Kind().String())
			}
		}
	}
	return template_error
}

func (templater Template_Templater) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	var template_handle_err error = nil
	if d.Id() == "" {
		for template_param_name, template_param_value := range template {
			if reflect.ValueOf(template_param_name).IsValid() && reflect.ValueOf(template_param_value).IsValid() {
				var (
					s_template_param_name   string      = reflect.ValueOf(template_param_name).String()
					interface_template_name interface{} = reflect.ValueOf(template_param_value).Interface()
					s_template_param_value  string      = reflect.ValueOf(template_param_value).String()
				)
				switch reflect.TypeOf(template_param_value).Kind() {
				case reflect.String:
					switch {
					case s_template_param_name == ID_FIELD:
					case s_template_param_name == OS_FIELD:
					case s_template_param_name == NAME_FIELD:
					case s_template_param_name == DATACENTER_FIELD:
					default:
						if d.Get(s_template_param_name) == "" {
							d.Set(s_template_param_name,
								s_template_param_value)
						}
					}
				case reflect.Int:
					switch {
					case s_template_param_name == ID_FIELD:
					default:
						if d.Get(s_template_param_name).(int) == 0 {
							d.Set(s_template_param_name,
								int(interface_template_name.(int)))
						}
					}
				case reflect.Slice:
					switch {
					case template_param_name == DISKS_FIELD:
					case template_param_name == NICS_FIELD:
						var (
							nic_map           map[string]interface{}
							schema_nics_slice []interface{}
						)
						for _, nic := range template_param_value.([]interface{}) {
							nic_map = map[string]interface{}{}
							for nicParamName, nicParamValue := range nic.(map[string]interface{}) {
								switch nicParamName {
								case VLAN_NAME_FIELD:
									nic_map[nicParamName] = nicParamValue
								case CONNECTED_FIELD:
									nic_map[nicParamName] = nicParamValue
								default:
								}
							}
							schema_nics_slice = append(schema_nics_slice, nic_map)
						}
						for _, nic := range d.Get(s_template_param_name).([]interface{}) {
							schema_nics_slice = append(schema_nics_slice,
								nic.(map[string]interface{}))
						}
						d.Set(s_template_param_name, schema_nics_slice)
					default:
					}
				default:
				}
			}
		}
	} else {
		template_handle_err = errors.New("Template field should not be set on " +
			"an existing resource, please review the configuration field." +
			"\n : The resource schema has not been updated.")
	}
	return template_handle_err
}

func (templater Template_Templater) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	vm := Template_created_VM_override{
		RAM:        d.Get(RAM_FIELD).(int),
		CPU:        d.Get(CPU_FIELD).(int),
		Vdc:        d.Get(VDC_FIELD).(string),
		Boot:       d.Get(BOOT_FIELD).(string),
		Backup:     d.Get(BACKUP_FIELD).(string),
		Disk_image: d.Get(DISK_IMAGE_FIELD).(string),
	}
	var (
		schemaer                Schema_Schemaer
		write_override_file_err error = nil
		read_list_value         []interface{}
		list_item               interface{}
		override_file           strings.Builder
	)
	switch {
	case d.Get(TEMPLATE_FIELD) == "":
		write_override_file_err = errors.New("Schema \"Template\" field is empty, " +
			"can not create a template override configuration.")
	default:
		override_file.WriteString(d.Get(TEMPLATE_FIELD).(string))
		override_file.WriteString("_template_override.tf.json")
		vm.OS = template[OS_FIELD].(string)
		logger := LoggerCreate("CreateTemplateOverrideConfig_" +
			d.Get(TEMPLATE_FIELD).(string) + "_.log")
		res1, err := os.Stat(override_file.String())
		logger.Println("os.Stat(override_file.String()),res1 = ", res1)
		logger.Println("os.Stat(override_file.String()),err = ", err)
		logger.Println("os.IsNotExist(err) = ", os.IsNotExist(err))
		if _, err := os.Stat(override_file.String()); os.IsNotExist(err) {
			for list_key, list_value := range template[DISKS_FIELD].([]interface{}) {
				list_item, _ = schemaer.Read_element(list_key,
					list_value,
					logger)
				disk := Disk_modifiable_fields{
					Name:          list_item.(map[string]interface{})["name"].(string),
					Size:          list_item.(map[string]interface{})["size"].(int),
					Storage_class: list_item.(map[string]interface{})["storage_class"].(string),
				}
				read_list_value = append(read_list_value, disk)
			}
			vm.Disks = read_list_value
			read_list_value = []interface{}{}
			for list_key, list_value := range d.Get(NICS_FIELD).([]interface{}) {
				list_item, _ = schemaer.Read_element(list_key,
					list_value,
					logger)
				nic := Nic_modifiable_fields{
					Vlan:      list_item.(map[string]interface{})["vlan"].(string),
					Connected: list_item.(map[string]interface{})["connected"].(bool),
				}
				read_list_value = append(read_list_value, nic)
			}
			vm.Nics = read_list_value
			vm_fields_map := map[string]interface{}{d.Get(NAME_FIELD).(string): vm}
			vm_map := map[string]interface{}{"sewan_clouddc_vm": vm_fields_map}
			resources_map := map[string]interface{}{"resource": vm_map}
			vm_json, _ := json.Marshal(resources_map)
			write_override_file_err = ioutil.WriteFile(override_file.String(),
				vm_json, 0644)
		}
	}
	return write_override_file_err, override_file.String()
}
