package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
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

func TestCheckCloudDcApiStatus(t *testing.T) {
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
		err := fakeApi_tools.CheckCloudDcApiStatus(testCase.Input_api,
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
	resourceName := "Unit test resource creation"
	testCases := []struct {
		Id              int
		TcClienter      Clienter
		ResourceType    string
		Creation_Err    error
		CreatedResource map[string]interface{}
	}{
		{
			1,
			VmCreationSuccessHttpClienterFake{},
			VmResourceType,
			nil,
			noTemplateVmMap,
		},
		{
			2,
			HttpClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			3,
			ResourceCreationFailureHttpClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(creationOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			4,
			HandleRespErrHttpClienterFake{},
			VmResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
	}
	apier := AirDrumResourcesApier{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		respCreationMap, err := apier.CreateResource(d,
			&fakeClientTooler,
			&fakeTemplates_tooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(testCase.CreatedResource, respCreationMap)
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

func TestReadResource(t *testing.T) {
	resourceName := "Unit test resource read"
	testCases := []struct {
		Id           int
		TcClienter   Clienter
		ResourceType string
		ReadError    error
		ReadResource map[string]interface{}
	}{
		{
			1,
			VmReadSuccessHttpClienterFake{},
			VmResourceType,
			nil,
			noTemplateVmMap,
		},
		{
			2,
			VdcReadSuccessHttpClienterFake{},
			VdcResourceType,
			nil,
			vdcReadResponseMap,
		},
		{
			3,
			HttpClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			4,
			ResourceReadFailureHttpClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(readOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			5,
			Error404HttpClienterFake{},
			VdcResourceType,
			ErrResourceNotExist,
			map[string]interface{}{},
		},
		{
			6,
			HandleRespErrHttpClienterFake{},
			VmResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
	}
	Apier := AirDrumResourcesApier{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		respCreationMap, err := Apier.ReadResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(respCreationMap, testCase.ReadResource)
		switch {
		case err == nil || testCase.ReadError == nil:
			if !((err == nil) && (testCase.ReadError == nil)) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.ReadError)
			} else {
				if diffs != "" {
					t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.ReadError != nil:
			if cmp.Diff(respCreationMap, map[string]interface{}{}) != "" {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, respCreationMap, testCase.ReadResource)
			}
			if err.Error() != testCase.ReadError.Error() {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.ReadError.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
				testCase.Id, diffs)
		}
	}
}

func TestUpdateResource(t *testing.T) {
	resourceName := "Unit test resource update"
	testCases := []struct {
		Id           int
		TcClienter   Clienter
		ResourceType string
		Update_Err   error
	}{
		{
			1,
			ResourceUpdateSuccessHttpClienterFake,
			VmResourceType,
			nil,
		},
		{
			2,
			HttpClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceUpdateFailureHttpClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(updateOperation,
				resourceName,
				errEmptyResp),
		},
	}
	Apier := AirDrumResourcesApier{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		err := Apier.UpdateResource(d,
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

func TestDeleteResource(t *testing.T) {
	resourceName := "Unit test resource deletion"
	testCases := []struct {
		Id           int
		TcClienter   Clienter
		ResourceType string
		Delete_Err   error
	}{
		{
			1,
			ResourceDeletionSuccessHttpClienterFake,
			VmResourceType,
			nil,
		},
		{
			2,
			HttpClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceDeletionFailureHttpClienterFake{},
			VdcResourceType,
			errDoCrudRequestsBuilder(deleteOperation,
				resourceName,
				errEmptyResp),
		},
	}
	Apier := AirDrumResourcesApier{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		err := Apier.DeleteResource(d,
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
