package sewan_go_sdk

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

type FakeAirDrumResource_APIer struct{}

func (apier FakeAirDrumResource_APIer) ValidateStatus(api *API,
	resourceType string,
	client ClientTooler) error {

	var err error
	switch {
	case api.URL != RIGHT_API_URL:
		err = errors.New(WRONG_API_URL_ERROR)
	case api.Token != RIGHT_API_TOKEN:
		err = errors.New(WRONG_TOKEN_ERROR)
	default:
		err = nil
	}
	return err
}

func (apier FakeAirDrumResource_APIer) ValidateResourceType(resourceType string) error {
	return nil
}

func (apier FakeAirDrumResource_APIer) ResourceInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceType string,
	api *API) (error, interface{}) {

	return nil, ""
}

func (apier FakeAirDrumResource_APIer) GetResourceCreationUrl(api *API,
	resourceType string) string {

	return ""
}

func (apier FakeAirDrumResource_APIer) GetResourceUrl(api *API,
	resourceType string,
	id string) string {

	return ""
}

func (apier FakeAirDrumResource_APIer) CreateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) (error, map[string]interface{}) {

	return nil, nil
}
func (apier FakeAirDrumResource_APIer) ReadResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) (error, map[string]interface{}, bool) {

	return nil, nil, true
}
func (apier FakeAirDrumResource_APIer) UpdateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) error {

	return nil
}
func (apier FakeAirDrumResource_APIer) DeleteResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) error {

	return nil
}
