package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		ID         int
		InputToken string
		InputURL   string
		OutputAPI  API
	}{
		{1,
			wrongAPIToken,
			rightAPIURL,
			API{wrongAPIToken, rightAPIURL, nil},
		},
		{2,
			rightAPIToken,
			wrongAPIURL,
			API{rightAPIToken, wrongAPIURL, nil},
		},
		{3,
			wrongAPIToken,
			wrongAPIURL,
			API{wrongAPIToken, wrongAPIURL, nil},
		},
		{4,
			rightAPIToken,
			rightAPIURL,
			API{rightAPIToken, rightAPIURL, nil},
		},
	}
	fakeAPItools := APITooler{
		APIImplementer: FakeAirDrumResourceAPIer{},
	}
	for _, testCase := range testCases {
		api := fakeAPItools.New(
			testCase.InputToken,
			testCase.InputURL,
		)
		switch {
		case api.Token != testCase.OutputAPI.Token:
			t.Errorf("\n\nTC %d : API token error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, api.Token, testCase.OutputAPI.Token)
		case api.URL != testCase.OutputAPI.URL:
			t.Errorf("\n\nTC %d : API token error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, api.URL, testCase.OutputAPI.URL)
		}
	}
}

func TestCheckCloudDcStatus(t *testing.T) {
	testCases := []struct {
		ID               int
		InputAPI         *API
		TCResourceTooler Resourceer
		Err              error
	}{
		{1,
			&API{
				wrongAPIToken,
				rightAPIURL,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongTokenError),
		},
		{2,
			&API{
				rightAPIToken,
				wrongAPIURL,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongAPIURLError),
		},
		{3,
			&API{
				wrongAPIToken,
				wrongAPIURL,
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongAPIURLError),
		},
		{4,
			&API{
				rightAPIToken,
				rightAPIURL,
				&http.Client{},
			},
			FakeResourceResourceer{},
			nil,
		},
	}
	fakeAPItools := APITooler{}
	fakeClientTooler := &ClientTooler{
		Client: HTTPClienter{},
	}
	fakeResourceTooler := &ResourceTooler{}
	for _, testCase := range testCases {
		fakeAPItools.APIImplementer = FakeAirDrumResourceAPIer{}
		fakeResourceTooler.Resource = testCase.TCResourceTooler
		err := fakeAPItools.CheckCloudDcStatus(testCase.InputAPI,
			fakeClientTooler,
			fakeResourceTooler)
		switch {
		case err == nil || testCase.Err == nil:
			if !(err == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : Check API error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err, testCase.Err)
			}
		case err.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : Check API error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.Err.Error())
		}
	}
}

func TestCreateResource(t *testing.T) {
	resourceName := "Unit test resource creation"
	testCases := []struct {
		ID              int
		TcClienter      Clienter
		ResourceType    string
		CreationErr     error
		CreatedResource map[string]interface{}
	}{
		{
			1,
			VMCreationSuccessHTTPClienterFake{},
			VMResourceType,
			nil,
			noTemplateVMMap,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			3,
			ResourceCreationFailureHTTPClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(creationOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			4,
			HandleRespErrHTTPClienterFake{},
			VMResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
	}
	apier := AirDrumResourcesAPI{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplatesTooler := TemplatesTooler{
		TemplatesTools: TemplateTemplater{},
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
			&fakeTemplatesTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(testCase.CreatedResource, respCreationMap)
		switch {
		case err == nil || testCase.CreationErr == nil:
			if !(err == nil && testCase.CreationErr == nil) {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.CreationErr)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.CreationErr != nil:
			if err.Error() != testCase.CreationErr.Error() {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.CreationErr.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) \n%s",
				testCase.ID, diffs)
		}
	}
}

func TestReadResource(t *testing.T) {
	resourceName := "Unit test resource read"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		ResourceType string
		ReadError    error
		ReadResource map[string]interface{}
	}{
		{
			1,
			VMReadSuccessHTTPClienterFake{},
			VMResourceType,
			nil,
			noTemplateVMMap,
		},
		{
			2,
			VdcReadSuccessHTTPClienterFake{},
			VdcResourceType,
			nil,
			vdcReadResponseMap,
		},
		{
			3,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			4,
			ResourceReadFailureHTTPClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(readOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			5,
			Error404HTTPClienterFake{},
			VdcResourceType,
			ErrResourceNotExist,
			map[string]interface{}{},
		},
		{
			6,
			HandleRespErrHTTPClienterFake{},
			VMResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
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
		respCreationMap, err := APIImplementerer.ReadResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(respCreationMap, testCase.ReadResource)
		switch {
		case err == nil || testCase.ReadError == nil:
			if !((err == nil) && (testCase.ReadError == nil)) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.ReadError)
			} else {
				if diffs != "" {
					t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.ReadError != nil:
			if cmp.Diff(respCreationMap, map[string]interface{}{}) != "" {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.ID, respCreationMap, testCase.ReadResource)
			}
			if err.Error() != testCase.ReadError.Error() {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.ReadError.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
				testCase.ID, diffs)
		}
	}
}

func TestUpdateResource(t *testing.T) {
	resourceName := "Unit test resource update"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		ResourceType string
		UpdateErr    error
	}{
		{
			1,
			ResourceUpdateSuccessHTTPClienterFake,
			VMResourceType,
			nil,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceUpdateFailureHTTPClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(updateOperation,
				resourceName,
				errEmptyResp),
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplatesTooler := TemplatesTooler{
		TemplatesTools: TemplateTemplater{},
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
		err := APIImplementerer.UpdateResource(d,
			&fakeClientTooler,
			&fakeTemplatesTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		switch {
		case err == nil || testCase.UpdateErr == nil:
			if !(err == nil && testCase.UpdateErr == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.UpdateErr)
			}
		case err.Error() != testCase.UpdateErr.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.UpdateErr.Error())
		}
	}
}

func TestDeleteResource(t *testing.T) {
	resourceName := "Unit test resource deletion"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		ResourceType string
		DeleteErr    error
	}{
		{
			1,
			ResourceDeletionSuccessHTTPClienterFake,
			VMResourceType,
			nil,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceDeletionFailureHTTPClienterFake{},
			VdcResourceType,
			errDoCrudRequestsBuilder(deleteOperation,
				resourceName,
				errEmptyResp),
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
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
		err := APIImplementerer.DeleteResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		switch {
		case err == nil || testCase.DeleteErr == nil:
			if !(err == nil && testCase.DeleteErr == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.DeleteErr)
			}
		case err.Error() != testCase.DeleteErr.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.DeleteErr.Error())
		}
	}
}
