package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"strings"
)

const (
	defaultResourceType = VmResourceType
)

type API struct {
	Token  string
	URL    string
	Client *http.Client
}
type APITooler struct {
	Api APIer
}
type APIer interface {
	CreateResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		resourceTooler *ResourceTooler,
		resourceType string,
		sewan *API) (map[string]interface{}, error)
	ReadResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		resourceTooler *ResourceTooler,
		resourceType string,
		sewan *API) (map[string]interface{}, error)
	UpdateResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		resourceTooler *ResourceTooler,
		resourceType string,
		sewan *API) error
	DeleteResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		resourceTooler *ResourceTooler,
		resourceType string,
		sewan *API) error
}
type AirDrumResourcesApier struct{}

func (apiTools *APITooler) New(token string, url string) *API {
	return &API{
		Token:  token,
		URL:    url,
		Client: &http.Client{},
	}
}

func (apiTools *APITooler) CheckCloudDcApiStatus(api *API,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler) error {
	var apiClientErr error
	apiClientErr = resourceTooler.Resource.ValidateStatus(api,
		defaultResourceType,
		*clientTooler)
	return apiClientErr
}

func (apier AirDrumResourcesApier) CreateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	var (
		instanceName string = d.Get(NameField).(string)
	)
	resourceInstance, err1 := resourceTooler.Resource.ResourceInstanceCreate(d,
		clientTooler,
		templatesTooler,
		resourceType,
		sewan)
	if err1 != nil {
		return map[string]interface{}{}, err1.(error)
	}
	resourceJson, err2 := json.Marshal(resourceInstance)
	if err2 != nil {
		return map[string]interface{}{}, err2
	}
	req, err3 := http.NewRequest("POST",
		resourceTooler.Resource.GetResourceCreationUrl(sewan, resourceType),
		bytes.NewBuffer(resourceJson))
	if err3 != nil {
		return map[string]interface{}{}, err3
	}
	req.Header.Add(httpAuthorization, httpTokenHeader+sewan.Token)
	req.Header.Add(httpReqContentType, httpJsonContentType)
	resp, err4 := clientTooler.Client.Do(sewan, req)
	switch {
	case err4 != nil:
		return map[string]interface{}{}, errDoCrudRequestsBuilder(creationOperation,
			instanceName, err4)
	default:
		createdResource, err5 := clientTooler.Client.HandleResponse(resp,
			http.StatusCreated,
			httpJsonContentType)
		if createdResource == nil {
			return map[string]interface{}{}, err5
		}
		return createdResource.(map[string]interface{}), err5
	}
}

func (apier AirDrumResourcesApier) ReadResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	err1 := resourceTooler.Resource.ValidateResourceType(resourceType)
	if err1 != nil {
		return map[string]interface{}{}, err1
	}
	req, err2 := http.NewRequest("GET",
		resourceTooler.Resource.GetResourceUrl(sewan, resourceType, d.Id()), nil)
	if err2 != nil {
		return map[string]interface{}{}, err2
	}
	req.Header.Add(httpAuthorization, httpTokenHeader+sewan.Token)
	resp, err3 := clientTooler.Client.Do(sewan, req)
	switch {
	case err3 != nil:
		return map[string]interface{}{}, errDoCrudRequestsBuilder(readOperation,
			d.Get(NameField).(string),
			err3)
	default:
		if (resp != nil) && (resp.StatusCode == http.StatusNotFound) {
			return map[string]interface{}{}, ErrResourceNotExist
		}
		readResource, err4 := clientTooler.Client.HandleResponse(resp,
			http.StatusOK,
			httpJsonContentType)
		if readResource == nil {
			return map[string]interface{}{}, err4
		}
		if resourceType == VdcResourceType {
			err5 := updateSchemaReadVdcResource(d,
				readResource.(map[string]interface{}))
			if err5 != nil {
				return map[string]interface{}{}, err5
			}
		}
		return readResource.(map[string]interface{}), err4
	}
}

func updateSchemaReadVdcResource(d *schema.ResourceData,
	readResource map[string]interface{}) error {
	var (
		resourceNamePrefix strings.Builder
		resourcesList      []interface{}
	)
	resourceNamePrefix.WriteString(readResource[EnterpriseField].(string))
	resourceNamePrefix.WriteString(monoField)
	for _, resource := range readResource[VdcResourceField].([]interface{}) {
		resource.(map[string]interface{})[ResourceField] = strings.TrimPrefix(resource.(map[string]interface{})[ResourceField].(string), resourceNamePrefix.String())
		resourcesList = append(resourcesList, resource)
	}
	return d.Set(VdcResourceField, resourcesList)
}

func (apier AirDrumResourcesApier) UpdateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	resourceInstance,
		err1 := resourceTooler.Resource.ResourceInstanceCreate(d,
		clientTooler,
		templatesTooler,
		resourceType,
		sewan)
	if err1 != nil {
		return err1
	}
	resourceJson, err2 := json.Marshal(resourceInstance)
	if err2 != nil {
		return err2
	}
	req, err3 := http.NewRequest("PUT",
		resourceTooler.Resource.GetResourceUrl(sewan, resourceType, d.Id()),
		bytes.NewBuffer(resourceJson))
	if err3 != nil {
		return err3
	}
	req.Header.Add(httpAuthorization, httpTokenHeader+sewan.Token)
	req.Header.Add(httpReqContentType, httpJsonContentType)
	resp, err4 := clientTooler.Client.Do(sewan, req)
	switch {
	case err4 != nil:
		return errDoCrudRequestsBuilder(updateOperation,
			d.Get(NameField).(string),
			err4)
	default:
		_, err5 := clientTooler.Client.HandleResponse(resp,
			http.StatusOK,
			httpJsonContentType)
		return err5
	}
}

func (apier AirDrumResourcesApier) DeleteResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	err1 := resourceTooler.Resource.ValidateResourceType(resourceType)
	if err1 != nil {
		return err1
	}
	req, err2 := http.NewRequest("DELETE",
		resourceTooler.Resource.GetResourceUrl(sewan, resourceType, d.Id()), nil)
	if err2 != nil {
		return err2
	}
	req.Header.Add(httpAuthorization, httpTokenHeader+sewan.Token)
	resp, err3 := clientTooler.Client.Do(sewan, req)
	switch {
	case err3 != nil:
		return errDoCrudRequestsBuilder(deleteOperation,
			d.Get(NameField).(string),
			err3)
	default:
		_, err4 := clientTooler.Client.HandleResponse(resp,
			http.StatusNoContent,
			httpJsonContentType)
		return err4
	}
}
