package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"os"
)

//------------------------------------------------------------------------------
func CompareJsonAndMap(jsonFile string,
	fileDataMap map[string]interface{}) (string, error) {
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
	return diffs, err
}

type TemplaterDummy struct{}

func (templaterFake TemplaterDummy) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (templaterFake TemplaterDummy) ValidateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake TemplaterDummy) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake TemplaterDummy) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type UnexistingTemplate_TemplaterFake struct{}

func (templaterFake UnexistingTemplate_TemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, errors.New("Unavailable template : windows95")
}
func (templaterFake UnexistingTemplate_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake UnexistingTemplate_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake UnexistingTemplate_TemplaterFake) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type Template_FormatError_TemplaterFake struct{}

func (templaterFake Template_FormatError_TemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (templaterFake Template_FormatError_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {
	return errors.New("Template missing fields : " + "\"" + NameField + "\" " +
		"\"" + OsField + "\" " +
		"\"" + RamField + "\" " +
		"\"" + CpuField + "\" " +
		"\"" + EnterpriseField + "\" " +
		"\"" + DisksField + "\" " +
		"\"" + DatacenterField + "\" ")
}
func (templaterFake Template_FormatError_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake Template_FormatError_TemplaterFake) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type existingTemplateNoAdditionalDiskVmMap_TemplaterFake struct{}

func (templaterFake existingTemplateNoAdditionalDiskVmMap_TemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IdField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{NameField: "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{VlanNameField: "unit test vlan1",
				MacAdressField: "00:50:56:21:7c:ab",
				ConnectedField: true,
			},
			map[string]interface{}{VlanNameField: "unit test vlan2",
				MacAdressField: "00:50:56:21:7c:ac",
				ConnectedField: true,
			},
		},
		"login":      "",
		"password":   "",
		DynamicField: "",
	}, nil
}
func (templaterFake existingTemplateNoAdditionalDiskVmMap_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake existingTemplateNoAdditionalDiskVmMap_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "Unit test template no disc add on vm resource")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RamField, 1)
	d.Set(CpuField, 1)
	d.Set(DisksField,
		[]interface{}{
			map[string]interface{}{NameField: "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
	)
	d.Set(NicsField, []interface{}{})
	return nil
}
func (templaterFake existingTemplateNoAdditionalDiskVmMap_TemplaterFake) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake struct{}

func (templaterFake instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IdField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{NameField: "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{VlanNameField: "unit test vlan1",
				MacAdressField: "00:50:56:21:7c:ab",
				ConnectedField: true,
			},
			map[string]interface{}{VlanNameField: "unit test vlan2",
				MacAdressField: "00:50:56:21:7c:ac",
				ConnectedField: true,
			},
		},
		"login":      "",
		"password":   "",
		DynamicField: "",
	}, nil
}
func (templaterFake instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "instanceNumberFieldUnitTest")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RamField, 1)
	d.Set(CpuField, 1)
	d.Set(DisksField,
		[]interface{}{
			map[string]interface{}{NameField: "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
	)
	d.Set(NicsField, []interface{}{})
	return nil
}
func (templaterFake instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake struct{}

func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IdField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{NameField: "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{VlanNameField: "unit test vlan1",
				MacAdressField: "00:50:56:21:7c:ab",
				ConnectedField: true,
			},
			map[string]interface{}{VlanNameField: "unit test vlan2",
				MacAdressField: "00:50:56:21:7c:ac",
				ConnectedField: true,
			},
		},
		"login":      "",
		"password":   "",
		DynamicField: "",
	}, nil
}
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake) ValidateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake) UpdateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RamField, 8)
	d.Set(CpuField, 4)
	d.Set(DisksField,
		[]interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
				SlugField:         "",
			},
			map[string]interface{}{
				NameField:         "template1 disk1",
				SizeField:         25,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
			},
		},
	)
	d.Set(NicsField, []interface{}{
		map[string]interface{}{
			VlanNameField:  "non template vlan 1",
			MacAdressField: "00:21:21:21:21:21",
			ConnectedField: true,
		},
		map[string]interface{}{
			VlanNameField:  "non template vlan 2",
			MacAdressField: "00:21:21:21:21:22",
			ConnectedField: true,
		},
	},
	)
	d.Set("vdc:          ", "vdc")
	d.Set("boot:         ", "on disk")
	d.Set(StorageClassField, "storage_enterprise")
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
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake) CreateVmTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}
