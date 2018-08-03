package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"testing"
)

func TestResourceInstanceCreate(t *testing.T) {
	testCases := []struct {
		Id            int
		D             *schema.ResourceData
		TC_Clienter   Clienter
		TC_Templater  Templater
		Resource_type string
		Error         error
		VmInstance    interface{}
	}{
		{
			1,
			vmSchemaInit(noTemplateVmMap),
			GetTemplatesList_Success_HttpClienterFake{},
			TemplaterDummy{},
			VmResourceField,
			nil,
			vmInstanceNoTemplateVmMap(),
		},
		{
			2,
			vmSchemaInit(existingTemplateNoAdditionalDiskVmMap),
			GetTemplatesList_Success_HttpClienterFake{},
			existingTemplateNoAdditionalDiskVmMap_TemplaterFake{},
			VmResourceField,
			nil,
			FakeVmInstanceExistingTemplateNoAdditionalDiskVmMap(),
		},
		{
			3,
			vmSchemaInit(existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap),
			GetTemplatesList_Success_HttpClienterFake{},
			existingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap_TemplaterFake{},
			VmResourceField,
			nil,
			FakeVmInstanceExistingTemplateWithAdditionalAndModifiedDisksAndNicsVmMap(),
		},
		{
			4,
			vmSchemaInit(nonExistingTemplateVmMap),
			GetTemplatesList_Success_HttpClienterFake{},
			UnexistingTemplate_TemplaterFake{},
			VmResourceField,
			errors.New("Unavailable template : windows95"),
			VM{},
		},
		{
			5,
			vdcSchemaInit(vdcCreationMap),
			nil,
			TemplaterDummy{},
			VdcResourceField,
			nil,
			FakeVdcInstanceVdcCreationMap(),
		},
		{
			6,
			vdcSchemaInit(vdcCreationMap),
			GetTemplatesList_Success_HttpClienterFake{},
			TemplaterDummy{},
			wrongResourceType,
			errors.New("Resource of type \"a_non_supportedResource_type\" not supported," +
				"list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
			nil,
		},
		{
			7,
			vmSchemaInit(nonExistingTemplateVmMap),
			GetTemplatesList_Failure_HttpClienterFake{},
			UnexistingTemplate_TemplaterFake{},
			VmResourceField,
			errors.New("GetTemplatesList() error"),
			VM{},
		},
		{
			8,
			vmSchemaInit(existingTemplateNoAdditionalDiskVmMap),
			GetTemplatesList_Success_HttpClienterFake{},
			Template_FormatError_TemplaterFake{},
			VmResourceField,
			errors.New("Template missing fields : " + "\"" + NameField + "\" " +
				"\"" + OsField + "\" " +
				"\"" + RamField + "\" " +
				"\"" + CpuField + "\" " +
				"\"" + EnterpriseField + "\" " +
				"\"" + DisksField + "\" " +
				"\"" + DatacenterField + "\" "),
			VM{},
		},
		{
			9,
			vmSchemaInit(instanceNumberFieldUnitTestVmInstance),
			GetTemplatesList_Success_HttpClienterFake{},
			instanceNumberFieldUnitTestVmInstance_MAP_TemplaterFake{},
			VmResourceField,
			nil,
			FakeVmInstanceInstanceNumberFieldUnitTestVmInstance_MAP(),
		},
	}

	var (
		sewan    *API
		err      error = nil
		instance interface{}
		diffs    string
	)
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	fakeClientTooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	for _, testCase := range testCases {
		fakeClientTooler.Client = testCase.TC_Clienter
		fakeTemplates_tooler.TemplatesTools = testCase.TC_Templater
		instance, err = fakeResourceTooler.Resource.ResourceInstanceCreate(testCase.D,
			&fakeClientTooler,
			&fakeTemplates_tooler,
			testCase.Resource_type,
			sewan)
		diffs = cmp.Diff(instance, testCase.VmInstance)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : ResourceInstanceCreate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Error)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong ResourceInstanceCreate() "+
						"created instance (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			switch {
			case diffs != "":
				t.Errorf("\n\nTC %d : Wrong ResourceInstanceCreate() "+
					"created instance (-got +want) :\n%s",
					testCase.Id, diffs)
			case err.Error() != testCase.Error.Error():
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			}
		}
	}
}

func TestGetResourceUrl(t *testing.T) {
	testCases := []struct {
		Id     int
		api    API
		vm_id  string
		vm_url string
	}{
		{1,
			API{
				rightApiToken,
				rightApiUrl,
				&http.Client{},
			},
			"42",
			rightVmUrlQuaranteDeux,
		},
		{2,
			API{
				rightApiToken,
				rightApiUrl,
				&http.Client{},
			},
			"PATATE",
			rightVmUrlPatate,
		},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		s_vm_url := fakeResourceTooler.Resource.GetResourceUrl(&testCase.api,
			VmResourceField,
			testCase.vm_id)
		switch {
		case s_vm_url != testCase.vm_url:
			t.Errorf("VM url was incorrect,\n\rgot: \"%s\"\n\rwant: \"%s\"",
				s_vm_url, testCase.vm_url)
		}
	}
}

func TestGetResourceCreationUrl(t *testing.T) {
	testCases := []struct {
		Id                    int
		api                   API
		resource_creation_url string
	}{
		{1,
			API{
				rightApiToken,
				rightApiUrl,
				&http.Client{},
			},
			rightVmCreationApiUrl,
		},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		s_resource_creation_url := fakeResourceTooler.Resource.GetResourceCreationUrl(&testCase.api,
			VmResourceField)
		switch {
		case s_resource_creation_url != testCase.resource_creation_url:
			t.Errorf("resource api creation url was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				s_resource_creation_url, testCase.resource_creation_url)
		}
	}
}

func TestValidateStatus(t *testing.T) {
	testCases := []struct {
		Id           int
		Api          API
		Err          error
		ResourceType string
	}{
		{1,
			API{
				rightApiToken,
				rightApiUrl,
				&http.Client{},
			},
			nil,
			VmResourceField,
		},
		{2,
			API{
				wrongApiToken,
				rightApiUrl,
				&http.Client{},
			},
			errors.New("401 Unauthorized{\"detail\":\"Invalid token.\"}"),
			VmResourceField,
		},
		{3,
			API{
				rightApiToken,
				wrongApiUrl,
				&http.Client{},
			},
			errors.New("Could not get a proper json response from \"" +
				wrongApiUrl + errApiDownOrWrongApiUrl),
			VmResourceField,
		},
		{4,
			API{
				wrongApiToken,
				wrongApiUrl,
				&http.Client{},
			},
			errors.New("Could not get a proper json response from \"" +
				wrongApiUrl + errApiDownOrWrongApiUrl),
			VmResourceField,
		},
		{5,
			API{
				rightApiToken,
				noRespBodyApiUrl,
				&http.Client{},
			},
			errors.New("Could not get a response body from \"" +
				noRespBodyApiUrl + errApiDownOrWrongApiUrl),
			VmResourceField,
		},
		{6,
			API{
				rightApiToken,
				noRespApiUrl,
				&http.Client{},
			},
			errors.New("Could not get a response from \"" +
				noRespApiUrl + errApiDownOrWrongApiUrl),
			VmResourceField,
		},
	}
	fakeResourceTooler := &ResourceTooler{
		Resource: ResourceResourceer{},
	}
	clientTooler := ClientTooler{
		Client: FakeHttpClienter{},
	}
	var apiClientErr error
	for _, testCase := range testCases {
		apiClientErr = fakeResourceTooler.Resource.ValidateStatus(&testCase.Api,
			testCase.ResourceType,
			clientTooler)
		switch {
		case apiClientErr == nil || testCase.Err == nil:
			if !(apiClientErr == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : ValidateStatus() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, apiClientErr, testCase.Err)
			}
		case apiClientErr.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : ValidateStatus() error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, apiClientErr.Error(), testCase.Err.Error())
		}
	}
}

func CreateTestResourceSchema(id interface{}) *schema.ResourceData {
	vm_res := resourceVm()
	d := vm_res.TestResourceData()
	d.SetId(id.(string))
	return d
}
