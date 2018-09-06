package sewansdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"testing"
)

func TestResourceInstanceCreate(t *testing.T) {
	testCases := []struct {
		ID           int
		D            *schema.ResourceData
		Clienter     Clienter
		Templater    Templater
		ResourceType string
		Error        error
		VMInstance   interface{}
	}{
		{
			1,
			vmSchemaInit(noTemplateVMMap),
			getListSuccessHTTPClienterFake{},
			TemplaterDummy{},
			VMResourceType,
			nil,
			vmInstanceNoTemplateVMMap(),
		},
		{
			2,
			vmSchemaInit(existingTemplateNoAdditionalDiskVMMap),
			getListSuccessHTTPClienterFake{},
			existingTemplateNoAdditionalDiskVMMapTemplaterFake{},
			VMResourceType,
			nil,
			fakeVMInstanceExistingTemplateNoAdditionalDiskVMMap(),
		},
		{
			3,
			vmSchemaInit(existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap),
			getListSuccessHTTPClienterFake{},
			existingTemplateWithAdditionalAndModifiedDisksAndNicsVMMapTemplaterFake{},
			VMResourceType,
			nil,
			fakeVMInstanceExistingTemplateWithAdditionalAndModifiedDisksAndNicsVMMap(),
		},
		{
			4,
			vmSchemaInit(nonExistingTemplateVMMap),
			getListSuccessHTTPClienterFake{},
			UnexistingTemplateTemplaterFake{},
			VMResourceType,
			errors.New("Unavailable template : windows95"),
			vmStruct{},
		},
		{
			5,
			vdcSchemaInit(vdcCreationMap),
			nil,
			TemplaterDummy{},
			VdcResourceType,
			nil,
			fakeVdcInstanceVdcCreationMap(),
		},
		{
			6,
			vdcSchemaInit(vdcCreationMap),
			getListSuccessHTTPClienterFake{},
			TemplaterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			nil,
		},
		{
			7,
			vmSchemaInit(nonExistingTemplateVMMap),
			getJSONListFailureHTTPClienterFake{},
			UnexistingTemplateTemplaterFake{},
			VMResourceType,
			errEmptyTemplateList,
			vmStruct{},
		},
		{
			8,
			vmSchemaInit(existingTemplateNoAdditionalDiskVMMap),
			getListSuccessHTTPClienterFake{},
			TemplateFormatErrorTemplaterFake{},
			VMResourceType,
			errors.New("Template missing fields : " + "\"" + NameField + "\" " +
				"\"" + OsField + "\" " +
				"\"" + RAMField + "\" " +
				"\"" + CPUField + "\" " +
				"\"" + EnterpriseField + "\" " +
				"\"" + DisksField + "\" " +
				"\"" + DatacenterField + "\" "),
			vmStruct{},
		},
		{
			9,
			vmSchemaInit(instanceNumberFieldUnitTestVMInstance),
			getListSuccessHTTPClienterFake{},
			instanceNumberFieldUnitTestVMInstanceMAPTemplaterFake{},
			VMResourceType,
			nil,
			fakeVMInstanceInstanceNumberFieldUnitTestVMInstanceMAP(),
		},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	fakeClientTooler := ClientTooler{}
	fakeTemplatesTooler := TemplatesTooler{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	for _, testCase := range testCases {
		fakeClientTooler.Client = testCase.Clienter
		fakeTemplatesTooler.TemplatesTools = testCase.Templater
		instance, err := fakeResourceTooler.Resource.resourceInstanceCreate(testCase.D,
			&fakeClientTooler,
			&fakeTemplatesTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(instance, testCase.VMInstance)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : resourceInstanceCreate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err, testCase.Error)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong resourceInstanceCreate() "+
						"created instance (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			switch {
			case diffs != "":
				t.Errorf("\n\nTC %d : Wrong resourceInstanceCreate() "+
					"created instance (-got +want) :\n%s",
					testCase.ID, diffs)
			case err.Error() != testCase.Error.Error():
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.Error.Error())
			}
		}
	}
}

func TestGetResourceCreationURLetResourceURL(t *testing.T) {
	testCases := []struct {
		ID    int
		api   API
		VMID  string
		VMURL string
	}{
		{1,
			API{
				rightAPIToken,
				rightAPIURL,
				enterpriseSlug,
				&http.Client{},
			},
			"42",
			rightVMURLQuaranteDeux,
		},
		{2,
			API{
				rightAPIToken,
				rightAPIURL,
				enterpriseSlug,
				&http.Client{},
			},
			"PATATE",
			rightVMURLPatate,
		},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		sVMURL := fakeResourceTooler.Resource.getResourceURL(&testCase.api,
			VMResourceType,
			testCase.VMID)
		switch {
		case sVMURL != testCase.VMURL:
			t.Errorf("VM url was incorrect,\n\rgot: \"%s\"\n\rwant: \"%s\"",
				sVMURL, testCase.VMURL)
		}
	}
}

func TestGetResourceCreationURL(t *testing.T) {
	testCases := []struct {
		ID                  int
		api                 API
		resourceCreationURL string
	}{
		{1,
			API{
				rightAPIToken,
				rightAPIURL,
				enterpriseSlug,
				&http.Client{},
			},
			rightVMCreationAPIURL,
		},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		sResourceCreationURL := fakeResourceTooler.Resource.getResourceCreationURL(&testCase.api,
			VMResourceType)
		switch {
		case sResourceCreationURL != testCase.resourceCreationURL:
			t.Errorf("resource api creation url was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				sResourceCreationURL, testCase.resourceCreationURL)
		}
	}
}

func TestValidateStatus(t *testing.T) {
	testCases := []struct {
		ID           int
		TcAPI        API
		Client       Clienter
		Err          error
		ResourceType string
	}{
		{1,
			API{
				rightAPIToken,
				rightAPIURL,
				enterpriseSlug,
				&http.Client{},
			},
			VMReadSuccessHTTPClienterFake{},
			nil,
			VMResourceType,
		},
		{2,
			API{
				rightAPIToken,
				rightAPIURL,
				enterpriseSlug,
				&http.Client{},
			},
			CheckRedirectReqFailureHTTPClienterFake{},
			errCheckRedirectFailure,
			VMResourceType,
		},
	}
	fakeResourceTooler := &ResourceTooler{
		Resource: ResourceResourceer{},
	}
	clientTooler := ClientTooler{}
	for _, testCase := range testCases {
		clientTooler.Client = testCase.Client
		apiClientErr := fakeResourceTooler.Resource.validateStatus(&testCase.TcAPI,
			testCase.ResourceType,
			clientTooler)
		switch {
		case apiClientErr == nil || testCase.Err == nil:
			if !(apiClientErr == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : validateStatus() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, apiClientErr, testCase.Err)
			}
		case apiClientErr.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : validateStatus() error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, apiClientErr.Error(), testCase.Err.Error())
		}
	}
}

func CreateTestResourceSchema(id interface{}) *schema.ResourceData {
	vmRes := resourceVM()
	d := vmRes.TestResourceData()
	d.SetId(id.(string))
	return d
}
