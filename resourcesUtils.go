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

type ResourceTooler struct {
	Resource Resourceer
}
type Resourceer interface {
	GetResourceCreationUrl(api *API,
		resourceType string) string
	GetResourceUrl(api *API,
		resourceType string,
		id string) string
	ValidateResourceType(resourceType string) error
	ValidateStatus(api *API,
		resourceType string,
		client ClientTooler) error
	ResourceInstanceCreate(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		resourceType string,
		api *API) (interface{}, error)
}
type ResourceResourceer struct{}

type DynamicFieldStruct struct {
	TerraformProvisioned    bool          `json:"terraform_provisioned"`
	CreationTemplate        string        `json:"creationTemplate"`
	TemplateDisksOnCreation []interface{} `json:"TemplateDisksOnCreation"`
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

func vdcInstanceCreate(d *schema.ResourceData) (VDC, error) {
	var (
		resourceName strings.Builder
	)
	vdc := VDC{
		Name:          d.Get(NameField).(string),
		Enterprise:    d.Get(EnterpriseField).(string),
		Datacenter:    d.Get(DatacenterField).(string),
		Vdc_resources: d.Get(VdcResourceField).([]interface{}),
		Slug:          d.Get(SlugField).(string),
		DynamicField:  d.Get(DynamicField).(string),
	}
	for index, resource := range vdc.Vdc_resources {
		resourceName.Reset()
		resourceName.WriteString(vdc.Enterprise)
		resourceName.WriteString(monoField)
		resourceName.WriteString(resource.(map[string]interface{})[ResourceField].(string))
		resource.(map[string]interface{})[ResourceField] = resourceName.String()
		vdc.Vdc_resources[index] = resource
	}
	return vdc, nil
}

func getTemplateAndUpdateSchema(templateName string,
	d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	api *API) (map[string]interface{}, error) {
	templateList, err1 := clientTooler.Client.GetTemplatesList(clientTooler,
		d.Get(EnterpriseField).(string),
		api)
	if err1 != nil {
		return map[string]interface{}{}, err1
	}
	template, err2 := templatesTooler.TemplatesTools.FetchTemplateFromList(templateName,
		templateList)
	if err2 != nil {
		return map[string]interface{}{}, err2
	}
	err3 := templatesTooler.TemplatesTools.ValidateTemplate(template)
	if err3 != nil {
		return map[string]interface{}{}, err3
	}
	err4 := templatesTooler.TemplatesTools.UpdateSchemaFromTemplateOnResourceCreation(d,
		template)
	if err4 != nil {
		return map[string]interface{}{}, err4
	}
	return template, nil
}

func vmInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	api *API) (VM, error) {
	var (
		templateError error                  = nil
		template      map[string]interface{} = nil
		templateName  string                 = d.Get(TemplateField).(string)
		vmName        strings.Builder
	)
	vmName.WriteString(d.Get(NameField).(string))
	if templateName != "" && d.Id() == "" {
		instanceNumber := d.Get(InstanceNumberField).(int)
		vmName.WriteString(resourceNameCountSeparator)
		vmName.WriteString(strconv.Itoa(instanceNumber))
		template,
			templateError = getTemplateAndUpdateSchema(templateName,
			d,
			clientTooler,
			templatesTooler,
			api)
	}
	if templateError != nil {
		return VM{}, templateError
	}
	vm := VM{
		Name:          vmName.String(),
		Enterprise:    d.Get(EnterpriseField).(string),
		State:         d.Get(StateField).(string),
		OS:            d.Get(OsField).(string),
		RAM:           d.Get(RamField).(int),
		CPU:           d.Get(CpuField).(int),
		Disks:         d.Get(DisksField).([]interface{}),
		Nics:          d.Get(NicsField).([]interface{}),
		Vdc:           d.Get(VdcField).(string),
		Boot:          d.Get(BootField).(string),
		Storage_class: d.Get(StorageClassField).(string),
		Slug:          d.Get(SlugField).(string),
		Token:         d.Get(TokenField).(string),
		Backup:        d.Get(BackupField).(string),
		Disk_image:    d.Get(DiskImageField).(string),
		Platform_name: d.Get(PlatformNameField).(string),
		Backup_size:   d.Get(BackupSizeField).(int),
		DynamicField:  d.Get(DynamicField).(string),
	}
	if d.Id() == "" {
		dynamicFieldStruct := DynamicFieldStruct{
			TerraformProvisioned:    true,
			CreationTemplate:        d.Get(TemplateField).(string),
			TemplateDisksOnCreation: nil,
		}
		if template != nil {
			dynamicFieldStruct.TemplateDisksOnCreation = template[DisksField].([]interface{})
			_, err := templatesTooler.TemplatesTools.CreateTemplateOverrideConfig(d,
				template)
			if err != nil {
				return VM{}, err
			}
			vm.Template = d.Get(TemplateField).(string)
		}
		dynamicFieldJson, err2 := json.Marshal(dynamicFieldStruct)
		if err2 != nil {
			return VM{}, err2
		}
		vm.DynamicField = string(dynamicFieldJson)
	}
	return vm, nil
}

func (resource ResourceResourceer) ResourceInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceType string,
	api *API) (interface{}, error) {
	switch resourceType {
	case VdcResourceType:
		return vdcInstanceCreate(d)
	case VmResourceType:
		return vmInstanceCreate(d,
			clientTooler,
			templatesTooler,
			api)
	default:
		return nil, resource.ValidateResourceType(resourceType)
	}
}

func (resource ResourceResourceer) ValidateResourceType(resourceType string) error {
	switch resourceType {
	case VdcResourceType:
		return nil
	case VmResourceType:
		return nil
	default:
		return errWrongResourceTypeBuilder(resourceType)
	}
}

func (resource ResourceResourceer) GetResourceCreationUrl(api *API,
	resourceType string) string {
	var resourceUrl strings.Builder
	resourceUrl.WriteString(api.URL)
	resourceUrl.WriteString(resourceType)
	resourceUrl.WriteString("/")
	return resourceUrl.String()
}

func (resource ResourceResourceer) GetResourceUrl(api *API,
	resourceType string,
	resourceId string) string {
	var resourceUrl strings.Builder
	s_create_url := resource.GetResourceCreationUrl(api, resourceType)
	resourceUrl.WriteString(s_create_url)
	resourceUrl.WriteString(resourceId)
	resourceUrl.WriteString("/")
	return resourceUrl.String()
}

func (resource ResourceResourceer) ValidateStatus(api *API,
	resourceType string,
	clientTooler ClientTooler) error {
	var apiError error
	var responseBody string
	req, _ := http.NewRequest("GET",
		resource.GetResourceCreationUrl(api, resourceType),
		nil)
	req.Header.Add(httpAuthorization, httpTokenHeader+api.Token)
	resp, apiError := clientTooler.Client.Do(api, req)
	if apiError == nil {
		if resp.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			responseBody = string(bodyBytes)
			switch {
			case resp.StatusCode == http.StatusUnauthorized:
				apiError = errors.New(resp.Status + responseBody)
			case resp.Header.Get(httpReqContentType) != httpJsonContentType:
				apiError = errors.New("Could not get a proper json response from \"" +
					api.URL + errApiDownOrWrongApiUrl)
			}
		} else {
			apiError = errors.New("Could not get a response body from \"" + api.URL +
				errApiDownOrWrongApiUrl)
		}
	} else {
		apiError = errors.New("Could not get a response from \"" + api.URL +
			errApiDownOrWrongApiUrl)
	}

	return apiError
}
