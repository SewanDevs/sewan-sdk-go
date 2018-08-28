package sewan_go_sdk

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

type FakeResourceResourceer struct{}

func (resourceer FakeResourceResourceer) validateStatus(api *API,
	resourceType string,
	clientTooler ClientTooler) error {
	var err error
	switch {
	case api.URL != rightAPIURL:
		err = errors.New(wrongAPIURLError)
	case api.Token != rightAPIToken:
		err = errors.New(wrongTokenError)
	default:
		err = nil
	}
	return err
}

func (resourceer FakeResourceResourceer) validateResourceType(resourceType string) error {
	return nil
}

func (resourceer FakeResourceResourceer) resourceInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceType string,
	api *API) (interface{}, error) {

	return "", nil
}

func (resourceer FakeResourceResourceer) getResourceCreationURL(api *API,
	resourceType string) string {

	return ""
}

func (resourceer FakeResourceResourceer) getResourceURL(api *API,
	resourceType string,
	id string) string {

	return ""
}
