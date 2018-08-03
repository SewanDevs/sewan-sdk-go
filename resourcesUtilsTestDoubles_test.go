package sewan_go_sdk

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
)

type FakeResourceResourceer struct{}

func (resourceer FakeResourceResourceer) ValidateStatus(api *API,
	resourceType string,
	clientTooler ClientTooler) error {
	var err error
	switch {
	case api.URL != rightApiUrl:
		err = errors.New(wrongApiUrlError)
	case api.Token != rightApiToken:
		err = errors.New(wrongTokenError)
	default:
		err = nil
	}
	return err
}

func (resourceer FakeResourceResourceer) ValidateResourceType(resourceType string) error {
	return nil
}

func (resourceer FakeResourceResourceer) ResourceInstanceCreate(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	resourceType string,
	api *API) (interface{}, error) {

	return "", nil
}

func (resourceer FakeResourceResourceer) GetResourceCreationUrl(api *API,
	resourceType string) string {

	return ""
}

func (resourceer FakeResourceResourceer) GetResourceUrl(api *API,
	resourceType string,
	id string) string {

	return ""
}
