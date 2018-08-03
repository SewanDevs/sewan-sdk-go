package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

var (
	resourceOverrideJsonMap = map[string]interface{}{
		"resource": map[string]interface{}{
			"sewan_clouddc_vm": map[string]interface{}{
				"CreateTemplateOverrideConfig Unit test": map[string]interface{}{
					"name": "CreateTemplateOverrideConfig Unit test",
					"disks": []interface{}{
						map[string]interface{}{
							"name":          "unit test disk template1",
							"size":          float64(20),
							"storage_class": "storage_enterprise",
						},
					},
					"disk_image": "",
					"os":         "CentOS",
					"ram":        float64(1),
					"cpu":        float64(1),
					"backup":     "",
					"nics": []interface{}{
						map[string]interface{}{
							"vlan":      "unit test vlan1",
							"connected": true,
						},
						map[string]interface{}{
							"vlan":      "unit test vlan2",
							"connected": true,
						},
					},
					"vdc":  "",
					"boot": "",
				},
			},
		},
	}
	resourceN42OverrideJsonMap = map[string]interface{}{
		"resource": map[string]interface{}{
			"sewan_clouddc_vm": map[string]interface{}{
				"CreateTemplateOverrideConfig Unit test": map[string]interface{}{
					"name": "CreateTemplateOverrideConfig Unit test-${count.index + 1}",
					"disks": []interface{}{
						map[string]interface{}{
							"name":          "unit test disk template1",
							"size":          float64(20),
							"storage_class": "storage_enterprise",
						},
					},
					"disk_image": "",
					"os":         "CentOS",
					"ram":        float64(1),
					"cpu":        float64(1),
					"backup":     "",
					"nics": []interface{}{
						map[string]interface{}{
							"vlan":      "unit test vlan1",
							"connected": true,
						},
						map[string]interface{}{
							"vlan":      "unit test vlan2",
							"connected": true,
						},
					},
					"vdc":  "",
					"boot": "",
				},
			},
		},
	}
	nonExistingErrorVmSchemaMap = map[string]interface{}{
		IdField:         "an id, same behaviour if it's an int or float",
		NameField:       "VM schema update unit test",
		TemplateField:   "template1",
		RamField:        2,
		EnterpriseField: "unit test enterprise",
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vm additional unit test vlan1",
				ConnectedField: true,
			},
		},
		DynamicField: "",
	}
	vmSchemaMapPreUpdateFromTemplate = map[string]interface{}{
		NameField:       "VM schema update unit test",
		TemplateField:   "template1",
		RamField:        2,
		EnterpriseField: "unit test enterprise",
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vm additional unit test vlan1",
				ConnectedField: true,
			},
		},
		DynamicField: "",
	}
	vmSchemaMapPostUpdateFromTemplate = map[string]interface{}{
		NameField:       "VM schema update unit test",
		TemplateField:   "template1",
		CpuField:        1,
		RamField:        2,
		EnterpriseField: "unit test enterprise",
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "unit test vlan1",
				MacAdressField: "",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "unit test vlan2",
				MacAdressField: "",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "vm additional unit test vlan1",
				MacAdressField: "",
				ConnectedField: true,
			},
		},
		DynamicField: "",
	}
	vmCreationFromTemplate1Schema = map[string]interface{}{
		NameField:       "CreateTemplateOverrideConfig Unit test",
		RamField:        1,
		CpuField:        1,
		EnterpriseField: "unit test enterprise",
		TemplateField:   "template1",
		OsField:         "Debian",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk template1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk slug",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{VlanNameField: "unit test vlan1",
				ConnectedField: true,
			},
			map[string]interface{}{VlanNameField: "unit test vlan2",
				ConnectedField: true,
			},
		},
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":null}",
	}
	vmCreationN42FromTemplate1Schema = map[string]interface{}{
		NameField:           "CreateTemplateOverrideConfig Unit test",
		instanceNumberField: 42,
		RamField:            1,
		CpuField:            1,
		EnterpriseField:     "unit test enterprise",
		TemplateField:       "template1",
		OsField:             "Debian",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk template1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk slug",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{VlanNameField: "unit test vlan1",
				ConnectedField: true,
			},
			map[string]interface{}{VlanNameField: "unit test vlan2",
				ConnectedField: true,
			},
		},
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":null}",
	}
	vmCreationFromTemplate1SchemaPreCreationWrongNicsInitMap = map[string]interface{}{
		RamField:        1,
		CpuField:        1,
		EnterpriseField: "unit test enterprise",
		NameField:       "template1",
		OsField:         "Debian",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk template1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk slug",
			},
		},
		NicsField:    "Wrong nics type",
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":null}",
	}
	vdcCreationMap = map[string]interface{}{
		NameField:       "Unit test vdc resource",
		EnterpriseField: "enterprise",
		VdcResourceField: []interface{}{
			map[string]interface{}{
				ResourceField: RamField,
				TotalField:    20,
			},
			map[string]interface{}{
				ResourceField: CpuField,
				TotalField:    1,
			},
			map[string]interface{}{
				ResourceField: "storage_enterprise",
				TotalField:    10,
			},
			map[string]interface{}{
				ResourceField: "storage_performance",
				TotalField:    10,
			},
			map[string]interface{}{
				ResourceField: "storage_high_performance",
				TotalField:    10,
			},
		},
	}
	vdcResourcesNamesPreUpdateMap = map[string]interface{}{
		NameField:       "Unit test vdc resource",
		EnterpriseField: "enterprise",
		VdcResourceField: []interface{}{
			map[string]interface{}{
				ResourceField: "enterprise-mono-ram",
				TotalField:    20,
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-cpu",
				TotalField:    1,
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_enterprise",
				TotalField:    10,
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_performance",
				TotalField:    10,
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_high_performance",
				TotalField:    10,
			},
		},
	}
	vdcResourcesNamesUpdatedMap = map[string]interface{}{
		NameField:       "Unit test vdc resource",
		EnterpriseField: "enterprise",
		VdcResourceField: []interface{}{
			map[string]interface{}{
				ResourceField: RamField,
				TotalField:    20,
				UsedField:     0,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: CpuField,
				TotalField:    1,
				UsedField:     0,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "storage_enterprise",
				TotalField:    10,
				UsedField:     0,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "storage_performance",
				TotalField:    10,
				UsedField:     0,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "storage_high_performance",
				TotalField:    10,
				UsedField:     0,
				SlugField:     "",
			},
		},
	}
	vdcReadResponseMap = map[string]interface{}{
		NameField:       "Unit test vdc",
		EnterpriseField: "unit test enterprise",
		VdcResourceField: []interface{}{
			map[string]interface{}{
				ResourceField: RamField,
				UsedField:     "0",
				TotalField:    "20",
				SlugField:     "unit test enterprise-dc1-vdc_te-ram",
			},
			map[string]interface{}{
				ResourceField: CpuField,
				UsedField:     "0",
				TotalField:    "1",
				SlugField:     "unit test enterprise-dc1-vdc_te-cpu",
			},
			map[string]interface{}{
				ResourceField: "storage_enterprise",
				UsedField:     "0",
				TotalField:    "10",
				SlugField:     "unit test enterprise-dc1-vdc_te-storage_enterprise",
			},
			map[string]interface{}{
				ResourceField: "storage_performance",
				UsedField:     "0",
				TotalField:    "10",
				SlugField:     "unit test enterprise-dc1-vdc_te-storage_performance",
			},
			map[string]interface{}{
				ResourceField: "storage_high_performance",
				UsedField:     "0",
				TotalField:    "10",
				SlugField:     "unit test enterprise-dc1-vdc_te-storage_high_performance",
			},
		},
		SlugField:    "unit test enterprise-dc1-vdc_te",
		DynamicField: "",
	}
	noTemplateVmMap = map[string]interface{}{
		NameField:  "Unit test no template vm resource",
		StateField: "UP",
		OsField:    "Debian",
		RamField:   8,
		CpuField:   4,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan 1 update",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "vlan 2",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: true,
			},
		},
		BootField:         BootField,
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
		DiskImageField:    "",
		PlatformNameField: "42",
		BackupSizeField:   42,
		CommentField:      "42",
	}
	existingTemplateNoAdditionalDiskVmMap = map[string]interface{}{
		NameField:         "Unit test template no disc add on vm resource",
		EnterpriseField:   "unit test enterprise",
		TemplateField:     "template1",
		StateField:        "UP",
		BootField:         BootField,
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
	}
	instanceNumberFieldUnitTestVmInstance = map[string]interface{}{
		NameField:           "instanceNumberFieldUnitTest",
		instanceNumberField: 42,
		EnterpriseField:     "unit test enterprise",
		TemplateField:       "template1",
		StateField:          "UP",
		BootField:           BootField,
		BootField:           "on disk",
		StorageClassField:   "storage_enterprise",
		SlugField:           "42",
		TokenField:          "424242",
		BackupField:         "backup_no_backup",
	}
	existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap = map[string]interface{}{
		NameField:       "existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap",
		EnterpriseField: "unit test enterprise",
		TemplateField:   "template1",
		StateField:      "UP",
		OsField:         "Debian",
		RamField:        8,
		CpuField:        4,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
			},
		},
		NicsField: []interface{}{
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
		BootField:         BootField,
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
		DiskImageField:    "",
		PlatformNameField: "42",
		BackupSizeField:   42,
		CommentField:      "42",
	}
	existingTemplateWithDeletedDiskVmMap = map[string]interface{}{
		IdField:         "EXISTING_TEMPLATE_AND_VM_INSTANCE_WITH_DELETED_DISK_VM_MAP",
		NameField:       "existingTemplateWithDeletedDiskVmMap",
		EnterpriseField: "unit test enterprise",
		TemplateField:   "template1",
		StateField:      "UP",
		OsField:         "Debian",
		RamField:        8,
		CpuField:        4,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "template1 disk1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
			},
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
			},
		},
		NicsField:         []interface{}{},
		BootField:         BootField,
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
		DiskImageField:    "",
		PlatformNameField: "42",
		BackupSizeField:   42,
		CommentField:      "42",
		DynamicField:      "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":null}",
	}
	nonExistingTemplateVmMap = map[string]interface{}{
		NameField:         "windows95 vm",
		EnterpriseField:   "unit test enterprise",
		TemplateField:     "windows95",
		StateField:        "UP",
		RamField:          8,
		CpuField:          4,
		BootField:         BootField,
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
		DiskImageField:    "",
	}
	template2Map = map[string]interface{}{
		IdField:         40,
		NameField:       "template2",
		SlugField:       "unit test disk goulouglougoulouglou",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk goulouglouglou",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk goulouglou slug",
			},
		},
		NicsField:    []interface{}{},
		"login":      "",
		"password":   "",
		DynamicField: "",
	}
	lastTemplateInTemplatesList = map[string]interface{}{
		IdField:         69,
		NameField:       "lastTemplate",
		SlugField:       "lastTemplate-slug",
		RamField:        1,
		CpuField:        1,
		OsField:         "Debian",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk-debian9-rd-1",
				SizeField:         10,
				StorageClassField: "storage_enterprise",
				SlugField:         "disk-debian9-rd-1",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "unit test vlan1",
				MacAdressField: "00:50:56:00:01:de",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "unit test vlan2",
				MacAdressField: "00:50:56:00:01:df",
				ConnectedField: true,
			},
		},
		"login":      nil,
		"password":   nil,
		DynamicField: nil,
	}
	template1Map = map[string]interface{}{
		IdField:         82,
		NameField:       "template1",
		SlugField:       "template1 slug",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		BootField:       "on disk",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk template1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk slug",
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
	}
	template1MapBis = map[string]interface{}{
		IdField:         82,
		NameField:       "template1",
		SlugField:       "template1 slug",
		RamField:        1,
		CpuField:        1,
		OsField:         "CentOS",
		BootField:       "on disk",
		EnterpriseField: "unit test enterprise",
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "unit test disk template1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "unit test disk slug",
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
		DynamicField: "",
	}
	templatesList = []interface{}{
		template2Map,
		template1Map,
		map[string]interface{}{
			IdField:         41,
			NameField:       "template3",
			SlugField:       "unit test template3 slug",
			RamField:        1,
			CpuField:        1,
			OsField:         "Debian",
			EnterpriseField: "unit test enterprise",
			DisksField: []interface{}{
				map[string]interface{}{
					NameField:         "unit test disk2",
					SizeField:         20,
					StorageClassField: "storage_enterprise",
					SlugField:         "unit test disk slug 2",
				},
			},
			NicsField:    []interface{}{},
			"login":      "",
			"password":   "",
			DynamicField: "",
		},
		map[string]interface{}{
			IdField:         43,
			NameField:       "template4",
			SlugField:       "tpl-centos7-rd",
			RamField:        1,
			CpuField:        1,
			OsField:         "CentOS",
			EnterpriseField: "unit test enterprise",
			DisksField: []interface{}{
				map[string]interface{}{
					NameField:         "unit test disk 1",
					SizeField:         20,
					StorageClassField: "storage_enterprise",
					SlugField:         "unit test disk slug",
				},
			},
			NicsField: []interface{}{
				map[string]interface{}{
					VlanNameField:  "unit test vlan1",
					MacAdressField: "00:50:56:00:00:23",
					ConnectedField: true,
				},
				map[string]interface{}{
					VlanNameField:  "unit test vlan2",
					MacAdressField: "00:50:56:00:00:24",
					ConnectedField: true,
				},
			},
			"login":      nil,
			"password":   nil,
			DynamicField: nil,
		},
		map[string]interface{}{
			IdField:         58,
			NameField:       "template windaube7",
			SlugField:       "slug windows7",
			RamField:        1,
			CpuField:        1,
			OsField:         "Windows Serveur 64bits",
			EnterpriseField: "unit test enterprise",
			DisksField: []interface{}{
				map[string]interface{}{
					NameField:         "disk-Template-Windows",
					SizeField:         60,
					StorageClassField: "storage_enterprise",
					SlugField:         "disk-template-windows7",
				},
			},
			NicsField:    []interface{}{},
			"login":      nil,
			"password":   nil,
			DynamicField: nil,
		},
		lastTemplateInTemplatesList,
	}
	wrongTemplatesList = []interface{}{
		template2Map,
		"Wrongly formated template",
		lastTemplateInTemplatesList,
	}
	testUpdateVmMap = map[string]interface{}{
		IdField:    "unit test vm",
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RamField:   16,
		CpuField:   8,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
			map[string]interface{}{
				NameField:         "disk 2 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan 1 update",
				MacAdressField: "42",
				ConnectedField: false,
			},
		},
		BootField:         "vdc update",
		BootField:         "on disk update",
		StorageClassField: "storage_enterprise update",
		SlugField:         "42 update",
		TokenField:        "424242 update",
		BackupField:       "backup_no_backup update",
		DiskImageField:    " update",
		PlatformNameField: "",
		BackupSizeField:   42,
		CommentField:      "",
	}
	testUpdateVmMapIntId = map[string]interface{}{
		IdField:    1212,
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RamField:   16,
		CpuField:   8,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
			map[string]interface{}{
				NameField:         "disk 2 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan 1 update",
				MacAdressField: "42",
				ConnectedField: false,
			},
		},
		BootField:         "vdc update",
		BootField:         "on disk update",
		StorageClassField: "storage_enterprise update",
		SlugField:         "42 update",
		TokenField:        "424242 update",
		BackupField:       "backup_no_backup update",
		DiskImageField:    " update",
		PlatformNameField: "",
		BackupSizeField:   43,
		CommentField:      "",
	}
	testUpdateVmMapFloatId = map[string]interface{}{
		IdField:    121212.12,
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RamField:   16,
		CpuField:   8,
		DisksField: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
			map[string]interface{}{
				NameField:         "disk 2 update",
				SizeField:         42,
				StorageClassField: "StorageClassField update",
				SlugField:         "slug update",
				VDiskField:        "",
			},
		},
		NicsField: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan 1 update",
				MacAdressField: "42",
				ConnectedField: false,
			},
		},
		BootField:         "vdc update",
		BootField:         "on disk update",
		StorageClassField: "storage_enterprise update",
		SlugField:         "42 update",
		TokenField:        "424242 update",
		BackupField:       "backup_no_backup update",
		DiskImageField:    " update",
		PlatformNameField: "",
		BackupSizeField:   42,
		CommentField:      "",
	}
)

func resourceVdcResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			ResourceField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			UsedField: &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			TotalField: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			SlugField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVdc() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NameField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			EnterpriseField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			DatacenterField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			VdcResourceField: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVdcResource(),
			},
			SlugField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			DynamicField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVmDisk() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NameField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			SizeField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			StorageClassField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			SlugField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			VDiskField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVmNic() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			VlanNameField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			MacAdressField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			ConnectedField: &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceVm() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NameField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			instanceNumberField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			EnterpriseField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			TemplateField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			StateField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			OsField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			RamField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			CpuField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			DisksField: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVmDisk(),
			},
			NicsField: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVmNic(),
			},
			BootField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			BootField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			StorageClassField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			SlugField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			TokenField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			BackupField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			DiskImageField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			PlatformNameField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			BackupSizeField: &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			CommentField: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			DynamicField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			outsourcingField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func FakeVdcInstanceVdcCreationMap() VDC {
	return VDC{
		Name:       "Unit test vdc resource",
		Enterprise: "enterprise",
		Vdc_resources: []interface{}{
			map[string]interface{}{
				ResourceField: "enterprise-mono-ram",
				UsedField:     0,
				TotalField:    20,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-cpu",
				UsedField:     0,
				TotalField:    1,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_enterprise",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_performance",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: "enterprise-mono-storage_high_performance",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "",
			},
		},
		Slug:         "",
		DynamicField: "",
	}
}

func vdcInstanceFake() VDC {
	return VDC{
		Name:       "Unit test vdc resource",
		Enterprise: "Unit Test value",
		Datacenter: "Unit Test value",
		Vdc_resources: []interface{}{
			map[string]interface{}{
				ResourceField: "Resource1",
				UsedField:     1,
				TotalField:    2,
				SlugField:     "Unit Test value1",
			},
			map[string]interface{}{
				ResourceField: "Resource2",
				UsedField:     1,
				TotalField:    2,
				SlugField:     "Unit Test value2",
			},
			map[string]interface{}{
				ResourceField: "Resource3",
				UsedField:     1,
				TotalField:    2,
				SlugField:     "Unit Test value3",
			},
		},
		Slug:         "Unit Test value",
		DynamicField: "Unit Test value",
	}
}

func vmInstanceNoTemplateVmMap() VM {
	return VM{
		Name:  "Unit test no template vm resource",
		State: "UP",
		OS:    "Debian",
		RAM:   8,
		CPU:   4,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
				SlugField:         "",
				VDiskField:        "",
			},
		},
		Nics: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan 1 update",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "vlan 2",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: true,
			},
		},
		Vdc:           BootField,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"\",\"TemplateDisksOnCreation\":null}",
	}
}

func FakeVmInstanceExistingTemplateNoAdditionalDiskVmMap() VM {
	return VM{
		Name:       "Unit test template no disc add on vm resource-0",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		RAM:        1,
		CPU:        1,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
				VDiskField:        "",
			},
		},
		Nics:          []interface{}{},
		Vdc:           BootField,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstanceInstanceNumberFieldUnitTestVmInstance_MAP() VM {
	return VM{
		Name:       "instanceNumberFieldUnitTest-42",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		RAM:        1,
		CPU:        1,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "template1 disk1",
				SizeField:         20,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
				VDiskField:        "",
			},
		},
		Nics:          []interface{}{},
		Vdc:           BootField,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstanceExistingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap() VM {
	return VM{
		Name:       "existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap-0",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		OS:         "Debian",
		RAM:        8,
		CPU:        4,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
				SlugField:         "",
				VDiskField:        "",
			},
			map[string]interface{}{
				NameField:         "template1 disk1",
				SizeField:         25,
				StorageClassField: "storage_enterprise",
				SlugField:         "template1 disk1 slug",
				VDiskField:        "",
			},
		},
		Nics: []interface{}{
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
		Vdc:           BootField,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstanceExistingTemplateWithDeletedDiskVmMap() VM {
	return VM{
		Name:       "existingTemplateWithDeletedDiskVmMap",
		Enterprise: "unit test enterprise",
		Template:   "",
		State:      "UP",
		OS:         "Debian",
		RAM:        8,
		CPU:        4,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "disk 1",
				SizeField:         24,
				StorageClassField: "storage_enterprise",
				SlugField:         "",
				VDiskField:        "",
			},
		},
		Nics:          []interface{}{},
		Vdc:           BootField,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":null}",
	}
}

func vmInstanceNoTemplateFake() VM {
	return VM{
		Name:       "Unit test vm resource",
		Enterprise: "Unit Test value",
		Template:   "",
		State:      "Unit Test value",
		OS:         "Unit Test value",
		RAM:        1,
		CPU:        1,
		Disks: []interface{}{
			map[string]interface{}{
				NameField:         "name1",
				SizeField:         10,
				StorageClassField: "Unit Test value",
				VDiskField:        "",
			},
			map[string]interface{}{
				NameField:         "name2",
				SizeField:         10,
				StorageClassField: "Unit Test value",
				VDiskField:        "",
			},
		},
		Nics: []interface{}{
			map[string]interface{}{
				VlanNameField:  "vlan1",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: true,
			},
			map[string]interface{}{
				VlanNameField:  "vlan1",
				MacAdressField: "00:21:21:21:21:21",
				ConnectedField: false,
			},
		},
		Vdc:           "Unit Test value",
		Boot:          "Unit Test value",
		Storage_class: "Unit Test value",
		Slug:          "Unit Test value",
		Token:         "Unit Test value",
		Backup:        "Unit Test value",
		Disk_image:    "Unit Test value",
		Platform_name: "Unit Test value",
		Backup_size:   42,
		Comment:       "",
	}
}

func vdcSchemaInit(vdc map[string]interface{}) *schema.ResourceData {
	d := resourceVdc().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.UpdateLocalResourceState(vdc, d, &schemaTooler)

	return d
}

func vmSchemaInit(vm map[string]interface{}) *schema.ResourceData {
	d := resourceVm().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.UpdateLocalResourceState(vm, d, &schemaTooler)

	return d
}

func resource(resourceType string) *schema.Resource {

	resource := &schema.Resource{}
	switch resourceType {
	case VdcResourceType:
		resource = resourceVdc()
	case VmResourceType:
		resource = resourceVm()
	default:
		resource = &schema.Resource{
			Schema: map[string]*schema.Schema{
				NameField: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		}
	}
	return resource
}

type Resp_Body struct {
	Detail string `json:"detail"`
}

func JsonStub() map[string]interface{} {

	var jsonStub interface{}
	simpleJson, _ := json.Marshal(Resp_Body{Detail: "a simple json"})
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(simpleJson))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonStub)

	return jsonStub.(map[string]interface{})
}

func JsonTemplateListFake() []interface{} {
	var jsonFake interface{}
	fakeJson, _ := json.Marshal(templatesList)
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(fakeJson))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonFake)

	return jsonFake.([]interface{})
}
