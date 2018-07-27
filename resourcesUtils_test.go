package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"testing"
)

//------------------------------------------------------------------------------
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
			vmSchemaInit(NO_TEMPLATE_VM_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			TemplaterDummy{},
			VM_RESOURCE_TYPE,
			nil,
			vmInstanceNO_TEMPLATE_VM_MAP(),
		},
		{
			2,
			vmSchemaInit(EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP_TemplaterFake{},
			VM_RESOURCE_TYPE,
			nil,
			FakeVmInstance_EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP(),
		},
		{
			3,
			vmSchemaInit(EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP_TemplaterFake{},
			VM_RESOURCE_TYPE,
			nil,
			FakeVmInstance_EXISTING_TEMPLATE_WITH_ADDITIONAL_AND_MODIFIED_NICS_AND_DISKS_VM_MAP(),
		},
		{
			4,
			vmSchemaInit(NON_EXISTING_TEMPLATE_VM_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			UnexistingTemplate_TemplaterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Unavailable template : windows95"),
			VM{},
		},
		{
			5,
			vdcSchemaInit(VDC_CREATION_MAP),
			nil,
			TemplaterDummy{},
			VDC_RESOURCE_TYPE,
			nil,
			FakeVdcInstance_VDC_CREATION_MAP(),
		},
		{
			6,
			vdcSchemaInit(VDC_CREATION_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			TemplaterDummy{},
			WRONG_RESOURCE_TYPE,
			errors.New("Resource of type \"a_non_supportedResource_type\" not supported," +
				"list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
			nil,
		},
		{
			7,
			vmSchemaInit(NON_EXISTING_TEMPLATE_VM_MAP),
			GetTemplatesList_Failure_HttpClienterFake{},
			UnexistingTemplate_TemplaterFake{},
			VM_RESOURCE_TYPE,
			errors.New("GetTemplatesList() error"),
			VM{},
		},
		{
			8,
			vmSchemaInit(EXISTING_TEMPLATE_NO_ADDITIONAL_DISK_VM_MAP),
			GetTemplatesList_Success_HttpClienterFake{},
			Template_FormatError_TemplaterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Template missing fields : " + "\"" + NAME_FIELD + "\" " +
				"\"" + OS_FIELD + "\" " +
				"\"" + RAM_FIELD + "\" " +
				"\"" + CPU_FIELD + "\" " +
				"\"" + ENTERPRISE_FIELD + "\" " +
				"\"" + DISKS_FIELD + "\" " +
				"\"" + DATACENTER_FIELD + "\" "),
			VM{},
		},
		{
			9,
			vmSchemaInit(INSTANCE_NUMBER_FIELD_UNIT_TEST_VM_INSTANCE),
			GetTemplatesList_Success_HttpClienterFake{},
			INSTANCE_NUMBER_FIELD_UNIT_TEST_VM_INSTANCE_MAP_TemplaterFake{},
			VM_RESOURCE_TYPE,
			nil,
			FakeVmInstance_INSTANCE_NUMBER_FIELD_UNIT_TEST_VM_INSTANCE_MAP(),
		},
	}

	var (
		sewan    *API
		err      error = nil
		instance interface{}
		diffs    string
	)
	apiTools := APITooler{
		Api: AirDrumResources_Apier{},
	}
	fake_client_tooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{}
	fake_schema_tooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}

	for _, testCase := range testCases {
		fake_client_tooler.Client = testCase.TC_Clienter
		fakeTemplates_tooler.TemplatesTools = testCase.TC_Templater
		err, instance = apiTools.Api.ResourceInstanceCreate(testCase.D,
			&fake_client_tooler,
			&fakeTemplates_tooler,
			&fake_schema_tooler,
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

//------------------------------------------------------------------------------
func TestGetResourceUrl(t *testing.T) {
	testCases := []struct {
		Id     int
		api    API
		vm_id  string
		vm_url string
	}{
		{1,
			API{
				RIGHT_API_TOKEN,
				RIGHT_API_URL,
				&http.Client{},
			},
			"42",
			RIGHT_VM_URL_42,
		},
		{2,
			API{
				RIGHT_API_TOKEN,
				RIGHT_API_URL,
				&http.Client{},
			},
			"PATATE",
			RIGHT_VM_URL_PATATE,
		},
	}

	apiTools := APITooler{
		Api: AirDrumResources_Apier{},
	}

	for _, testCase := range testCases {
		s_vm_url := apiTools.Api.GetResourceUrl(&testCase.api,
			VM_RESOURCE_TYPE,
			testCase.vm_id)

		switch {
		case s_vm_url != testCase.vm_url:
			t.Errorf("VM url was incorrect,\n\rgot: \"%s\"\n\rwant: \"%s\"",
				s_vm_url, testCase.vm_url)
		}
	}
}

//------------------------------------------------------------------------------
func TestGetResourceCreationUrl(t *testing.T) {
	testCases := []struct {
		Id                    int
		api                   API
		resource_creation_url string
	}{
		{1,
			API{
				RIGHT_API_TOKEN,
				RIGHT_API_URL,
				&http.Client{},
			},
			RIGHT_VM_CREATION_API_URL,
		},
	}

	apiTools := APITooler{
		Api: AirDrumResources_Apier{},
	}

	for _, testCase := range testCases {
		s_resource_creation_url := apiTools.Api.GetResourceCreationUrl(&testCase.api,
			VM_RESOURCE_TYPE)

		switch {
		case s_resource_creation_url != testCase.resource_creation_url:
			t.Errorf("resource api creation url was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				s_resource_creation_url, testCase.resource_creation_url)
		}
	}
}

//------------------------------------------------------------------------------
func TestValidateStatus(t *testing.T) {
	testCases := []struct {
		Id           int
		Api          API
		Err          error
		ResourceType string
	}{
		{1,
			API{
				RIGHT_API_TOKEN,
				RIGHT_API_URL,
				&http.Client{},
			},
			nil,
			VM_RESOURCE_TYPE,
		},
		{2,
			API{
				WRONG_API_TOKEN,
				RIGHT_API_URL,
				&http.Client{},
			},
			errors.New("401 Unauthorized{\"detail\":\"Invalid token.\"}"),
			VM_RESOURCE_TYPE,
		},
		{3,
			API{
				RIGHT_API_TOKEN,
				WRONG_API_URL,
				&http.Client{},
			},
			errors.New("Could not get a proper json response from \"" +
				WRONG_API_URL + "\", the api is down or this url is wrong."),
			VM_RESOURCE_TYPE,
		},
		{4,
			API{
				WRONG_API_TOKEN,
				WRONG_API_URL,
				&http.Client{},
			},
			errors.New("Could not get a proper json response from \"" +
				WRONG_API_URL + "\", the api is down or this url is wrong."),
			VM_RESOURCE_TYPE,
		},
		{5,
			API{
				RIGHT_API_TOKEN,
				NO_RESP_BODY_API_URL,
				&http.Client{},
			},
			errors.New("Could not get a response body from \"" +
				NO_RESP_BODY_API_URL + "\", the api is down or this url is wrong."),
			VM_RESOURCE_TYPE,
		},
		{6,
			API{
				RIGHT_API_TOKEN,
				NO_RESP_API_URL,
				&http.Client{},
			},
			errors.New("Could not get a response from \"" +
				NO_RESP_API_URL + "\", the api is down or this url is wrong."),
			VM_RESOURCE_TYPE,
		},
	}

	apiTooler := APITooler{
		Api: AirDrumResources_Apier{},
	}
	clientTooler := ClientTooler{
		Client: FakeHttpClienter{},
	}
	var apiClientErr error

	for _, testCase := range testCases {
		apiClientErr = apiTooler.Api.ValidateStatus(&testCase.Api,
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

func Create_test_resource_schema(id interface{}) *schema.ResourceData {
	vm_res := resource_vm()
	d := vm_res.TestResourceData()
	d.SetId(id.(string))
	return d
}

func TestDeleteTerraformResource(t *testing.T) {
	d := Create_test_resource_schema("resource to delete")
	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	schemaTooler.SchemaTools.DeleteTerraformResource(d)
	if d.Id() != "" {
		t.Errorf("Deletion of unit test resource failed.")
	}
}

func TestUpdateLocalResourceState_AND_ReadElement(t *testing.T) {
	testCases := []struct {
		Id           int
		Vm_map       map[string]interface{}
		Vm_Id_string string
	}{
		{
			1,
			TEST_UPDATE_VM_MAP,
			"unit test vm",
		},
		{
			2,
			TEST_UPDATE_VM_MAP_FLOATID,
			"121212.12",
		},
		{
			3,
			TEST_UPDATE_VM_MAP_INTID,
			"1212",
		},
	}
	var (
		d     *schema.ResourceData
		diffs string
	)
	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	for _, testCase := range testCases {
		d = Create_test_resource_schema(testCase.Vm_Id_string)
		schemaTooler.SchemaTools.UpdateLocalResourceState(testCase.Vm_map,
			d,
			&schemaTooler)
		for key, value := range testCase.Vm_map {
			diffs = cmp.Diff(d.Get(key), value)
			switch {
			case key != ID_FIELD:
				if diffs != "" {
					t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
						testCase.Id, key, diffs)
				}
			default:
				if d.Id() != testCase.Vm_Id_string {
					t.Errorf("\n\nTC %d : Update of Id reserved field failed :\n\rGot :%s\n\rWant :%s",
						testCase.Id, d.Id(), testCase.Vm_Id_string)
				}
			}
		}
	}
}
