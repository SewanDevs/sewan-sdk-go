package sewan_go_sdk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type FakeAirDrumResource_APIer struct{}

func (apier FakeAirDrumResource_APIer) CreateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (apier FakeAirDrumResource_APIer) ReadResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (apier FakeAirDrumResource_APIer) UpdateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	return nil
}
func (apier FakeAirDrumResource_APIer) DeleteResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	return nil
}
