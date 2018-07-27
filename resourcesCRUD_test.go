package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"testing"
)

func TestCreateResource(t *testing.T) {
	testCases := []struct {
		Id              int
		TC_clienter     Clienter
		ResourceType    string
		Creation_Err    error
		CreatedResource map[string]interface{}
	}{
		{
			1,
			ErrorResponse_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(REQ_ERR),
			nil,
		},
		{
			2,
			BadBodyResponse_StatusCreated_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("Creation of \"Unit test resource\" failed, " +
				"the response body is not a properly formated json :" +
				"\n\r\"invalid character '\"' after object key\""),
			nil,
		},
		{
			3,
			Error401_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(UNAUTHORIZED_MSG),
			nil,
		},
		{
			4,
			VDC_CreationSuccess_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			nil,
			VDC_READ_RESPONSE_MAP,
		},
		{
			5,
			VM_CreationSuccess_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			nil,
			NO_TEMPLATE_VM_MAP,
		},
		{
			6,
			CheckRedirectReqFailure_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Creation of \"Unit test resource\" failed, response reception " +
				"error : CheckRedirectReqFailure"),
			nil,
		},
		{
			7,
			BadBodyResponseContentType_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Unhandled api response type : image" +
				"\nPlease validate the configuration api url."),
			nil,
		},
		{
			8,
			StatusInternalServerError_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("<h1>Server Error (500)</h1>"),
			nil,
		},
		{
			8,
			StatusInternalServerError_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("<h1>Server Error (500)</h1>"),
			nil,
		},
		{
			9,
			VM_CreationSuccess_HttpClienterFake{},
			WRONG_RESOURCE_TYPE,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
			nil,
		},
	}
	var (
		sewan             *API
		err               error
		resp_creation_map map[string]interface{}
		resource_res      *schema.Resource
		d                 *schema.ResourceData
		diffs             string
	)
	apier := AirDrumResources_Apier{}

	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fake_client_tooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fake_schema_tooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}

	for _, testCase := range testCases {
		resource_res = resource(testCase.ResourceType)
		d = resource_res.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NAME_FIELD, "Unit test resource")
		fake_client_tooler.Client = testCase.TC_clienter
		err, resp_creation_map = apier.CreateResource(d,
			&fake_client_tooler,
			&fakeTemplates_tooler,
			&fake_schema_tooler,
			testCase.ResourceType,
			sewan)

		diffs = cmp.Diff(testCase.CreatedResource, resp_creation_map)
		switch {
		case err == nil || testCase.Creation_Err == nil:
			if !(err == nil && testCase.Creation_Err == nil) {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Creation_Err)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Creation_Err != nil:
			if resp_creation_map != nil {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, resp_creation_map, testCase.CreatedResource)
			}
			if err.Error() != testCase.Creation_Err.Error() {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Creation_Err.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) \n%s",
				testCase.Id, diffs)
		}
	}
}

//------------------------------------------------------------------------------
func TestReadResource(t *testing.T) {
	testCases := []struct {
		Id              int
		TC_clienter     Clienter
		ResourceType    string
		Read_Err        error
		ReadResource    map[string]interface{}
		Resource_exists bool
	}{
		{
			1,
			ErrorResponse_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(REQ_ERR),
			nil,
			true,
		},
		{
			2,
			BadBodyResponse_StatusOK_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Read of \"Unit test resource\" failed, response body json " +
				"error :\n\r\"invalid character '\"' after object key\""),
			nil,
			true,
		},
		{
			3,
			Error401_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New(UNAUTHORIZED_MSG),
			nil,
			true,
		},
		{
			4,
			Error404_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			nil,
			nil,
			false,
		},
		{
			5,
			CheckRedirectReqFailure_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("Read of \"Unit test resource\" state failed, response reception " +
				"error : CheckRedirectReqFailure"),
			nil,
			true,
		},
		{
			6,
			VDC_ReadSuccess_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			nil,
			VDC_READ_RESPONSE_MAP,
			true,
		},
		{
			7,
			VM_ReadSuccess_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			nil,
			NO_TEMPLATE_VM_MAP,
			true,
		},
		{
			8,
			BadBodyResponseContentType_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Unhandled api response type : image" +
				"\nPlease validate the configuration api url."),
			nil,
			true,
		},
		{
			9,
			StatusInternalServerError_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("<h1>Server Error (500)</h1>"),
			nil,
			true,
		},
		{
			10,
			VDC_ReadSuccess_HttpClienterFake{},
			WRONG_RESOURCE_TYPE,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
			nil,
			true,
		},
	}
	var (
		sewan             *API
		err               error
		resp_creation_map map[string]interface{}
		res_exists        bool
		resource_res      *schema.Resource
		d                 *schema.ResourceData
		diffs             string
	)
	Apier := AirDrumResources_Apier{}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fake_client_tooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fake_schema_tooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}

	for _, testCase := range testCases {
		resource_res = resource(testCase.ResourceType)
		d = resource_res.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NAME_FIELD, "Unit test resource")
		fake_client_tooler.Client = testCase.TC_clienter
		err, resp_creation_map, res_exists = Apier.ReadResource(d,
			&fake_client_tooler,
			&fakeTemplates_tooler,
			&fake_schema_tooler,
			testCase.ResourceType,
			sewan)
		diffs = cmp.Diff(testCase.ReadResource, resp_creation_map)
		switch {
		case err == nil || testCase.Read_Err == nil:
			if !((err == nil) && (testCase.Read_Err == nil)) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Read_Err)
			} else {
				switch {
				case res_exists != testCase.Resource_exists:
					t.Errorf("\n\nTC %d : Wrong read resource exists value"+
						"\n\rgot: \"%v\"\n\rwant: \"%v\"",
						testCase.Id, res_exists, testCase.Resource_exists)
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Read_Err != nil:
			if resp_creation_map != nil {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, resp_creation_map, testCase.ReadResource)
			}
			if err.Error() != testCase.Read_Err.Error() {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Read_Err.Error())
			}
		case res_exists != testCase.Resource_exists:
			t.Errorf("\n\nTC %d : Wrong read resource exists value"+
				"\n\rgot: \"%v\"\n\rwant: \"%v\"",
				testCase.Id, res_exists, testCase.Resource_exists)
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
				testCase.Id, diffs)
		}
	}
}

//------------------------------------------------------------------------------
func TestUpdateResource(t *testing.T) {
	testCases := []struct {
		Id           int
		TC_clienter  Clienter
		ResourceType string
		Update_Err   error
	}{
		{
			1,
			ErrorResponse_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(REQ_ERR),
		},
		{
			2,
			BadBodyResponse_StatusOK_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("Read of \"Unit test resource\" failed, response body json " +
				"error :\n\r\"invalid character '\"' after object key"),
		},
		{
			3,
			Error401_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(UNAUTHORIZED_MSG),
		},
		{
			4,
			Error404_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New(NOT_FOUND_MSG),
		},
		{
			5,
			CheckRedirectReqFailure_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Update of \"Unit test resource\" state failed, response reception " +
				"error : CheckRedirectReqFailure"),
		},
		{
			6,
			VDC_UpdateSuccess_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			nil,
		},
		{
			7,
			VM_UpdateSuccess_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			nil,
		},
		{
			8,
			BadBodyResponseContentType_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Unhandled api response type : image" +
				"\nPlease validate the configuration api url."),
		},
		{
			9,
			StatusInternalServerError_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("<h1>Server Error (500)</h1>"),
		},
		{
			10,
			VDC_UpdateSuccess_HttpClienterFake{},
			WRONG_RESOURCE_TYPE,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
		},
	}
	Apier := AirDrumResources_Apier{}
	var (
		sewan        *API
		err          error
		resource_res *schema.Resource
		d            *schema.ResourceData
	)
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fake_client_tooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fake_schema_tooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}

	for _, testCase := range testCases {
		resource_res = resource(testCase.ResourceType)
		d = resource_res.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NAME_FIELD, "Unit test resource")
		fake_client_tooler.Client = testCase.TC_clienter
		err = Apier.UpdateResource(d,
			&fake_client_tooler,
			&fakeTemplates_tooler,
			&fake_schema_tooler,
			testCase.ResourceType,
			sewan)

		switch {
		case err == nil || testCase.Update_Err == nil:
			if !(err == nil && testCase.Update_Err == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Update_Err)
			}
		case err.Error() != testCase.Update_Err.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, err.Error(), testCase.Update_Err.Error())
		}
	}
}

////------------------------------------------------------------------------------
func TestDeleteResource(t *testing.T) {
	testCases := []struct {
		Id           int
		TC_clienter  Clienter
		ResourceType string
		Delete_Err   error
	}{
		{
			1,
			ErrorResponse_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(REQ_ERR),
		},
		{
			2,
			BadBodyResponse_StatusOK_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New("Read of \"Unit test resource\" failed, response body json " +
				"error :\n\r\"invalid character '\"' after object key"),
		},
		{
			3,
			Error401_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New(UNAUTHORIZED_MSG),
		},
		{
			4,
			Error404_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New(NOT_FOUND_MSG),
		},
		{
			5,
			CheckRedirectReqFailure_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Deletion of \"Unit test resource\" state failed, response reception " +
				"error : CheckRedirectReqFailure"),
		},
		{
			6,
			VDC_DeleteSuccess_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			nil,
		},
		{
			7,
			VM_DeleteSuccess_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			nil,
		},
		{
			8,
			DeleteWRONGResponseBody_HttpClienterFake{},
			VDC_RESOURCE_TYPE,
			errors.New(DESTROY_WRONG_MSG),
		},
		{
			9,
			BadBodyResponseContentType_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("Unhandled api response type : image" +
				"\nPlease validate the configuration api url."),
		},
		{
			10,
			StatusInternalServerError_HttpClienterFake{},
			VM_RESOURCE_TYPE,
			errors.New("<h1>Server Error (500)</h1>"),
		},
		{
			7,
			VM_DeleteSuccess_HttpClienterFake{},
			WRONG_RESOURCE_TYPE,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
		},
	}
	var (
		sewan        *API
		err          error
		resource_res *schema.Resource
		d            *schema.ResourceData
	)
	Apier := AirDrumResources_Apier{}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fake_client_tooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fake_schema_tooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}

	for _, testCase := range testCases {
		resource_res = resource(testCase.ResourceType)
		d = resource_res.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NAME_FIELD, "Unit test resource")
		fake_client_tooler.Client = testCase.TC_clienter
		err = Apier.DeleteResource(d,
			&fake_client_tooler,
			&fakeTemplates_tooler,
			&fake_schema_tooler,
			testCase.ResourceType,
			sewan)

		switch {
		case err == nil || testCase.Delete_Err == nil:
			if !(err == nil && testCase.Delete_Err == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Delete_Err)
			}
		case err.Error() != testCase.Delete_Err.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, err.Error(), testCase.Delete_Err.Error())
		}
	}
}
