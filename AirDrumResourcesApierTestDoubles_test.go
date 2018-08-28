package sewan_go_sdk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type FakeAirDrumResourceAPIer struct{}

func (apier FakeAirDrumResourceAPIer) CreateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (apier FakeAirDrumResourceAPIer) ReadResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (apier FakeAirDrumResourceAPIer) UpdateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	return nil
}
func (apier FakeAirDrumResourceAPIer) DeleteResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	resourceTooler *ResourceTooler,
	resourceType string,
	sewan *API) error {
	return nil
}
