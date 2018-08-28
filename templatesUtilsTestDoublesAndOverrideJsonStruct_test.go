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
func compareJSONAndMap(jsonFile string,
	fileDataMap map[string]interface{}) (string, error) {
	var (
		err        error
		jsonData   []byte
		diffs      string
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
func (templaterFake TemplaterDummy) validateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake TemplaterDummy) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake TemplaterDummy) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type UnexistingTemplateTemplaterFake struct{}

func (templaterFake UnexistingTemplateTemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, errors.New("Unavailable template : windows95")
}
func (templaterFake UnexistingTemplateTemplaterFake) validateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake UnexistingTemplateTemplaterFake) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake UnexistingTemplateTemplaterFake) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type TemplateFormatErrorTemplaterFake struct{}

func (templaterFake TemplateFormatErrorTemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (templaterFake TemplateFormatErrorTemplaterFake) validateTemplate(template map[string]interface{}) error {
	return errors.New("Template missing fields : " + "\"" + NameField + "\" " +
		"\"" + OsField + "\" " +
		"\"" + RAMField + "\" " +
		"\"" + CPUField + "\" " +
		"\"" + EnterpriseField + "\" " +
		"\"" + DisksField + "\" " +
		"\"" + DatacenterField + "\" ")
}
func (templaterFake TemplateFormatErrorTemplaterFake) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	return nil
}
func (templaterFake TemplateFormatErrorTemplaterFake) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type existingTemplateNoAdditionalDiskVMMapTemplaterFake struct{}

func (templaterFake existingTemplateNoAdditionalDiskVMMapTemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IDField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RAMField:        1,
		CPUField:        1,
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
func (templaterFake existingTemplateNoAdditionalDiskVMMapTemplaterFake) validateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake existingTemplateNoAdditionalDiskVMMapTemplaterFake) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "Unit test template no disc add on vm resource")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RAMField, 1)
	d.Set(CPUField, 1)
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
func (templaterFake existingTemplateNoAdditionalDiskVMMapTemplaterFake) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake struct{}

func (templaterFake instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IDField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RAMField:        1,
		CPUField:        1,
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
func (templaterFake instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake) validateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "instanceNumberFieldUnitTest")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RAMField, 1)
	d.Set(CPUField, 1)
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
func (templaterFake instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}

type existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake struct{}

func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake) FetchTemplateFromList(templateName string,
	templateList []interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		IDField:         82,
		NameField:       "template1",
		SlugField:       "centos7-rd-dc1",
		RAMField:        1,
		CPUField:        1,
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
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake) validateTemplate(template map[string]interface{}) error {
	return nil
}
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake) updateSchemaFromTemplateOnResourceCreation(d *schema.ResourceData,
	template map[string]interface{}) error {
	d.Set(NameField, "existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap")
	d.Set(EnterpriseField, "unit test enterprise")
	d.Set(TemplateField, "template1")
	d.Set(RAMField, 8)
	d.Set(CPUField, 4)
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
func (templaterFake existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake) createVMTemplateOverrideConfig(d *schema.ResourceData,
	template map[string]interface{}) (string, error) {
	return "", nil
}
