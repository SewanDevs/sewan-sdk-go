package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type DynamicField_struct struct {
	Terraform_provisioned      bool          `json:"terraform_provisioned"`
	CreationTemplate           string        `json:"creationTemplate"`
	Template_disks_on_creation []interface{} `json:"Template_disks_on_creation"`
}

type VDC_resource struct {
	Resource string `json:"vdc_resources"`
	Used     int    `json:"used"`
	Total    int    `json:"total"`
	Slug     string `json:"slug"`
}

type VDC struct {
	Name          string        `json:"name"`
	Enterprise    string        `json:"enterprise"`
	Datacenter    string        `json:"datacenter"`
	Vdc_resources []interface{} `json:"vdc_resources"`
	Slug          string        `json:"slug"`
	DynamicField  string        `json:"dynamic_field"`
}

type VM_DISK struct {
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Storage_class string `json:"storage_class"`
	Slug          string `json:"slug"`
	V_disk        string `json:"v_disk,omitempty"`
}

type VM_NIC struct {
	Vlan        string `json:"vlan"`
	Mac_address string `json:"mac_address"`
	Connected   bool   `json:"connected"`
}

type VM struct {
	Name          string        `json:"name"`
	Enterprise    string        `json:"enterprise"`
	Template      string        `json:"template,omitempty"`
	State         string        `json:"state"`
	OS            string        `json:"os,omitempty"`
	RAM           int           `json:"ram"`
	CPU           int           `json:"cpu"`
	Disks         []interface{} `json:"disks,omitempty"`
	Nics          []interface{} `json:"nics,omitempty"`
	Vdc           string        `json:"vdc"`
	Boot          string        `json:"boot"`
	Storage_class string        `json:"storage_class"`
	Slug          string        `json:"slug"`
	Token         string        `json:"token"`
	Backup        string        `json:"backup"`
	Disk_image    string        `json:"disk_image"`
	Platform_name string        `json:"platform_name"`
	Backup_size   int           `json:"backup_size"`
	Comment       string        `json:"comment,omitempty"`
	DynamicField  string        `json:"dynamic_field"`
	Outsourcing   string        `json:"outsourcing"`
}

func vdcInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	api *API) (VDC, error) {

	var (
		vdc          VDC
		resourceName strings.Builder
	)
	vdc = VDC{
		Name:          d.Get(NAME_FIELD).(string),
		Enterprise:    d.Get(ENTERPRISE_FIELD).(string),
		Datacenter:    d.Get(DATACENTER_FIELD).(string),
		Vdc_resources: d.Get(VDC_RESOURCE_FIELD).([]interface{}),
		Slug:          d.Get(SLUG_FIELD).(string),
		DynamicField:  d.Get(DYNAMIC_FIELD).(string),
	}

	for index, resource := range vdc.Vdc_resources {
		resourceName.Reset()
		resourceName.WriteString(vdc.Enterprise)
		resourceName.WriteString(MONO_FIELD)
		resourceName.WriteString(resource.(map[string]interface{})[RESOURCE_FIELD].(string))
		resource.(map[string]interface{})[RESOURCE_FIELD] = resourceName.String()
		vdc.Vdc_resources[index] = resource
	}

	return vdc, nil
}

func vmInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	api *API) (VM, error) {

	var (
		vm                         VM
		getTemplatesListError      error                  = nil
		fetchTemplateFromListError error                  = nil
		templateFormatError        error                  = nil
		instanceCreationError      error                  = nil
		template                   map[string]interface{} = nil
		templateName               string                 = d.Get(TEMPLATE_FIELD).(string)
		enterprise                 string                 = d.Get(ENTERPRISE_FIELD).(string)
		vmName                     strings.Builder
		instanceNumber             int
	)
	logger := LoggerCreate("vminstanceCreate" + d.Id() + ".log")
	vmName.WriteString(d.Get(NAME_FIELD).(string))
	if templateName != "" && d.Id() == "" {
		var templateList []interface{}
		instanceNumber = d.Get(INSTANCE_NUMBER_FIELD).(int)
		vmName.WriteString(RESOURCE_NAME_COUNT_SEPARATOR)
		vmName.WriteString(strconv.Itoa(instanceNumber))
		templateList,
			getTemplatesListError = clientTooler.Client.GetTemplatesList(clientTooler,
			enterprise, api)
		if getTemplatesListError == nil {
			template,
				fetchTemplateFromListError = templatesTooler.TemplatesTools.FetchTemplateFromList(templateName,
				templateList)
			templateFormatError = templatesTooler.TemplatesTools.ValidateTemplate(template)
			switch {
			case fetchTemplateFromListError != nil:
				instanceCreationError = fetchTemplateFromListError
			case templateFormatError != nil:
				instanceCreationError = templateFormatError
			default:
				instanceCreationError = templatesTooler.TemplatesTools.UpdateSchemaFromTemplateOnResourceCreation(d,
					template)
			}
		} else {
			instanceCreationError = getTemplatesListError
		}
	}
	logger.Println("instanceCreationError = ", instanceCreationError)
	if instanceCreationError == nil {
		vm = VM{
			Name:          vmName.String(),
			Enterprise:    d.Get(ENTERPRISE_FIELD).(string),
			State:         d.Get(STATE_FIELD).(string),
			OS:            d.Get(OS_FIELD).(string),
			RAM:           d.Get(RAM_FIELD).(int),
			CPU:           d.Get(CPU_FIELD).(int),
			Disks:         d.Get(DISKS_FIELD).([]interface{}),
			Nics:          d.Get(NICS_FIELD).([]interface{}),
			Vdc:           d.Get(VDC_FIELD).(string),
			Boot:          d.Get(BOOT_FIELD).(string),
			Storage_class: d.Get(STORAGE_CLASS_FIELD).(string),
			Slug:          d.Get(SLUG_FIELD).(string),
			Token:         d.Get(TOKEN_FIELD).(string),
			Backup:        d.Get(BACKUP_FIELD).(string),
			Disk_image:    d.Get(DISK_IMAGE_FIELD).(string),
			Platform_name: d.Get(PLATFORM_NAME_FIELD).(string),
			Backup_size:   d.Get(BACKUP_SIZE_FIELD).(int),
			DynamicField:  d.Get(DYNAMIC_FIELD).(string),
		}
		logger.Println("vm.Name =", vm.Name)
		if d.Id() == "" {
			dynamicFieldStruct := DynamicField_struct{
				Terraform_provisioned:      true,
				CreationTemplate:           d.Get(TEMPLATE_FIELD).(string),
				Template_disks_on_creation: nil,
			}
			if template != nil {
				dynamicFieldStruct.Template_disks_on_creation = template[DISKS_FIELD].([]interface{})
				overrideError, _ := templatesTooler.TemplatesTools.CreateTemplateOverrideConfig(d, template)
				if overrideError != nil {
					instanceCreationError = overrideError
				}
				vm.Template = d.Get(TEMPLATE_FIELD).(string)
			}
			dynamicFieldJson, _ := json.Marshal(dynamicFieldStruct)
			vm.DynamicField = string(dynamicFieldJson)
		}
	}
	logger.Println("vm = ", vm)
	logger.Println("instanceCreationError = ", instanceCreationError)
	return vm, instanceCreationError
}

func (apier AirDrumResources_Apier) ResourceInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	api *API) (error, interface{}) {

	var (
		resourceInstance interface{} = nil
		instanceError    error       = nil
	)

	switch resourceType {
	case VDC_FIELD:
		resourceInstance, instanceError = vdcInstanceCreate(d,
			clientTooler,
			api)
	case VM_RESOURCE_TYPE:
		resourceInstance, instanceError = vmInstanceCreate(d,
			clientTooler,
			templatesTooler,
			schemaTools,
			api)
	default:
		instanceError = apier.ValidateResourceType(resourceType)
	}

	return instanceError, resourceInstance
}

func (apier AirDrumResources_Apier) ValidateResourceType(resourceType string) error {
	var err error

	switch resourceType {
	case VDC_FIELD:
		err = nil
	case VM_RESOURCE_TYPE:
		err = nil
	default:
		err = errors.New("Resource of type \"" + resourceType + "\" not supported," +
			"list of accepted resource types :\n\r" +
			"- \"vdc\"\n\r" +
			"- \"vm\"")
	}

	return err
}

func (apier AirDrumResources_Apier) GetResourceCreationUrl(api *API,
	resourceType string) string {

	var resourceUrl strings.Builder
	resourceUrl.WriteString(api.URL)
	resourceUrl.WriteString(resourceType)
	resourceUrl.WriteString("/")
	return resourceUrl.String()
}

func (apier AirDrumResources_Apier) GetResourceUrl(api *API,
	resourceType string,
	resourceId string) string {

	var resourceUrl strings.Builder
	apiTools := APITooler{
		Api: apier,
	}
	s_create_url := apiTools.Api.GetResourceCreationUrl(api, resourceType)
	resourceUrl.WriteString(s_create_url)
	resourceUrl.WriteString(resourceId)
	resourceUrl.WriteString("/")
	return resourceUrl.String()
}

func (apier AirDrumResources_Apier) ValidateStatus(api *API,
	resourceType string,
	clientTooler ClientTooler) error {

	var apiError error
	var responseBody string
	apiTools := APITooler{
		Api: apier,
	}
	req, _ := http.NewRequest("GET",
		apiTools.Api.GetResourceCreationUrl(api, resourceType),
		nil)
	req.Header.Add("authorization", "Token "+api.Token)
	resp, apiError := clientTooler.Client.Do(api, req)

	if apiError == nil {
		if resp.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			responseBody = string(bodyBytes)
			switch {
			case resp.StatusCode == http.StatusUnauthorized:
				apiError = errors.New(resp.Status + responseBody)
			case resp.Header.Get("content-type") != "application/json":
				apiError = errors.New("Could not get a proper json response from \"" +
					api.URL + "\", the api is down or this url is wrong.")
			}
		} else {
			apiError = errors.New("Could not get a response body from \"" + api.URL +
				"\", the api is down or this url is wrong.")
		}
	} else {
		apiError = errors.New("Could not get a response from \"" + api.URL +
			"\", the api is down or this url is wrong.")
	}

	return apiError
}
