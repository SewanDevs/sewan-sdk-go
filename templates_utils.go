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
	UpdateSchemaFromTemplate(d *schema.ResourceData,
		template map[string]interface{},
		templatesTooler *TemplatesTooler,
		schemaTools *SchemaTooler) error
	CreateTemplateOverrideConfig(d *schema.ResourceData, template map[string]interface{}) error
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
	OS        string           `json:"os"`
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
		template            map[string]interface{} = nil
		template_list_valid error                  = nil
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
			template_list_valid = errors.New("Wrong template list format.\n" +
				"got :" + reflect.TypeOf(templateList[i]).Kind().String() +
				"want :" + reflect.Map.String())
		}
	}
	if template == nil {
		template_list_valid = errors.New("Template \"" + template_name +
			"\" does not exists, please validate it's name.")
	}
	return template, template_list_valid
}

func (templater Template_Templater) UpdateSchemaFromTemplate(d *schema.ResourceData,
	template map[string]interface{},
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler) error {
	var template_handle_err error = nil
	for template_param_name, template_param_value := range template {
		if reflect.ValueOf(template_param_name).IsValid() && reflect.ValueOf(template_param_value).IsValid() {
			logger.Println("--")
			var (
				s_template_param_name   string      = reflect.ValueOf(template_param_name).String()
				interface_template_name interface{} = reflect.ValueOf(template_param_value).Interface()
				s_template_param_value  string      = reflect.ValueOf(template_param_value).String()
			)
			switch reflect.TypeOf(template_param_value).Kind() {
			case reflect.String:
				logger.Println("case String : ", template_param_name)
				if d.Id() == "" {
					switch {
					case s_template_param_name == OS_FIELD && d.Id()=="":
						logger.Println("Case name")
					case s_template_param_name == NAME_FIELD:
						logger.Println("Case name")
					default:
						if d.Get(s_template_param_name) == "" {
							logger.Println("Case : ", s_template_param_name)
							d.Set(s_template_param_name,
								s_template_param_value)
						}
					}
				} else {
					switch {
					case s_template_param_name == NAME_FIELD:
						logger.Println("Case template name")
						data := &Dynamic_field_struct{}
						dynamicfield_read_err := json.Unmarshal([]byte(d.Get(DYNAMIC_FIELD).(string)), data)
						if dynamicfield_read_err == nil {
							if s_template_param_value != data.Creation_template {
								if data.Creation_template == "" {
									template_handle_err = errors.New("This resource has not been " +
										"created with a template. Please remove template field from" +
										"the configuration file.")
								} else {
									template_handle_err = errors.New("This resource has been " +
										"created with \"" + data.Creation_template +
										"\" template. This value can not be changed, please set it back.")
								}
							}
						} else {
							template_handle_err = errors.New(d.Get(NAME_FIELD).(string) +
								"'s resource dynamic field is not a valid json, please make sure" +
								" this resource is modified only by a terraform session, \n" +
								"json error :" + dynamicfield_read_err.Error())
						}
					default:
						if d.Get(s_template_param_name) == "" {
							logger.Println("3a : ", s_template_param_name)
							d.Set(s_template_param_name,
								s_template_param_value)
						}
					}
				}
			case reflect.Float64:
				logger.Println("case float 64 : ", template_param_name, " = ",
					d.Get(s_template_param_name))
				switch {
				case s_template_param_name == ID_FIELD:
					logger.Println("2, d.Id() = ", d.Id())
				default:
					if d.Get(s_template_param_name).(int) == 0 {
						logger.Println("3, val to set = ",
							int(interface_template_name.(float64)))
						d.Set(s_template_param_name,
							int(interface_template_name.(float64)))
					}
				}
			case reflect.Int:
				logger.Println("case Int : ", template_param_name, " = ",
					d.Get(s_template_param_name))
				switch {
				case s_template_param_name == ID_FIELD:
					logger.Println("2")
				default:
					if d.Get(s_template_param_name).(int) == 0 {
						logger.Println("3, val to set = ",
							int(interface_template_name.(int)))
						d.Set(s_template_param_name,
							int(interface_template_name.(int)))
					}
				}
			case reflect.Slice:
				logger.Println("case Slice : ", template_param_name, " = ",
					d.Get(s_template_param_name))
				switch {
				case template_param_name == DISKS_FIELD && d.Id()=="":
				default:
					d.Set(s_template_param_name, template_param_value.([]interface{}))
				}
				if template_handle_err != nil {
					logger.Println(template_param_name, "=",
						d.Get(s_template_param_name),
						"error :", template_handle_err)
					break
				}
			default:
			}
		}
	}
	return template_handle_err
}

func (templater Template_Templater) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) error {
	vm := Template_created_VM_override{
		OS:        template[OS_FIELD].(string),
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
	if d.Get(TEMPLATE_FIELD) != "" {
		override_file.WriteString(d.Get(TEMPLATE_FIELD).(string))
		override_file.WriteString("_template_override.tf.json")
		if _, err := os.Stat(override_file.String()); os.IsNotExist(err) {
			logger := LoggerCreate("CreateTemplateOverrideConfig_" +
				d.Get(TEMPLATE_FIELD).(string) + "_.log")
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
			vm_fields_map := map[string]interface{}{"template-server": vm}
			vm_map := map[string]interface{}{"sewan_clouddc_vm": vm_fields_map}
			resources_map := map[string]interface{}{"resource": vm_map}
			vm_json, _ := json.Marshal(resources_map)
			write_override_file_err = ioutil.WriteFile(override_file.String(),
				vm_json, 0644)
		}
	} else {
		write_override_file_err = errors.New("Template field is empty, " +
			"can not create a template override configuration.")
	}

	return write_override_file_err
}
