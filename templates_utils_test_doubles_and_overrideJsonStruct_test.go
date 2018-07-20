package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
	//"reflect"
	"github.com/google/go-cmp/cmp"
)

//------------------------------------------------------------------------------
func CompareJsonAndMap(jsonFile string,
	fileDataMap map[string]interface{}) (error, string) {
	var (
		err        error
		jsonData   []byte
		diffs      string = ""
		dataStruct interface{}
	)
	jsonData, err = ioutil.ReadFile(jsonFile)
	if jsonFile != "" {
		err = os.Remove(jsonFile)
	}
	if err == nil {
		err = json.Unmarshal(jsonData, &dataStruct)
		if err == nil {
			diffs = cmp.Diff(fileDataMap, dataStruct)

		}
	}
	return err, diffs
}

//------------------------------------------------------------------------------
type TemplaterDummy struct{}

func (templaterFake TemplaterDummy) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	return nil, nil
}
func (templaterFake TemplaterDummy) ValidateTemplate(template map[string]interface{}) error {

	return nil
}
func (templaterFake TemplaterDummy) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {

	return nil
}
func (templaterFake TemplaterDummy) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	return nil, ""
}

//------------------------------------------------------------------------------
type Unexisting_template_TemplaterFake struct{}

func (templaterFake Unexisting_template_TemplaterFake) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	return nil, errors.New("Unavailable template : windows95")
}
func (templaterFake Unexisting_template_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {

	return nil
}
func (templaterFake Unexisting_template_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {

	return nil
}
func (templaterFake Unexisting_template_TemplaterFake) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	return nil, ""
}

//------------------------------------------------------------------------------
type Template_Format_error_TemplaterFake struct{}

func (templaterFake Template_Format_error_TemplaterFake) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	return nil, nil
}
func (templaterFake Template_Format_error_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {

	return errors.New("Template missing fields : " + "\"" + NAME_FIELD + "\" " +
		"\"" + OS_FIELD + "\" " +
		"\"" + RAM_FIELD + "\" " +
		"\"" + CPU_FIELD + "\" " +
		"\"" + ENTERPRISE_FIELD + "\" " +
		"\"" + DISKS_FIELD + "\" " +
		"\"" + DATACENTER_FIELD + "\" ")
}
func (templaterFake Template_Format_error_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {

	return nil
}
func (templaterFake Template_Format_error_TemplaterFake) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	return nil, ""
}

//------------------------------------------------------------------------------
type EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake struct{}

func (templaterFake EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	return map[string]interface{}{
		ID_FIELD:         82,
		NAME_FIELD:       "template1",
		SLUG_FIELD:       "centos7-rd-dc1",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "CentOS",
		ENTERPRISE_FIELD: "unit test enterprise",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{NAME_FIELD: "template1 disk1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan1",
				MAC_ADRESS_FIELD: "00:50:56:21:7c:ab",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan2",
				MAC_ADRESS_FIELD: "00:50:56:21:7c:ac",
				CONNECTED_FIELD:  true,
			},
		},
		"login":       "",
		"password":    "",
		DYNAMIC_FIELD: "",
	}, nil
}
func (templaterFake EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {

	return nil
}
func (templaterFake EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {

	d.Set(NAME_FIELD, "Unit test template no disc add on vm resource")
	d.Set(ENTERPRISE_FIELD, "unit test enterprise")
	d.Set(TEMPLATE_FIELD, "template1")
	d.Set(RAM_FIELD, 1)
	d.Set(CPU_FIELD, 1)
	d.Set(DISKS_FIELD,
		[]interface{}{
			map[string]interface{}{NAME_FIELD: "template1 disk1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
	)
	d.Set(NICS_FIELD, []interface{}{})
	return nil
}
func (templaterFake EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	return nil, ""
}

//------------------------------------------------------------------------------
type EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake struct{}

func (templaterFake EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake) FetchTemplateFromList(template_name string,
	templateList []interface{}) (map[string]interface{}, error) {

	return map[string]interface{}{
		ID_FIELD:         82,
		NAME_FIELD:       "template1",
		SLUG_FIELD:       "centos7-rd-dc1",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "CentOS",
		ENTERPRISE_FIELD: "unit test enterprise",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{NAME_FIELD: "template1 disk1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan1",
				MAC_ADRESS_FIELD: "00:50:56:21:7c:ab",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan2",
				MAC_ADRESS_FIELD: "00:50:56:21:7c:ac",
				CONNECTED_FIELD:  true,
			},
		},
		"login":       "",
		"password":    "",
		DYNAMIC_FIELD: "",
	}, nil
}
func (templaterFake EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {

	return nil
}
func (templaterFake EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {

	d.Set(NAME_FIELD, "EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP")
	d.Set(ENTERPRISE_FIELD, "unit test enterprise")
	d.Set(TEMPLATE_FIELD, "template1")
	d.Set(RAM_FIELD, 8)
	d.Set(CPU_FIELD, 4)
	d.Set(DISKS_FIELD,
		[]interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "",
			},
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          25,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
	)
	d.Set(NICS_FIELD, []interface{}{
		map[string]interface{}{
			VLAN_NAME_FIELD:  "non template vlan 1",
			MAC_ADRESS_FIELD: "00:21:21:21:21:21",
			CONNECTED_FIELD:  true,
		},
		map[string]interface{}{
			VLAN_NAME_FIELD:  "non template vlan 2",
			MAC_ADRESS_FIELD: "00:21:21:21:21:22",
			CONNECTED_FIELD:  true,
		},
	},
	)
	d.Set("vdc:          ", VDC_FIELD)
	d.Set("boot:         ", "on disk")
	d.Set(STORAGE_CLASS_FIELD, "storage_enterprise")
	d.Set("slug:         ", "42")
	d.Set("token:        ", "424242")
	d.Set("backup:       ", "backup_no_backup")
	d.Set("disk_image:   ", "")
	d.Set("platform_name:", "42")
	d.Set("backup_size:  ", 42)
	d.Set("comment:      ", "42")
	d.Set("dynamic_field:", "42")
	return nil
}
func (templaterFake EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake) CreateTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (error, string) {
	return nil, ""
}
