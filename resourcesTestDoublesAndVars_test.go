package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

const (
	REQ_ERR                 = "Creation request response error."
	NOT_FOUND_STATUS        = "404 Not Found"
	NOT_FOUND_MSG           = "404 Not Found{\"detail\":\"Not found.\"}"
	UNAUTHORIZED_STATUS     = "401 Unauthorized"
	UNAUTHORIZED_MSG        = "401 Unauthorized{\"detail\":\"Token non valide.\"}"
	DESTROY_WRONG_MSG       = "{\"detail\":\"Destroying resource wrong body message\"}"
	CHECK_REDIRECT_FAILURE  = "CheckRedirectReqFailure"
	VDC_DESTROY_FAILURE_MSG = "Destroying the VDC now"
	VM_DESTROY_FAILURE_MSG  = "Destroying the VM now"
	WRONG_RESOURCE_TYPE     = "a_non_supportedResource_type"
	ENTERPRISE_SLUG         = "unit test enterprise"
)

var (
	RESOURCE_OVERRIDE_JSON_MAP = map[string]interface{}{
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
	RESOURCE_N42_OVERRIDE_JSON_MAP = map[string]interface{}{
		"resource": map[string]interface{}{
			"sewan_clouddc_vm": map[string]interface{}{
				"CreateTemplateOverrideConfig Unit test": map[string]interface{}{
					"name": "CreateTemplateOverrideConfig Unit test-${count.index}",
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
	NON_EXISTING_ERROR_VM_SCHEMA_MAP = map[string]interface{}{
		ID_FIELD:         "an id, same behaviour if it's an int or float",
		NAME_FIELD:       "VM schema update unit test",
		TEMPLATE_FIELD:   "template1",
		RAM_FIELD:        2,
		ENTERPRISE_FIELD: "unit test enterprise",
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD: "vm additional unit test vlan1",
				CONNECTED_FIELD: true,
			},
		},
		DYNAMIC_FIELD: "",
	}
	VM_SCHEMA_MAP_PRE_UPDATE_FROM_TEMPLATE = map[string]interface{}{
		NAME_FIELD:       "VM schema update unit test",
		TEMPLATE_FIELD:   "template1",
		RAM_FIELD:        2,
		ENTERPRISE_FIELD: "unit test enterprise",
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD: "vm additional unit test vlan1",
				CONNECTED_FIELD: true,
			},
		},
		DYNAMIC_FIELD: "",
	}
	VM_SCHEMA_MAP_POST_UPDATE_FROM_TEMPLATE = map[string]interface{}{
		NAME_FIELD:       "VM schema update unit test",
		TEMPLATE_FIELD:   "template1",
		CPU_FIELD:        1,
		RAM_FIELD:        2,
		ENTERPRISE_FIELD: "unit test enterprise",
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan1",
				MAC_ADRESS_FIELD: "",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan2",
				MAC_ADRESS_FIELD: "",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vm additional unit test vlan1",
				MAC_ADRESS_FIELD: "",
				CONNECTED_FIELD:  true,
			},
		},
		DYNAMIC_FIELD: "",
	}
	VM_CREATION_FROM_TEMPLATE1_SCHEMA = map[string]interface{}{
		NAME_FIELD:       "CreateTemplateOverrideConfig Unit test",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		ENTERPRISE_FIELD: "unit test enterprise",
		TEMPLATE_FIELD:   "template1",
		OS_FIELD:         "Debian",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk template1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk slug",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan1",
				CONNECTED_FIELD: true,
			},
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan2",
				CONNECTED_FIELD: true,
			},
		},
		DYNAMIC_FIELD: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":null}",
	}
	VM_CREATION_N42_FROM_TEMPLATE1_SCHEMA = map[string]interface{}{
		NAME_FIELD:            "CreateTemplateOverrideConfig Unit test",
		INSTANCE_NUMBER_FIELD: 42,
		RAM_FIELD:             1,
		CPU_FIELD:             1,
		ENTERPRISE_FIELD:      "unit test enterprise",
		TEMPLATE_FIELD:        "template1",
		OS_FIELD:              "Debian",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk template1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk slug",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan1",
				CONNECTED_FIELD: true,
			},
			map[string]interface{}{VLAN_NAME_FIELD: "unit test vlan2",
				CONNECTED_FIELD: true,
			},
		},
		DYNAMIC_FIELD: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":null}",
	}
	VM_CREATION_FROM_TEMPLATE_SCHEMA_PRE_CREATION_WRONG_NICS_INIT_MAP = map[string]interface{}{
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		ENTERPRISE_FIELD: "unit test enterprise",
		NAME_FIELD:       "template1",
		OS_FIELD:         "Debian",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk template1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk slug",
			},
		},
		NICS_FIELD:    "Wrong nics type",
		DYNAMIC_FIELD: "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":null}",
	}
	VDC_CREATION_MAP = map[string]interface{}{
		NAME_FIELD:       "Unit test vdc resource",
		ENTERPRISE_FIELD: "enterprise",
		VDC_RESOURCE_FIELD: []interface{}{
			map[string]interface{}{
				RESOURCE_FIELD: RAM_FIELD,
				TOTAL_FIELD:    20,
			},
			map[string]interface{}{
				RESOURCE_FIELD: CPU_FIELD,
				TOTAL_FIELD:    1,
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_enterprise",
				TOTAL_FIELD:    10,
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_performance",
				TOTAL_FIELD:    10,
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_high_performance",
				TOTAL_FIELD:    10,
			},
		},
	}
	VDC_READ_RESPONSE_MAP = map[string]interface{}{
		NAME_FIELD:       "Unit test vdc",
		ENTERPRISE_FIELD: "unit test enterprise",
		VDC_RESOURCE_FIELD: []interface{}{
			map[string]interface{}{
				RESOURCE_FIELD: RAM_FIELD,
				USED_FIELD:     "0",
				TOTAL_FIELD:    "20",
				SLUG_FIELD:     "unit test enterprise-dc1-vdc_te-ram",
			},
			map[string]interface{}{
				RESOURCE_FIELD: CPU_FIELD,
				USED_FIELD:     "0",
				TOTAL_FIELD:    "1",
				SLUG_FIELD:     "unit test enterprise-dc1-vdc_te-cpu",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_enterprise",
				USED_FIELD:     "0",
				TOTAL_FIELD:    "10",
				SLUG_FIELD:     "unit test enterprise-dc1-vdc_te-storage_enterprise",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_performance",
				USED_FIELD:     "0",
				TOTAL_FIELD:    "10",
				SLUG_FIELD:     "unit test enterprise-dc1-vdc_te-storage_performance",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_high_performance",
				USED_FIELD:     "0",
				TOTAL_FIELD:    "10",
				SLUG_FIELD:     "unit test enterprise-dc1-vdc_te-storage_high_performance",
			},
		},
		SLUG_FIELD:    "unit test enterprise-dc1-vdc_te",
		DYNAMIC_FIELD: "",
	}
	NO_TEMPLATE_VM_MAP = map[string]interface{}{
		NAME_FIELD:  "Unit test no template vm resource",
		STATE_FIELD: "UP",
		OS_FIELD:    "Debian",
		RAM_FIELD:   8,
		CPU_FIELD:   4,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 1 update",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 2",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  true,
			},
		},
		VDC_FIELD:           VDC_FIELD,
		BOOT_FIELD:          "on disk",
		STORAGE_CLASS_FIELD: "storage_enterprise",
		SLUG_FIELD:          "42",
		TOKEN_FIELD:         "424242",
		BACKUP_FIELD:        "backup_no_backup",
		DISK_IMAGE_FIELD:    "",
		PLATFORM_NAME_FIELD: "42",
		BACKUP_SIZE_FIELD:   42,
		COMMENT_FIELD:       "42",
	}
	EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP = map[string]interface{}{
		NAME_FIELD:          "Unit test template no disc add on vm resource",
		ENTERPRISE_FIELD:    "unit test enterprise",
		TEMPLATE_FIELD:      "template1",
		STATE_FIELD:         "UP",
		VDC_FIELD:           VDC_FIELD,
		BOOT_FIELD:          "on disk",
		STORAGE_CLASS_FIELD: "storage_enterprise",
		SLUG_FIELD:          "42",
		TOKEN_FIELD:         "424242",
		BACKUP_FIELD:        "backup_no_backup",
	}
	INSTANCE_NUMBER_FIELD_UNIT_TEST_VM_INSTANCE = map[string]interface{}{
		NAME_FIELD:            "INSTANCE_NUMBER_FIELDUnitTest",
		INSTANCE_NUMBER_FIELD: 42,
		ENTERPRISE_FIELD:      "unit test enterprise",
		TEMPLATE_FIELD:        "template1",
		STATE_FIELD:           "UP",
		VDC_FIELD:             VDC_FIELD,
		BOOT_FIELD:            "on disk",
		STORAGE_CLASS_FIELD:   "storage_enterprise",
		SLUG_FIELD:            "42",
		TOKEN_FIELD:           "424242",
		BACKUP_FIELD:          "backup_no_backup",
	}
	EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP = map[string]interface{}{
		NAME_FIELD:       "EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP",
		ENTERPRISE_FIELD: "unit test enterprise",
		TEMPLATE_FIELD:   "template1",
		STATE_FIELD:      "UP",
		OS_FIELD:         "Debian",
		RAM_FIELD:        8,
		CPU_FIELD:        4,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
			},
		},
		NICS_FIELD: []interface{}{
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
		VDC_FIELD:           VDC_FIELD,
		BOOT_FIELD:          "on disk",
		STORAGE_CLASS_FIELD: "storage_enterprise",
		SLUG_FIELD:          "42",
		TOKEN_FIELD:         "424242",
		BACKUP_FIELD:        "backup_no_backup",
		DISK_IMAGE_FIELD:    "",
		PLATFORM_NAME_FIELD: "42",
		BACKUP_SIZE_FIELD:   42,
		COMMENT_FIELD:       "42",
	}
	EXISTING_TEMPLATE_WITH_DELETED_DISK_VM_MAP = map[string]interface{}{
		ID_FIELD:         "EXISTING_TEMPLATE_AND_VM_INSTANCE_WITH_DELETED_DISK_VM_MAP",
		NAME_FIELD:       "EXISTING_TEMPLATE_WITH_DELETED_DISK_VM_MAP",
		ENTERPRISE_FIELD: "unit test enterprise",
		TEMPLATE_FIELD:   "template1",
		STATE_FIELD:      "UP",
		OS_FIELD:         "Debian",
		RAM_FIELD:        8,
		CPU_FIELD:        4,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
			},
		},
		NICS_FIELD:          []interface{}{},
		VDC_FIELD:           VDC_FIELD,
		BOOT_FIELD:          "on disk",
		STORAGE_CLASS_FIELD: "storage_enterprise",
		SLUG_FIELD:          "42",
		TOKEN_FIELD:         "424242",
		BACKUP_FIELD:        "backup_no_backup",
		DISK_IMAGE_FIELD:    "",
		PLATFORM_NAME_FIELD: "42",
		BACKUP_SIZE_FIELD:   42,
		COMMENT_FIELD:       "42",
		DYNAMIC_FIELD:       "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":null}",
	}
	NON_EXISTING_TEMPLATE_VM_MAP = map[string]interface{}{
		NAME_FIELD:          "windows95 vm",
		ENTERPRISE_FIELD:    "unit test enterprise",
		TEMPLATE_FIELD:      "windows95",
		STATE_FIELD:         "UP",
		RAM_FIELD:           8,
		CPU_FIELD:           4,
		VDC_FIELD:           VDC_FIELD,
		BOOT_FIELD:          "on disk",
		STORAGE_CLASS_FIELD: "storage_enterprise",
		SLUG_FIELD:          "42",
		TOKEN_FIELD:         "424242",
		BACKUP_FIELD:        "backup_no_backup",
		DISK_IMAGE_FIELD:    "",
	}
	TEMPLATE2_MAP = map[string]interface{}{
		ID_FIELD:         40,
		NAME_FIELD:       "template2",
		SLUG_FIELD:       "unit test disk goulouglougoulouglou",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "CentOS",
		ENTERPRISE_FIELD: "unit test enterprise",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk goulouglouglou",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk goulouglou slug",
			},
		},
		NICS_FIELD:    []interface{}{},
		"login":       "",
		"password":    "",
		DYNAMIC_FIELD: "",
	}
	LAST_TEMPLATE_IN_LIST = map[string]interface{}{
		ID_FIELD:         69,
		NAME_FIELD:       "lastTemplate",
		SLUG_FIELD:       "lastTemplate-slug",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "Debian",
		ENTERPRISE_FIELD: "unit test enterprise",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk-debian9-rd-1",
				SIZE_FIELD:          10,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "disk-debian9-rd-1",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan1",
				MAC_ADRESS_FIELD: "00:50:56:00:01:de",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan2",
				MAC_ADRESS_FIELD: "00:50:56:00:01:df",
				CONNECTED_FIELD:  true,
			},
		},
		"login":       nil,
		"password":    nil,
		DYNAMIC_FIELD: nil,
	}
	TEMPLATE1_MAP = map[string]interface{}{
		ID_FIELD:         82,
		NAME_FIELD:       "template1",
		SLUG_FIELD:       "template1 slug",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "CentOS",
		BOOT_FIELD:       "on disk",
		ENTERPRISE_FIELD: "unit test enterprise",
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk template1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk slug",
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
	}
	TEMPLATE1_MAP_BIS = map[string]interface{}{
		ID_FIELD:         82,
		NAME_FIELD:       "template1",
		SLUG_FIELD:       "template1 slug",
		RAM_FIELD:        1,
		CPU_FIELD:        1,
		OS_FIELD:         "CentOS",
		BOOT_FIELD:       "on disk",
		ENTERPRISE_FIELD: "unit test enterprise",
		//DISKS_FIELD:      []interface{}{},
		//NICS_FIELD:       []interface{}{},
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "unit test disk template1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "unit test disk slug",
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
		DYNAMIC_FIELD: "",
	}
	TEMPLATES_LIST = []interface{}{
		TEMPLATE2_MAP,
		TEMPLATE1_MAP,
		map[string]interface{}{
			ID_FIELD:         41,
			NAME_FIELD:       "template3",
			SLUG_FIELD:       "unit test template3 slug",
			RAM_FIELD:        1,
			CPU_FIELD:        1,
			OS_FIELD:         "Debian",
			ENTERPRISE_FIELD: "unit test enterprise",
			DISKS_FIELD: []interface{}{
				map[string]interface{}{
					NAME_FIELD:          "unit test disk2",
					SIZE_FIELD:          20,
					STORAGE_CLASS_FIELD: "storage_enterprise",
					SLUG_FIELD:          "unit test disk slug 2",
				},
			},
			NICS_FIELD:    []interface{}{},
			"login":       "",
			"password":    "",
			DYNAMIC_FIELD: "",
		},
		map[string]interface{}{
			ID_FIELD:         43,
			NAME_FIELD:       "template4",
			SLUG_FIELD:       "tpl-centos7-rd",
			RAM_FIELD:        1,
			CPU_FIELD:        1,
			OS_FIELD:         "CentOS",
			ENTERPRISE_FIELD: "unit test enterprise",
			DISKS_FIELD: []interface{}{
				map[string]interface{}{
					NAME_FIELD:          "unit test disk 1",
					SIZE_FIELD:          20,
					STORAGE_CLASS_FIELD: "storage_enterprise",
					SLUG_FIELD:          "unit test disk slug",
				},
			},
			NICS_FIELD: []interface{}{
				map[string]interface{}{
					VLAN_NAME_FIELD:  "unit test vlan1",
					MAC_ADRESS_FIELD: "00:50:56:00:00:23",
					CONNECTED_FIELD:  true,
				},
				map[string]interface{}{
					VLAN_NAME_FIELD:  "unit test vlan2",
					MAC_ADRESS_FIELD: "00:50:56:00:00:24",
					CONNECTED_FIELD:  true,
				},
			},
			"login":       nil,
			"password":    nil,
			DYNAMIC_FIELD: nil,
		},
		map[string]interface{}{
			ID_FIELD:         58,
			NAME_FIELD:       "template windaube7",
			SLUG_FIELD:       "slug windows7",
			RAM_FIELD:        1,
			CPU_FIELD:        1,
			OS_FIELD:         "Windows Serveur 64bits",
			ENTERPRISE_FIELD: "unit test enterprise",
			DISKS_FIELD: []interface{}{
				map[string]interface{}{
					NAME_FIELD:          "disk-Template-Windows",
					SIZE_FIELD:          60,
					STORAGE_CLASS_FIELD: "storage_enterprise",
					SLUG_FIELD:          "disk-template-windows7",
				},
			},
			NICS_FIELD:    []interface{}{},
			"login":       nil,
			"password":    nil,
			DYNAMIC_FIELD: nil,
		},
		LAST_TEMPLATE_IN_LIST,
	}
	WRONG_TEMPLATES_LIST = []interface{}{
		TEMPLATE2_MAP,
		"Wrongly formated template",
		LAST_TEMPLATE_IN_LIST,
	}
	TEST_UPDATE_VM_MAP = map[string]interface{}{
		ID_FIELD:    "unit test vm",
		NAME_FIELD:  "Unit test vm",
		STATE_FIELD: "DOWN",
		OS_FIELD:    "CentOS",
		RAM_FIELD:   16,
		CPU_FIELD:   8,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 1 update",
				MAC_ADRESS_FIELD: "42",
				CONNECTED_FIELD:  false,
			},
		},
		VDC_FIELD:           "vdc update",
		BOOT_FIELD:          "on disk update",
		STORAGE_CLASS_FIELD: "storage_enterprise update",
		SLUG_FIELD:          "42 update",
		TOKEN_FIELD:         "424242 update",
		BACKUP_FIELD:        "backup_no_backup update",
		DISK_IMAGE_FIELD:    " update",
		PLATFORM_NAME_FIELD: "",
		BACKUP_SIZE_FIELD:   42,
		COMMENT_FIELD:       "",
	}
	TEST_UPDATE_VM_MAP_INTID = map[string]interface{}{
		ID_FIELD:    1212,
		NAME_FIELD:  "Unit test vm",
		STATE_FIELD: "DOWN",
		OS_FIELD:    "CentOS",
		RAM_FIELD:   16,
		CPU_FIELD:   8,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 1 update",
				MAC_ADRESS_FIELD: "42",
				CONNECTED_FIELD:  false,
			},
		},
		VDC_FIELD:           "vdc update",
		BOOT_FIELD:          "on disk update",
		STORAGE_CLASS_FIELD: "storage_enterprise update",
		SLUG_FIELD:          "42 update",
		TOKEN_FIELD:         "424242 update",
		BACKUP_FIELD:        "backup_no_backup update",
		DISK_IMAGE_FIELD:    " update",
		PLATFORM_NAME_FIELD: "",
		BACKUP_SIZE_FIELD:   43,
		COMMENT_FIELD:       "",
	}
	TEST_UPDATE_VM_MAP_FLOATID = map[string]interface{}{
		ID_FIELD:    121212.12,
		NAME_FIELD:  "Unit test vm",
		STATE_FIELD: "DOWN",
		OS_FIELD:    "CentOS",
		RAM_FIELD:   16,
		CPU_FIELD:   8,
		DISKS_FIELD: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 1 update",
				MAC_ADRESS_FIELD: "42",
				CONNECTED_FIELD:  false,
			},
		},
		VDC_FIELD:           "vdc update",
		BOOT_FIELD:          "on disk update",
		STORAGE_CLASS_FIELD: "storage_enterprise update",
		SLUG_FIELD:          "42 update",
		TOKEN_FIELD:         "424242 update",
		BACKUP_FIELD:        "backup_no_backup update",
		DISK_IMAGE_FIELD:    " update",
		PLATFORM_NAME_FIELD: "",
		BACKUP_SIZE_FIELD:   42,
		COMMENT_FIELD:       "",
	}
)

func resourceVdcResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			RESOURCE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			USED_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			TOTAL_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			SLUG_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVdc() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NAME_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			ENTERPRISE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			DATACENTER_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			VDC_RESOURCE_FIELD: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceVdcResource(),
			},
			SLUG_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			DYNAMIC_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resource_vm_disk() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NAME_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			SIZE_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			STORAGE_CLASS_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			SLUG_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resource_vm_nic() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			VLAN_NAME_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			MAC_ADRESS_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			CONNECTED_FIELD: &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resource_vm() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			NAME_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			INSTANCE_NUMBER_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			ENTERPRISE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			TEMPLATE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			STATE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			OS_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			RAM_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			CPU_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			DISKS_FIELD: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resource_vm_disk(),
			},
			NICS_FIELD: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resource_vm_nic(),
			},
			VDC_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			BOOT_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			STORAGE_CLASS_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			SLUG_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			TOKEN_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			BACKUP_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			DISK_IMAGE_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			PLATFORM_NAME_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			BACKUP_SIZE_FIELD: &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			COMMENT_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			DYNAMIC_FIELD: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func FakeVdcInstance_VDC_CREATION_MAP() VDC {
	return VDC{
		Name:       "Unit test vdc resource",
		Enterprise: "enterprise",
		Vdc_resources: []interface{}{
			map[string]interface{}{
				RESOURCE_FIELD: "enterprise-mono-ram",
				USED_FIELD:     0,
				TOTAL_FIELD:    20,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "enterprise-mono-cpu",
				USED_FIELD:     0,
				TOTAL_FIELD:    1,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "enterprise-mono-storage_enterprise",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "enterprise-mono-storage_performance",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "enterprise-mono-storage_high_performance",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
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
				RESOURCE_FIELD: "Resource1",
				USED_FIELD:     1,
				TOTAL_FIELD:    2,
				SLUG_FIELD:     "Unit Test value1",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "Resource2",
				USED_FIELD:     1,
				TOTAL_FIELD:    2,
				SLUG_FIELD:     "Unit Test value2",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "Resource3",
				USED_FIELD:     1,
				TOTAL_FIELD:    2,
				SLUG_FIELD:     "Unit Test value3",
			},
		},
		Slug:         "Unit Test value",
		DynamicField: "Unit Test value",
	}
}

func vmInstanceNO_TEMPLATE_VM_MAP() VM {
	return VM{
		Name:  "Unit test no template vm resource",
		State: "UP",
		OS:    "Debian",
		RAM:   8,
		CPU:   4,
		Disks: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "",
			},
		},
		Nics: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 1 update",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan 2",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  true,
			},
		},
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"\",\"Template_disks_on_creation\":null}",
	}
}

func FakeVmInstance_EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP() VM {
	return VM{
		Name:       "Unit test template no disc add on vm resource-0",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		RAM:        1,
		CPU:        1,
		Disks: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
		Nics:          []interface{}{},
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstance_INSTANCE_NUMBER_FIELD_UNIT_TEST_VM_INSTANCE_MAP() VM {
	return VM{
		Name:       "INSTANCE_NUMBER_FIELDUnitTest-42",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		RAM:        1,
		CPU:        1,
		Disks: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          20,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
		Nics:          []interface{}{},
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstance_EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP() VM {
	return VM{
		Name:       "EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP-0",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		OS:         "Debian",
		RAM:        8,
		CPU:        4,
		Disks: []interface{}{
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
		Nics: []interface{}{
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
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func FakeVmInstance_EXISTING_TEMPLATE_WITH_DELETED_DISK_VM_MAP() VM {
	return VM{
		Name:       "EXISTING_TEMPLATE_WITH_DELETED_DISK_VM_MAP",
		Enterprise: "unit test enterprise",
		Template:   "",
		State:      "UP",
		OS:         "Debian",
		RAM:        8,
		CPU:        4,
		Disks: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "disk 1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "",
			},
		},
		Nics:          []interface{}{},
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Disk_image:    "",
		Platform_name: "42",
		Backup_size:   42,
		Comment:       "",
		DynamicField:  "{\"terraform_provisioned\":true,\"creationTemplate\":\"template1\",\"Template_disks_on_creation\":null}",
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
				NAME_FIELD:          "name1",
				SIZE_FIELD:          10,
				STORAGE_CLASS_FIELD: "Unit Test value",
			},
			map[string]interface{}{
				NAME_FIELD:          "name2",
				SIZE_FIELD:          10,
				STORAGE_CLASS_FIELD: "Unit Test value",
			},
		},
		Nics: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan1",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "vlan1",
				MAC_ADRESS_FIELD: "00:21:21:21:21:21",
				CONNECTED_FIELD:  false,
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
	d := resource_vm().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.UpdateLocalResourceState(vm, d, &schemaTooler)

	return d
}

func resource(resourceType string) *schema.Resource {

	resource := &schema.Resource{}
	switch resourceType {
	case VDC_FIELD:
		resource = resourceVdc()
	case VM_RESOURCE_TYPE:
		resource = resource_vm()
	default:
		//return a false resource
		resource = &schema.Resource{
			Schema: map[string]*schema.Schema{
				NAME_FIELD: &schema.Schema{
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
	fakeJson, _ := json.Marshal(TEMPLATES_LIST)
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(fakeJson))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonFake)

	return jsonFake.([]interface{})
}
