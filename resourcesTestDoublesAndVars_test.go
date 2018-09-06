package sewansdk

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

var (
	resourceOverrideJSONMap = map[string]interface{}{
		"resource": map[string]interface{}{
			"sewan_clouddc_vm": map[string]interface{}{
				"createVMTemplateOverrideConfig Unit test": map[string]interface{}{
					"name": "createVMTemplateOverrideConfig Unit test",
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
	resourceN42OverrideJSONMap = map[string]interface{}{
		"resource": map[string]interface{}{
			"sewan_clouddc_vm": map[string]interface{}{
				"createVMTemplateOverrideConfig Unit test": map[string]interface{}{
					"name": "createVMTemplateOverrideConfig Unit test-${count.index + 1}",
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
	nonExistingErrorVMSchemaMap = map[string]interface{}{
		IDField:         "an id, same behaviour if it's an int or float",
		NameField:       "VM schema update unit test",
		TemplateField:   "template1",
		RAMField:        2,
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
		RAMField:        2,
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
		CPUField:        1,
		RAMField:        2,
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
		NameField:       "createVMTemplateOverrideConfig Unit test",
		RAMField:        1,
		CPUField:        1,
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
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":nil}",
	}
	vmCreationN42FromTemplate1Schema = map[string]interface{}{
		NameField:           "createVMTemplateOverrideConfig Unit test",
		InstanceNumberField: 42,
		RAMField:            1,
		CPUField:            1,
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
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":nil}",
	}
	vmCreationFromTemplate1SchemaPreCreationWrongNicsInitMap = map[string]interface{}{
		RAMField:        1,
		CPUField:        1,
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
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":nil}",
	}
	vdcCreationMap = map[string]interface{}{
		NameField:       "Unit test vdc resource",
		EnterpriseField: "enterprise",
		VdcResourceField: []interface{}{
			map[string]interface{}{
				ResourceField: RAMField,
				TotalField:    20,
			},
			map[string]interface{}{
				ResourceField: CPUField,
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
				ResourceField: RAMField,
				TotalField:    20,
				UsedField:     0,
				SlugField:     "",
			},
			map[string]interface{}{
				ResourceField: CPUField,
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
				ResourceField: RAMField,
				UsedField:     0,
				TotalField:    20,
				SlugField:     "unit test enterprise-mono-ram",
			},
			map[string]interface{}{
				ResourceField: CPUField,
				UsedField:     0,
				TotalField:    1,
				SlugField:     "unit test enterprise-mono-cpu",
			},
			map[string]interface{}{
				ResourceField: "storage_enterprise",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "unit test enterprise-mono-storage_enterprise",
			},
			map[string]interface{}{
				ResourceField: "storage_performance",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "unit test enterprise-mono-storage_performance",
			},
			map[string]interface{}{
				ResourceField: "storage_high_performance",
				UsedField:     0,
				TotalField:    10,
				SlugField:     "unit test enterprise-dc1-mono-storage_high_performance",
			},
		},
		SlugField:    "unit test enterprise-dc1-vdc_te",
		DynamicField: "",
	}
	noTemplateVMMap = map[string]interface{}{
		NameField:  "Unit test no template vm resource",
		StateField: "UP",
		OsField:    "Debian",
		RAMField:   8,
		CPUField:   4,
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
		VdcField:          "vdc unit test",
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
	existingTemplateNoAdditionalDiskVMMap = map[string]interface{}{
		NameField:         "Unit test template no disc add on vm resource",
		EnterpriseField:   "unit test enterprise",
		TemplateField:     "template1",
		StateField:        "UP",
		VdcField:          "vdc unit test",
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
	}
	instanceNumberFieldUnitTestVMInstance = map[string]interface{}{
		NameField:           "instanceNumberFieldUnitTest",
		InstanceNumberField: 42,
		EnterpriseField:     "unit test enterprise",
		TemplateField:       "template1",
		StateField:          "UP",
		VdcField:            "vdc unit test",
		BootField:           "on disk",
		StorageClassField:   "storage_enterprise",
		SlugField:           "42",
		TokenField:          "424242",
		BackupField:         "backup_no_backup",
	}
	existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap = map[string]interface{}{
		NameField:       "existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap",
		EnterpriseField: "unit test enterprise",
		TemplateField:   "template1",
		StateField:      "UP",
		OsField:         "Debian",
		RAMField:        8,
		CPUField:        4,
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
		VdcField:          "vdc unit test",
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
	nonExistingTemplateVMMap = map[string]interface{}{
		NameField:         "windows95 vm",
		EnterpriseField:   "unit test enterprise",
		TemplateField:     "windows95",
		StateField:        "UP",
		RAMField:          8,
		CPUField:          4,
		VdcField:          "vdc unit test",
		BootField:         "on disk",
		StorageClassField: "storage_enterprise",
		SlugField:         "42",
		TokenField:        "424242",
		BackupField:       "backup_no_backup",
		DiskImageField:    "",
	}
	template2Map = map[string]interface{}{
		IDField:         40,
		NameField:       "template2",
		SlugField:       "unit test disk goulouglougoulouglou",
		RAMField:        1,
		CPUField:        1,
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
		IDField:         69,
		NameField:       "lastTemplate",
		SlugField:       "lastTemplate-slug",
		RAMField:        1,
		CPUField:        1,
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
		IDField:         82,
		NameField:       "template1",
		SlugField:       "template1 slug",
		RAMField:        1,
		CPUField:        1,
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
		IDField:         82,
		NameField:       "template1",
		SlugField:       "template1 slug",
		RAMField:        1,
		CPUField:        1,
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
	nonCriticalResourceMetaDataList = []interface{}{
		map[string]interface{}{
			"id":            4,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-05-28T12:28:42+02:00",
			"cos":           "Mono",
			"name":          "ram",
			"used":          324,
			"total":         350,
			"slug":          "unit-test-enterprise-mono-ram",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            5,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-05-28T12:28:32+02:00",
			"cos":           "Mono",
			"name":          "cpu",
			"used":          275,
			"total":         300,
			"slug":          "unit-test-enterprise-mono-cpu",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            6,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-14T17:32:15+01:00",
			"cos":           "Mono",
			"name":          "storage_enterprise",
			"used":          7708,
			"total":         8000,
			"slug":          "unit-test-enterprise-mono-storage_enterprise",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            7,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-07-31T15:55:06+02:00",
			"cos":           "Mono",
			"name":          "storage_performance",
			"used":          630,
			"total":         700,
			"slug":          "unit-test-enterprise-mono-storage_performance",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            8,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-06T11:02:17+01:00",
			"cos":           "Mono",
			"name":          "storage_high_performance",
			"used":          10,
			"total":         20,
			"slug":          "unit-test-enterprise-mono-storage_high_performance",
			"dynamic_field": nil,
			"service":       1,
		},
	}
	criticalResourceMetaDataList = []interface{}{
		map[string]interface{}{
			"id":            305,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-10-10T12:19:51+02:00",
			"modified":      "2017-10-10T12:19:51+02:00",
			"cos":           "HA",
			"name":          "cpu",
			"used":          5,
			"total":         10,
			"slug":          "resource-cpu-rd-ha",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            306,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-10-10T12:20:11+02:00",
			"modified":      "2017-10-10T12:20:11+02:00",
			"cos":           "HA",
			"name":          "ram",
			"used":          5,
			"total":         10,
			"slug":          "ram-ha-rd",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            314,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-03T16:09:32+02:00",
			"modified":      "2018-04-24T15:50:56+02:00",
			"cos":           "HA",
			"name":          "storage_enterprise",
			"used":          60,
			"total":         100,
			"slug":          "storage_enterprise-ha",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            315,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-24T12:35:55+02:00",
			"modified":      "2018-04-24T15:51:04+02:00",
			"cos":           "HA",
			"name":          "storage_performance",
			"used":          55,
			"total":         100,
			"slug":          "unit-test-enterprise-ha-storage_performance",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            316,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-24T12:36:15+02:00",
			"modified":      "2018-04-24T15:51:13+02:00",
			"cos":           "HA",
			"name":          "storage_high_performance",
			"used":          0,
			"total":         100,
			"slug":          "unit-test-enterprise-ha-storage_high_performance",
			"dynamic_field": nil,
			"service":       1,
		},
	}
	otherResourceMetaDataList = []interface{}{
		map[string]interface{}{
			"id":            55,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2017-08-10T05:01:03+02:00",
			"cos":           "Global",
			"name":          "backup",
			"used":          220,
			"total":         220,
			"slug":          "unit-test-enterprise-clouddc-backup",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            196,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-21T12:45:28+01:00",
			"cos":           "Global",
			"name":          "license_win_server",
			"used":          7,
			"total":         20,
			"slug":          "unit-test-enterprise-global-license_win_server",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            313,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-02-15T18:39:17+01:00",
			"modified":      "2018-02-16T15:50:16+01:00",
			"cos":           "Global",
			"name":          "license_redhat",
			"used":          2,
			"total":         3,
			"slug":          "sewan-rd-cloud-daatcenter-vdc-rd-licence-redhat",
			"dynamic_field": nil,
			"service":       1,
		},
	}
	resourceMetaDataList = []interface{}{
		map[string]interface{}{
			"id":            4,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-05-28T12:28:42+02:00",
			"cos":           "Mono",
			"name":          "ram",
			"used":          324,
			"total":         350,
			"slug":          "unit-test-enterprise-mono-ram",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            5,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-05-28T12:28:32+02:00",
			"cos":           "Mono",
			"name":          "cpu",
			"used":          275,
			"total":         300,
			"slug":          "unit-test-enterprise-mono-cpu",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            6,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-14T17:32:15+01:00",
			"cos":           "Mono",
			"name":          "storage_enterprise",
			"used":          7708,
			"total":         8000,
			"slug":          "unit-test-enterprise-mono-storage_enterprise",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            7,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-07-31T15:55:06+02:00",
			"cos":           "Mono",
			"name":          "storage_performance",
			"used":          630,
			"total":         700,
			"slug":          "unit-test-enterprise-mono-storage_performance",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            8,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-06T11:02:17+01:00",
			"cos":           "Mono",
			"name":          "storage_high_performance",
			"used":          10,
			"total":         20,
			"slug":          "unit-test-enterprise-mono-storage_high_performance",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            55,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2017-08-10T05:01:03+02:00",
			"cos":           "Global",
			"name":          "backup",
			"used":          220,
			"total":         220,
			"slug":          "unit-test-enterprise-clouddc-backup",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            196,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-06-29T12:10:35+02:00",
			"modified":      "2018-02-21T12:45:28+01:00",
			"cos":           "Global",
			"name":          "license_win_server",
			"used":          7,
			"total":         20,
			"slug":          "unit-test-enterprise-global-license_win_server",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            305,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-10-10T12:19:51+02:00",
			"modified":      "2017-10-10T12:19:51+02:00",
			"cos":           "HA",
			"name":          "cpu",
			"used":          5,
			"total":         10,
			"slug":          "resource-cpu-rd-ha",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            306,
			"enterprise":    "unit-test-enterprise",
			"created":       "2017-10-10T12:20:11+02:00",
			"modified":      "2017-10-10T12:20:11+02:00",
			"cos":           "HA",
			"name":          "ram",
			"used":          5,
			"total":         10,
			"slug":          "ram-ha-rd",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            313,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-02-15T18:39:17+01:00",
			"modified":      "2018-02-16T15:50:16+01:00",
			"cos":           "Global",
			"name":          "license_redhat",
			"used":          2,
			"total":         3,
			"slug":          "sewan-rd-cloud-daatcenter-vdc-rd-licence-redhat",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            314,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-03T16:09:32+02:00",
			"modified":      "2018-04-24T15:50:56+02:00",
			"cos":           "HA",
			"name":          "storage_enterprise",
			"used":          60,
			"total":         100,
			"slug":          "storage_enterprise-ha",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            315,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-24T12:35:55+02:00",
			"modified":      "2018-04-24T15:51:04+02:00",
			"cos":           "HA",
			"name":          "storage_performance",
			"used":          55,
			"total":         100,
			"slug":          "unit-test-enterprise-ha-storage_performance",
			"dynamic_field": nil,
			"service":       1,
		},
		map[string]interface{}{
			"id":            316,
			"enterprise":    "unit-test-enterprise",
			"created":       "2018-04-24T12:36:15+02:00",
			"modified":      "2018-04-24T15:51:13+02:00",
			"cos":           "HA",
			"name":          "storage_high_performance",
			"used":          0,
			"total":         100,
			"slug":          "unit-test-enterprise-ha-storage_high_performance",
			"dynamic_field": nil,
			"service":       1,
		},
	}
	templatesList = []interface{}{
		template2Map,
		template1Map,
		map[string]interface{}{
			IDField:         41,
			NameField:       "template3",
			SlugField:       "unit test template3 slug",
			RAMField:        1,
			CPUField:        1,
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
			IDField:         43,
			NameField:       "template4",
			SlugField:       "tpl-centos7-rd",
			RAMField:        1,
			CPUField:        1,
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
			IDField:         58,
			NameField:       "template windaube7",
			SlugField:       "slug windows7",
			RAMField:        1,
			CPUField:        1,
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
	testUpdateVMMap = map[string]interface{}{
		IDField:    "unit test vm",
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RAMField:   16,
		CPUField:   8,
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
		VdcField:          "vdc update",
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
	testUpdateVMMapIntID = map[string]interface{}{
		IDField:    1212,
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RAMField:   16,
		CPUField:   8,
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
		VdcField:          "vdc update",
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
	testUpdateVMMapFloatID = map[string]interface{}{
		IDField:    121212.12,
		NameField:  "Unit test vm",
		StateField: "DOWN",
		OsField:    "CentOS",
		RAMField:   16,
		CPUField:   8,
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
		VdcField:          "vdc update",
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

func resourceVMDisk() *schema.Resource {
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

func resourceVMNic() *schema.Resource {
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

func resourceVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NameField: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			InstanceNumberField: &schema.Schema{
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
			RAMField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			CPUField: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			DisksField: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVMDisk(),
			},
			NicsField: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVMNic(),
			},
			VdcField: &schema.Schema{
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
			OutsourcingField: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fakeVdcInstanceVdcCreationMap() vdcStruct {
	return vdcStruct{
		Name:       "Unit test vdc resource",
		Enterprise: "enterprise",
		VdcResources: []interface{}{
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

func vdcInstanceFake() vdcStruct {
	return vdcStruct{
		Name:       "Unit test vdc resource",
		Enterprise: "Unit Test value",
		Datacenter: "Unit Test value",
		VdcResources: []interface{}{
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

func vmInstanceNoTemplateVMMap() vmStruct {
	return vmStruct{
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
		Vdc:          "vdc unit test",
		Boot:         "on disk",
		StorageClass: "storage_enterprise",
		Slug:         "42",
		Token:        "424242",
		Backup:       "backup_no_backup",
		DiskImage:    "",
		PlatformName: "42",
		BackupSize:   42,
		Comment:      "",
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"\",\"TemplateDisksOnCreation\":null}",
	}
}

func fakeVMInstanceExistingTemplateNoAdditionalDiskVMMap() vmStruct {
	return vmStruct{
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
		Nics:         []interface{}{},
		Vdc:          "vdc unit test",
		Boot:         "on disk",
		StorageClass: "storage_enterprise",
		Slug:         "42",
		Token:        "424242",
		Backup:       "backup_no_backup",
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func fakeVMInstanceInstanceNumberFieldUnitTestVMInstanceMAP() vmStruct {
	return vmStruct{
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
		Nics:         []interface{}{},
		Vdc:          "vdc unit test",
		Boot:         "on disk",
		StorageClass: "storage_enterprise",
		Slug:         "42",
		Token:        "424242",
		Backup:       "backup_no_backup",
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func fakeVMInstanceExistingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap() vmStruct {
	return vmStruct{
		Name:       "existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap-0",
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
		Vdc:          "vdc unit test",
		Boot:         "on disk",
		StorageClass: "storage_enterprise",
		Slug:         "42",
		Token:        "424242",
		Backup:       "backup_no_backup",
		DiskImage:    "",
		PlatformName: "42",
		BackupSize:   42,
		Comment:      "",
		DynamicField: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"TemplateDisksOnCreation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func vmInstanceNoTemplateFake() vmStruct {
	return vmStruct{
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
		Vdc:          "Unit Test value",
		Boot:         "Unit Test value",
		StorageClass: "Unit Test value",
		Slug:         "Unit Test value",
		Token:        "Unit Test value",
		Backup:       "Unit Test value",
		DiskImage:    "Unit Test value",
		PlatformName: "Unit Test value",
		BackupSize:   42,
		Comment:      "",
	}
}

func vdcSchemaInit(vdc map[string]interface{}) *schema.ResourceData {
	d := resourceVdc().TestResourceData()
	schemaTooler := SchemaTooler{
		SchemaTools: SchemaSchemaer{},
	}
	schemaTooler.SchemaTools.UpdateLocalResourceState(vdc, d, &schemaTooler)
	return d
}

func vmSchemaInit(vm map[string]interface{}) *schema.ResourceData {
	d := resourceVM().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: SchemaSchemaer{},
	}
	schemaTooler.SchemaTools.UpdateLocalResourceState(vm, d, &schemaTooler)

	return d
}

func resource(resourceType string) *schema.Resource {

	resource := &schema.Resource{}
	switch resourceType {
	case VdcResourceType:
		resource = resourceVdc()
	case VMResourceType:
		resource = resourceVM()
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

type RespBody struct {
	Detail string `json:"detail"`
}

func JSONStub() map[string]interface{} {

	var jsonStub interface{}
	simpleJSON, _ := json.Marshal(RespBody{Detail: "a simple json"})
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(simpleJSON))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonStub)

	return jsonStub.(map[string]interface{})
}

func JSONTemplateListFake() []interface{} {
	var jsonFake interface{}
	fakeJSON, _ := json.Marshal(templatesList)
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(fakeJSON))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonFake)

	return jsonFake.([]interface{})
}
