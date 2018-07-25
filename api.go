package sewan_go_sdk

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

const (
	DEFAULT_RESOURCE_TYPE = "vm"
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
		schemaTools *SchemaTooler,
		resourceType string,
		api *API) (error,
		interface{})
	CreateResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		schemaTools *SchemaTooler,
		resourceType string,
		sewan *API) (error,
		map[string]interface{})
	ReadResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		schemaTools *SchemaTooler,
		resourceType string,
		sewan *API) (error,
		map[string]interface{},
		bool)
	UpdateResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		schemaTools *SchemaTooler,
		resourceType string,
		sewan *API) error
	DeleteResource(d *schema.ResourceData,
		clientTooler *ClientTooler,
		templatesTooler *TemplatesTooler,
		schemaTools *SchemaTooler,
		resourceType string,
		sewan *API) error
}
type AirDrumResources_Apier struct{}

func (apiTools *APITooler) New(token string, url string) *API {
	return &API{
		Token:  token,
		URL:    url,
		Client: &http.Client{},
	}
}

func (apiTools *APITooler) CheckStatus(api *API) error {
	var apiClientErr error
	clientTooler := ClientTooler{Client: HttpClienter{}}
	apiClientErr = apiTools.Api.ValidateStatus(api, DEFAULT_RESOURCE_TYPE, clientTooler)
	return apiClientErr
}
