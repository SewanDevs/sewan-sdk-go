package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func TestFetchTemplateFromList(t *testing.T) {
	testCases := []struct {
		Id           int
		TemplateName string
		TemplateList []interface{}
		Template     map[string]interface{}
		Error        error
	}{
		{
			1,
			"",
			[]interface{}{},
			map[string]interface{}(nil),
			errors.New("Template \"\" does not exists, please validate it's name."),
		},
		{
			2,
			"template2",
			templatesList,
			template2Map,
			nil,
		},
		{
			3,
			"lastTemplate",
			templatesList,
			lastTemplateInTemplatesList,
			nil,
		},
		{
			4,
			"template 42",
			templatesList,
			map[string]interface{}(nil),
			errors.New("Template \"template 42\" does not exists, please validate it's name."),
		},
		{
			5,
			"template 42",
			wrongTemplatesList,
			map[string]interface{}(nil),
			errors.New("One of the fetch template " +
				"has a wrong format." +
				"\ngot : string" +
				"\nwant : map"),
		},
	}
	var (
		err        error
		template   map[string]interface{}
		diffs      string
		errorError bool
	)
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, testCase := range testCases {
		template, err = fakeTemplates_tooler.TemplatesTools.FetchTemplateFromList(testCase.TemplateName,
			testCase.TemplateList)
		diffs = cmp.Diff(template, testCase.Template)
		switch {
		case err == nil || testCase.Error == nil:
			errorError = !(err == nil && testCase.Error == nil)
			switch {
			case errorError && (diffs != ""):
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"\n"+
					"\n\nAND Wrong FetchTemplateFromList() template"+
					" (-got +want) :\n%s",
					testCase.Id, err, testCase.Error,
					diffs)
			case errorError && (diffs == ""):
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Error)
			case !errorError && (diffs != ""):
				t.Errorf("\n\nTC %d : Wrong FetchTemplateFromList() template"+
					" (-got +want) :\n%s", testCase.Id, diffs)
			}
		case err != nil && testCase.Error != nil:
			errorError = err.Error() != testCase.Error.Error()
			switch {
			case errorError && (diffs != ""):
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\""+
					"AND Wrong FetchTemplateFromList() template (-got +want) :\n%s",
					testCase.Id, err.Error(), testCase.Error.Error(), diffs)
			case errorError && (diffs == ""):
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			case !errorError && (diffs != ""):
				t.Errorf("\n\nTC %d : Wrong FetchTemplateFromList() template"+
					" (-got +want) :\n%s", testCase.Id, diffs)
			}
		}
	}
}
func TestValidateTemplate(t *testing.T) {
	testCases := []struct {
		Id       int
		Template map[string]interface{}
		Error    error
	}{
		{
			1,
			map[string]interface{}{},
			errors.New("Template missing fields : " + "\"" + NameField + "\" " +
				"\"" + OsField + "\" " +
				"\"" + RamField + "\" " +
				"\"" + CpuField + "\" " +
				"\"" + EnterpriseField + "\" " +
				"\"" + DisksField + "\" "),
		},
		{
			2,
			vmCreationFromTemplate1Schema,
			nil,
		},
		{
			3,
			vmCreationFromTemplate1SchemaPreCreationWrongNicsInitMap,
			errors.New("Template nics is not a list as required but a string"),
		},
	}
	var (
		err error
	)
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, testCase := range testCases {
		err = fakeTemplates_tooler.TemplatesTools.ValidateTemplate(testCase.Template)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : ValidateTemplate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Error)
			}
		case err != nil && testCase.Error != nil:
			switch {
			case err.Error() != testCase.Error.Error():
				t.Errorf("\n\nTC %d : ValidateTemplate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			default:
			}
		}
	}
}

func TestUpdateSchemaFromTemplateOnResourceCreation(t *testing.T) {
	testCases := []struct {
		Id        int
		Dinit     *schema.ResourceData
		Template  map[string]interface{}
		DfinalMap map[string]interface{}
		Error     error
	}{
		{
			1,
			vmSchemaInit(nonExistingErrorVmSchemaMap),
			template2Map,
			nonExistingErrorVmSchemaMap,
			errors.New("Template field should not be set on " +
				"an existing resource, please review the configuration field." +
				"\n : The resource schema has not been updated."),
		},
		{
			2,
			vmSchemaInit(vmSchemaMapPreUpdateFromTemplate),
			lastTemplateInTemplatesList,
			vmSchemaMapPostUpdateFromTemplate,
			nil,
		},
	}
	var (
		err         error
		testVal     interface{}
		expectedVal interface{}
		diffs       string
	)
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, testCase := range testCases {
		err = fakeTemplates_tooler.TemplatesTools.UpdateSchemaFromTemplateOnResourceCreation(testCase.Dinit,
			testCase.Template)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : UpdateSchemaFromTemplateOnResourceCreation() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Error)
			}
			for fieldName, fieldValue := range testCase.DfinalMap {
				if fieldName == IdField {
					testVal = testCase.Dinit.Id()
				} else {
					testVal = testCase.Dinit.Get(fieldName)
				}
				expectedVal = fieldValue
				diffs = cmp.Diff(testVal, expectedVal)
				if diffs != "" {
					t.Errorf("TC %d : Schema update error  (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			switch {
			case err.Error() != testCase.Error.Error():
				t.Errorf("\n\nTC %d : UpdateSchemaFromTemplateOnResourceCreation() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			default:
			}
		}
	}
}

func TestCreateTemplateOverrideConfig(t *testing.T) {
	testCases := []struct {
		Id                     int
		D                      *schema.ResourceData
		Template               map[string]interface{}
		OverrideFile           string
		OverrideFileDataStruct map[string]interface{}
		Error                  error
	}{
		{
			1,
			vmSchemaInit(noTemplateVmMap),
			map[string]interface{}{},
			"",
			map[string]interface{}{},
			errors.New("Schema \"Template\" field is empty, " +
				"can not create a template override configuration."),
		},
		{
			2,
			vmSchemaInit(vmCreationFromTemplate1Schema),
			template1MapBis,
			"template1Template_override.tf.json",
			resourceOverrideJsonMap,
			nil,
		},
		{
			3,
			vmSchemaInit(vmCreationN42FromTemplate1Schema),
			template1MapBis,
			"template1Template_override.tf.json",
			resourceN42OverrideJsonMap,
			nil,
		},
	}
	fakeTemplates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, testCase := range testCases {
		overrideFile,
			err := fakeTemplates_tooler.TemplatesTools.CreateTemplateOverrideConfig(testCase.D,
			testCase.Template)
		switch {
		case overrideFile != testCase.OverrideFile:
			t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() created overrideFile"+
				" error."+
				"\n\rcreatedFile: \"%s\"\n\rexpected: \"%s\"",
				testCase.Id, overrideFile, testCase.OverrideFile)
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err, testCase.Error)
			} else {
				jsonDiffs, err2 := CompareJsonAndMap(overrideFile,
					testCase.OverrideFileDataStruct)
				if err2 != nil {
					t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() "+
						" json file and test data struct failed."+
						"\n\rJson file error : \"%s",
						testCase.Id, err2.Error())
				}
				if jsonDiffs != "" {
					t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() generated"+
						" json file is incorrect,"+
						"\n\rDiffs (-got +want) : \"%s",
						testCase.Id, jsonDiffs)
				}
			}
		case err != nil && testCase.Error != nil:
			switch {
			case err.Error() != testCase.Error.Error():
				t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			default:
			}
		}
	}
}
