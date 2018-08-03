package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		Id          int
		Input_token string
		Input_url   string
		Output_api  API
	}{
		{1,
			wrongApiToken,
			rightApiUrl,
			API{wrongApiToken, rightApiUrl, nil},
		},
		{2,
			rightApiToken,
			wrongApiUrl,
			API{rightApiToken, wrongApiUrl, nil},
		},
		{3,
			wrongApiToken,
			wrongApiUrl,
			API{wrongApiToken, wrongApiUrl, nil},
		},
		{4,
			rightApiToken,
			rightApiUrl,
			API{rightApiToken, rightApiUrl, nil},
		},
	}
	fakeApi_tools := APITooler{
		Api: FakeAirDrumResource_APIer{},
	}
	for _, testCase := range testCases {
		api := fakeApi_tools.New(
			testCase.Input_token,
			testCase.Input_url,
		)
		switch {
		case api.Token != testCase.Output_api.Token:
			t.Errorf("\n\nTC %d : API token error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, api.Token, testCase.Output_api.Token)
		case api.URL != testCase.Output_api.URL:
			t.Errorf("\n\nTC %d : API token error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, api.URL, testCase.Output_api.URL)
		}
	}
}

func TestCheckStatus(t *testing.T) {
	testCases := []struct {
		Id               int
		Input_api        *API
		TCResourceTooler Resourceer
		Err              error
	}{
		{1,
			&API{
				wrongApiToken,
				rightApiUrl,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongTokenError),
		},
		{2,
			&API{
				rightApiToken,
				wrongApiUrl,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongApiUrlError),
		},
		{3,
			&API{
				wrongApiToken,
				wrongApiUrl,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongApiUrlError),
		},
		{4,
			&API{
				rightApiToken,
				rightApiUrl,
				&http.Client{},
			},
			FakeResourceResourceer{},
			nil,
		},
	}
	fakeApi_tools := APITooler{}
	fakeClientTooler := &ClientTooler{
		Client: HttpClienter{},
	}
	fakeResourceTooler := &ResourceTooler{}
	for _, testCase := range testCases {
		fakeApi_tools.Api = FakeAirDrumResource_APIer{}
		fakeResourceTooler.Resource = testCase.TCResourceTooler
		err := fakeApi_tools.CheckStatus(testCase.Input_api,
			fakeClientTooler,
			fakeResourceTooler)
		switch {
		case err == nil || testCase.Err == nil:
			if !(err == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : Check API error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Err)
			}
		case err.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : Check API error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.Id, err.Error(), testCase.Err.Error())
		}
	}
}

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
			VmResourceField,
			errors.New("Creation of \"" + "Unit test resource" +
				"\" failed, response reception error : " + reqErr),
			map[string]interface{}{},
		},
		//{
		//	2,
		//	BadBodyResponse_StatusCreated_HttpClienterFake{},
		//	VdcResourceField,
		//	errors.New(errJsonFormat +
		//		"\n\r\"invalid character '\"' after object key\""),
		//	map[string]interface{}{},
		//},
		//{
		//	3,
		//	Error401_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New(unauthorizedMsg),
		//	map[string]interface{}{},
		//},
		//{
		//	4,
		//	VDC_CreationSuccess_HttpClienterFake{},
		//	VdcResourceField,
		//	nil,
		//	vdcReadResponseMap,
		//},
		//{
		//	5,
		//	VM_CreationSuccess_HttpClienterFake{},
		//	VmResourceField,
		//	nil,
		//	noTemplateVmMap,
		//},
		//{
		//	6,
		//	CheckRedirectReqFailure_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New("Creation of \"Unit test resource\" failed, response reception " +
		//		"error : CheckRedirectReqFailure"),
		//	map[string]interface{}{},
		//},
		//{
		//	7,
		//	BadBodyResponseContentType_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New(errorApiUnhandledImageType +
		//		errValidateApiUrl),
		//	map[string]interface{}{},
		//},
		//{
		//	8,
		//	StatusInternalServerError_HttpClienterFake{},
		//	VdcResourceField,
		//	errors.New("<h1>Server Error (500)</h1>"),
		//	map[string]interface{}{},
		//},
		//{
		//	8,
		//	StatusInternalServerError_HttpClienterFake{},
		//	VdcResourceField,
		//	errors.New("<h1>Server Error (500)</h1>"),
		//	map[string]interface{}{},
		//},
		//{
		//	9,
		//	VM_CreationSuccess_HttpClienterFake{},
		//	wrongResourceType,
		//	errors.New("Resource of type \"a_non_supportedResource_type\"" +
		//		" not supported,list of accepted resource types :\n\r" +
		//		"- \"vdc\"\n\r" +
		//		"- \"vm\""),
		//	map[string]interface{}{},
		//},
	}
	var (
		sewan            *API
		resourceResponse *schema.Resource
		d                *schema.ResourceData
		diffs            string
	)
	apier := AirDrumResources_Apier{}

	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse = resource(testCase.ResourceType)
		d = resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, "Unit test resource")
		fakeClientTooler.Client = testCase.TC_clienter
		respCreationMap, err := apier.CreateResource(d,
			&fakeClientTooler,
			&fakeTemplates_tooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)

		diffs = cmp.Diff(testCase.CreatedResource, respCreationMap)
		switch {
		case err == nil || testCase.Creation_Err == nil:
			if !(err == nil && testCase.Creation_Err == nil) {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Creation_Err)
			} else {
				switch {
				//case (testCase.CreateResource == map[string]interface{}{}) && (respCreationMap == map[string]interface{}{}):
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Creation_Err != nil:
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
		//{
		//	1,
		//	ErrorResponse_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New(reqErr),
		//	nil,
		//	true,
		//},
		//{
		//	2,
		//	BadBodyResponse_StatusOK_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New("Read of \"Unit test resource\" failed, response body json " +
		//		"error :\n\r\"invalid character '\"' after object key\""),
		//	nil,
		//	true,
		//},
		//{
		//	3,
		//	Error401_HttpClienterFake{},
		//	VdcResourceField,
		//	errors.New(unauthorizedMsg),
		//	nil,
		//	true,
		//},
		//{
		//	4,
		//	Error404_HttpClienterFake{},
		//	VmResourceField,
		//	nil,
		//	nil,
		//	false,
		//},
		//{
		//	5,
		//	CheckRedirectReqFailure_HttpClienterFake{},
		//	VdcResourceField,
		//	errors.New("Read of \"Unit test resource\" state failed, response reception " +
		//		"error : CheckRedirectReqFailure"),
		//	nil,
		//	true,
		//},
		//{
		//	6,
		//	VDC_ReadSuccess_HttpClienterFake{},
		//	VdcResourceField,
		//	nil,
		//	vdcReadResponseMap,
		//	true,
		//},
		//{
		//	7,
		//	VM_ReadSuccess_HttpClienterFake{},
		//	VmResourceField,
		//	nil,
		//	noTemplateVmMap,
		//	true,
		//},
		//{
		//	8,
		//	BadBodyResponseContentType_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New(errorApiUnhandledImageType +
		//		errValidateApiUrl),
		//	nil,
		//	true,
		//},
		//{
		//	9,
		//	StatusInternalServerError_HttpClienterFake{},
		//	VmResourceField,
		//	errors.New("<h1>Server Error (500)</h1>"),
		//	nil,
		//	true,
		//},
		//{
		//	10,
		//	VDC_ReadSuccess_HttpClienterFake{},
		//	wrongResourceType,
		//	errors.New("Resource of type \"a_non_supportedResource_type\"" +
		//		" not supported,list of accepted resource types :\n\r" +
		//		"- \"vdc\"\n\r" +
		//		"- \"vm\""),
		//	nil,
		//	true,
		//},
	}
	var (
		sewan            *API
		res_exists       bool
		resourceResponse *schema.Resource
		d                *schema.ResourceData
		diffs            string
	)
	Apier := AirDrumResources_Apier{}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse = resource(testCase.ResourceType)
		d = resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, "Unit test resource")
		fakeClientTooler.Client = testCase.TC_clienter
		respCreationMap, err := Apier.ReadResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs = cmp.Diff(testCase.ReadResource, respCreationMap)
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
			if respCreationMap != nil {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, respCreationMap, testCase.ReadResource)
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
			VmResourceField,
			errors.New(reqErr),
		},
		{
			2,
			BadBodyResponse_StatusOK_HttpClienterFake{},
			VdcResourceField,
			errors.New("Read of \"Unit test resource\" failed, response body json " +
				"error :\n\r\"invalid character '\"' after object key"),
		},
		{
			3,
			Error401_HttpClienterFake{},
			VmResourceField,
			errors.New(unauthorizedMsg),
		},
		{
			4,
			Error404_HttpClienterFake{},
			VdcResourceField,
			errors.New(notFoundRespMsg),
		},
		{
			5,
			CheckRedirectReqFailure_HttpClienterFake{},
			VmResourceField,
			errors.New("Update of \"Unit test resource\" state failed, response reception " +
				"error : CheckRedirectReqFailure"),
		},
		{
			6,
			VDC_UpdateSuccess_HttpClienterFake{},
			VdcResourceField,
			nil,
		},
		{
			7,
			VM_UpdateSuccess_HttpClienterFake{},
			VmResourceField,
			nil,
		},
		{
			8,
			BadBodyResponseContentType_HttpClienterFake{},
			VmResourceField,
			errors.New(errorApiUnhandledImageType +
				errValidateApiUrl),
		},
		{
			9,
			StatusInternalServerError_HttpClienterFake{},
			VmResourceField,
			errors.New("<h1>Server Error (500)</h1>"),
		},
		{
			10,
			VDC_UpdateSuccess_HttpClienterFake{},
			wrongResourceType,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
		},
	}
	Apier := AirDrumResources_Apier{}
	var (
		sewan            *API
		err              error
		resourceResponse *schema.Resource
		d                *schema.ResourceData
	)
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse = resource(testCase.ResourceType)
		d = resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, "Unit test resource")
		fakeClientTooler.Client = testCase.TC_clienter
		err = Apier.UpdateResource(d,
			&fakeClientTooler,
			&fakeTemplates_tooler,
			&fakeResourceTooler,
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
			VmResourceField,
			errors.New(reqErr),
		},
		{
			2,
			BadBodyResponse_StatusOK_HttpClienterFake{},
			VdcResourceField,
			errors.New("Read of \"Unit test resource\" failed, response body json " +
				"error :\n\r\"invalid character '\"' after object key"),
		},
		{
			3,
			Error401_HttpClienterFake{},
			VmResourceField,
			errors.New(unauthorizedMsg),
		},
		{
			4,
			Error404_HttpClienterFake{},
			VdcResourceField,
			errors.New(notFoundRespMsg),
		},
		{
			5,
			CheckRedirectReqFailure_HttpClienterFake{},
			VmResourceField,
			errors.New("Deletion of \"Unit test resource\" state failed, response reception " +
				"error : CheckRedirectReqFailure"),
		},
		{
			6,
			VDC_DeleteSuccess_HttpClienterFake{},
			VdcResourceField,
			nil,
		},
		{
			7,
			VM_DeleteSuccess_HttpClienterFake{},
			VmResourceField,
			nil,
		},
		{
			8,
			DeleteWRONGResponseBody_HttpClienterFake{},
			VdcResourceField,
			errors.New(destroyWrongMsg),
		},
		{
			9,
			BadBodyResponseContentType_HttpClienterFake{},
			VmResourceField,
			errors.New(errorApiUnhandledImageType +
				errValidateApiUrl),
		},
		{
			10,
			StatusInternalServerError_HttpClienterFake{},
			VmResourceField,
			errors.New("<h1>Server Error (500)</h1>"),
		},
		{
			7,
			VM_DeleteSuccess_HttpClienterFake{},
			wrongResourceType,
			errors.New("Resource of type \"a_non_supportedResource_type\"" +
				" not supported,list of accepted resource types :\n\r" +
				"- \"vdc\"\n\r" +
				"- \"vm\""),
		},
	}
	var (
		sewan            *API
		err              error
		resourceResponse *schema.Resource
		d                *schema.ResourceData
	)
	Apier := AirDrumResources_Apier{}
	sewan = &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse = resource(testCase.ResourceType)
		d = resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, "Unit test resource")
		fakeClientTooler.Client = testCase.TC_clienter
		err = Apier.DeleteResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
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
