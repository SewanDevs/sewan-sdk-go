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
	VM_RESOURCE_TYPE        = "vm"
	VDC_RESOURCE_TYPE       = VDC_FIELD
	WRONG_RESOURCE_TYPE     = "a_non_supported_resource_type"
	ENTERPRISE_SLUG         = "unit test enterprise"
)

var (
	VDC_CREATION_MAP = map[string]interface{}{
		NAME_FIELD:       "Unit test vdc resource",
		ENTERPRISE_FIELD: "unit test enterprise",
		DATACENTER_FIELD: "dc1",
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
		DATACENTER_FIELD: "dc1",
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
	EXISTING_TEMPLATE_WITH_MODIFIED_NIC_AND_DISK_VM_MAP = map[string]interface{}{
		NAME_FIELD:       "EXISTING_TEMPLATE_WITH_MODIFIED_NIC_AND_DISK_VM_MAP",
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
				SLUG_FIELD:          "template1 disk1 slug",
			},
		},
		NICS_FIELD: []interface{}{
			map[string]interface{}{
				MAC_ADRESS_FIELD: "00:21:21:21:21:22",
				CONNECTED_FIELD:  true,
				VLAN_NAME_FIELD:  "unit test vlan1",
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan2",
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
				DELETION_FIELD:      true,
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
		DYNAMIC_FIELD:       "{\"terraform_provisioned\":true,\"creation_template\":\"template1\",\"disks_created_from_template\":null}",
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
	TEMPLATES_LIST = []interface{}{
		map[string]interface{}{
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
			DATACENTER_FIELD: "dc2",
			NICS_FIELD:       []interface{}{},
			"login":          "",
			"password":       "",
			DYNAMIC_FIELD:    "",
		},
		map[string]interface{}{
			ID_FIELD:         82,
			NAME_FIELD:       "template1",
			SLUG_FIELD:       "template1 slug",
			RAM_FIELD:        1,
			CPU_FIELD:        1,
			OS_FIELD:         "CentOS",
			ENTERPRISE_FIELD: "unit test enterprise",
			DISKS_FIELD: []interface{}{
				map[string]interface{}{
					NAME_FIELD:          "unit test disk template1",
					SIZE_FIELD:          20,
					STORAGE_CLASS_FIELD: "storage_enterprise",
					SLUG_FIELD:          "unit test disk slug",
				},
			},
			DATACENTER_FIELD: "dc1",
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
		},
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
			DATACENTER_FIELD: "dc2",
			NICS_FIELD:       []interface{}{},
			"login":          "",
			"password":       "",
			DYNAMIC_FIELD:    "",
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
			DATACENTER_FIELD: "dc2",
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
			DATACENTER_FIELD: "dc2",
			NICS_FIELD:       []interface{}{},
			"login":          nil,
			"password":       nil,
			DYNAMIC_FIELD:    nil,
		},
		map[string]interface{}{
			ID_FIELD:         69,
			NAME_FIELD:       "template5",
			SLUG_FIELD:       "template5-slug",
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
			DATACENTER_FIELD: "dc2",
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
		},
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
				DELETION_FIELD:      false,
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",

				DELETION_FIELD: false,
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

				DELETION_FIELD: false,
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",

				DELETION_FIELD: false,
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

				DELETION_FIELD: false,
			},
			map[string]interface{}{
				NAME_FIELD:          "disk 2 update",
				SIZE_FIELD:          42,
				STORAGE_CLASS_FIELD: "STORAGE_CLASS_FIELD update",
				SLUG_FIELD:          "slug update",

				DELETION_FIELD: false,
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

func resource_vdc_resource() *schema.Resource {
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

func resource_vdc() *schema.Resource {
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
				Elem:     resource_vdc_resource(),
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
			DELETION_FIELD: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
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

func Fake_vdcInstance_VDC_CREATION_MAP() VDC {
	return VDC{
		Name:       "Unit test vdc resource",
		Enterprise: "unit test enterprise",
		Datacenter: "dc1",
		Vdc_resources: []interface{}{
			map[string]interface{}{
				RESOURCE_FIELD: RAM_FIELD,
				USED_FIELD:     0,
				TOTAL_FIELD:    20,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: CPU_FIELD,
				USED_FIELD:     0,
				TOTAL_FIELD:    1,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_enterprise",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_performance",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
			},
			map[string]interface{}{
				RESOURCE_FIELD: "storage_high_performance",
				USED_FIELD:     0,
				TOTAL_FIELD:    10,
				SLUG_FIELD:     "",
			},
		},
		Slug:          "",
		Dynamic_field: "",
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
		Slug:          "Unit Test value",
		Dynamic_field: "Unit Test value",
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

				DELETION_FIELD: false,
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
		Dynamic_field: "{\"terraform_provisioned\":true,\"creation_template\":\"\",\"disks_created_from_template\":null}",
	}
}

func Fake_vmInstance_EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP() VM {
	return VM{
		Name:       "Unit test template no disc add on vm resource",
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
				DELETION_FIELD:      false,
			},
		},
		Nics:          []interface{}{},
		Vdc:           VDC_FIELD,
		Boot:          "on disk",
		Storage_class: "storage_enterprise",
		Slug:          "42",
		Token:         "424242",
		Backup:        "backup_no_backup",
		Dynamic_field: "{\"terraform_provisioned\":true,\"creation_template\":\"template1\",\"disks_created_from_template\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func Fake_vmInstance_EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP() VM {
	return VM{
		Name:       "EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP",
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

				DELETION_FIELD: false,
			},
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          25,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",

				DELETION_FIELD: false,
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
		Dynamic_field: "{\"terraform_provisioned\":true,\"creation_template\":\"template1\",\"disks_created_from_template\":[{\"name\":\"template1 disk1\",\"size\":20,\"slug\":\"template1 disk1 slug\",\"storage_class\":\"storage_enterprise\"}]}",
	}
}

func Fake_vmInstance_EXISTING_TEMPLATE_WITH_MODIFIED_NIC_AND_DISK_VM_MAP() VM {
	return VM{
		Name:       "EXISTING_TEMPLATE_WITH_MODIFIED_NIC_AND_DISK_VM_MAP",
		Enterprise: "unit test enterprise",
		Template:   "template1",
		State:      "UP",
		OS:         "Debian",
		RAM:        8,
		CPU:        4,
		Disks: []interface{}{
			map[string]interface{}{
				NAME_FIELD:          "template1 disk1",
				SIZE_FIELD:          24,
				STORAGE_CLASS_FIELD: "storage_enterprise",
				SLUG_FIELD:          "template1 disk1 slug",

				DELETION_FIELD: false,
			},
		},
		Nics: []interface{}{
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan1",
				MAC_ADRESS_FIELD: "00:21:21:21:21:22",
				CONNECTED_FIELD:  true,
			},
			map[string]interface{}{
				VLAN_NAME_FIELD:  "unit test vlan2",
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
		Dynamic_field: "{\"terraform_provisioned\":true,\"creation_template\":\"template1\",\"disks_created_from_template\":null}",
	}
}

func Fake_vmInstance_EXISTING_TEMPLATE_WITH_DELETED_DISK_VM_MAP() VM {
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

				DELETION_FIELD: false,
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
		Dynamic_field: "{\"terraform_provisioned\":true,\"creation_template\":\"template1\",\"disks_created_from_template\":null}",
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

func vdc_schema_init(vdc map[string]interface{}) *schema.ResourceData {
	d := resource_vdc().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.Update_local_resource_state(vdc, d, &schemaTooler)

	return d
}

func vm_schema_init(vm map[string]interface{}) *schema.ResourceData {
	d := resource_vm().TestResourceData()

	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.Update_local_resource_state(vm, d, &schemaTooler)

	return d
}

func resource(resourceType string) *schema.Resource {

	resource := &schema.Resource{}
	switch resourceType {
	case VDC_FIELD:
		resource = resource_vdc()
	case "vm":
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
	simple_json, _ := json.Marshal(Resp_Body{Detail: "a simple json"})
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(simple_json))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonStub)

	return jsonStub.(map[string]interface{})
}

func JsonTemplateListFake() []interface{} {
	var jsonFake interface{}
	fake_json, _ := json.Marshal(TEMPLATES_LIST)
	jsonBytes := ioutil.NopCloser(bytes.NewBuffer(fake_json))
	readBytes, _ := ioutil.ReadAll(jsonBytes)
	_ = json.Unmarshal(readBytes, &jsonFake)

	return jsonFake.([]interface{})
}
