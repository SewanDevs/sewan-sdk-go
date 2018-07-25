package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func TestFetchTemplateFromList(t *testing.T) {
	test_cases := []struct {
		Id            int
		Template_name string
		TemplateList  []interface{}
		Template      map[string]interface{}
		Error         error
	}{
		{
			1,
			"",
			[]interface{}{},
			map[string]interface{}{},
			errors.New("Template \"\" does not exists, please validate it's name."),
		},
		{
			2,
			"template2",
			TEMPLATES_LIST,
			TEMPLATE2_MAP,
			nil,
		},
		{
			3,
			"lastTemplate",
			TEMPLATES_LIST,
			LAST_TEMPLATE_IN_LIST,
			nil,
		},
		{
			4,
			"template 42",
			TEMPLATES_LIST,
			map[string]interface{}{},
			errors.New("Template \"template 42\" does not exists, please validate it's name."),
		},
		{
			5,
			"template 42",
			WRONG_TEMPLATES_LIST,
			map[string]interface{}{"intin":"toutou"},
			errors.New("One of the fetch template " +
				"has a wrong format." +
				"\ngot : string" +
				"\nwant : map"),
		},
	}
	var (
		err      error
		template map[string]interface{}
		diffs    string
	)
	fake_templates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, test_case := range test_cases {
		template, err = fake_templates_tooler.TemplatesTools.FetchTemplateFromList(test_case.Template_name,
			test_case.TemplateList)
		diffs = cmp.Diff(test_case.Template, template)
		switch {
		case err == nil || test_case.Error == nil:
			if !(err == nil && test_case.Error == nil) {
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err, test_case.Error)
			}
			if diffs != "" {
				t.Errorf("\n\nTC %d : Wrong FetchTemplateFromList() template"+
					" (-got +want) :\n%s", test_case.Id, diffs)
			}
		case err != nil && test_case.Error != nil:
			if err.Error() != test_case.Error.Error() {
				t.Errorf("\n\nTC %d : FetchTemplateFromList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err.Error(), test_case.Error.Error())
			}
		}
	}
}
func TestValidateTemplate(t *testing.T) {
	test_cases := []struct {
		Id       int
		Template map[string]interface{}
		Error    error
	}{
		{
			1,
			map[string]interface{}{},
			errors.New("Template missing fields : " + "\"" + NAME_FIELD + "\" " +
				"\"" + OS_FIELD + "\" " +
				"\"" + RAM_FIELD + "\" " +
				"\"" + CPU_FIELD + "\" " +
				"\"" + ENTERPRISE_FIELD + "\" " +
				"\"" + DISKS_FIELD + "\" "),
		},
		{
			2,
			VM_CREATION_FROM_TEMPLATE1_SCHEMA,
			nil,
		},
		{
			3,
			VM_CREATION_FROM_TEMPLATE_SCHEMA_PRE_CREATION_WRONG_NICS_INIT_MAP,
			errors.New("Template nics is not a list as required but a string"),
		},
	}
	var (
		err error
	)
	fake_templates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, test_case := range test_cases {
		err = fake_templates_tooler.TemplatesTools.ValidateTemplate(test_case.Template)
		switch {
		case err == nil || test_case.Error == nil:
			if !(err == nil && test_case.Error == nil) {
				t.Errorf("\n\nTC %d : ValidateTemplate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err, test_case.Error)
			}
		case err != nil && test_case.Error != nil:
			switch {
			case err.Error() != test_case.Error.Error():
				t.Errorf("\n\nTC %d : ValidateTemplate() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err.Error(), test_case.Error.Error())
			default:
			}
		}
	}
}

func TestUpdateSchemaFromTemplateOnResourceCreation(t *testing.T) {
	test_cases := []struct {
		Id        int
		Dinit     *schema.ResourceData
		Template  map[string]interface{}
		DfinalMap map[string]interface{}
		Error     error
	}{
		{
			1,
			vm_schema_init(NON_EXISTING_ERROR_VM_SCHEMA_MAP),
			TEMPLATE2_MAP,
			NON_EXISTING_ERROR_VM_SCHEMA_MAP,
			errors.New("Template field should not be set on " +
				"an existing resource, please review the configuration field." +
				"\n : The resource schema has not been updated."),
		},
		{
			2,
			vm_schema_init(VM_SCHEMA_MAP_PRE_UPDATE_FROM_TEMPLATE),
			LAST_TEMPLATE_IN_LIST,
			VM_SCHEMA_MAP_POST_UPDATE_FROM_TEMPLATE,
			nil,
		},
	}
	var (
		err          error
		test_val     interface{}
		expected_val interface{}
		diffs        string
	)
	fake_templates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, test_case := range test_cases {
		err = fake_templates_tooler.TemplatesTools.UpdateSchemaFromTemplateOnResourceCreation(test_case.Dinit,
			test_case.Template)
		diffs = cmp.Diff(test_val, expected_val)
		switch {
		case err == nil || test_case.Error == nil:
			if !(err == nil && test_case.Error == nil) {
				t.Errorf("\n\nTC %d : UpdateSchemaFromTemplateOnResourceCreation() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err, test_case.Error)
			}
			for fieldName, fieldValue := range test_case.DfinalMap {
				if fieldName == ID_FIELD {
					test_val = test_case.Dinit.Id()
				} else {
					test_val = test_case.Dinit.Get(fieldName)
				}
				expected_val = fieldValue
				if diffs != "" {
					t.Errorf("TC %d : Schema update error  (-got +want) :\n%s",
						test_case.Id, diffs)
				}
			}
		case err != nil && test_case.Error != nil:
			switch {
			case err.Error() != test_case.Error.Error():
				t.Errorf("\n\nTC %d : UpdateSchemaFromTemplateOnResourceCreation() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err.Error(), test_case.Error.Error())
			default:
			}
		}
	}
}

func TestCreateTemplateOverrideConfig(t *testing.T) {
	test_cases := []struct {
		Id                     int
		D                      *schema.ResourceData
		Template               map[string]interface{}
		OverrideFile           string
		OverrideFileDataStruct map[string]interface{}
		Error                  error
	}{
		{
			1,
			vm_schema_init(NO_TEMPLATE_VM_MAP),
			map[string]interface{}{},
			"",
			map[string]interface{}{},
			errors.New("Schema \"Template\" field is empty, " +
				"can not create a template override configuration."),
		},
		{
			2,
			vm_schema_init(VM_CREATION_FROM_TEMPLATE1_SCHEMA),
			TEMPLATE1_MAP_BIS,
			"template1_template_override.tf.json",
			RESOURCE_OVERRIDE_JSON_MAP,
			nil,
		},
	}
	var (
		err          error
		overrideFile string
		jsonDiffs    string
	)
	fake_templates_tooler := TemplatesTooler{
		TemplatesTools: Template_Templater{},
	}
	for _, test_case := range test_cases {
		err,
			overrideFile = fake_templates_tooler.TemplatesTools.CreateTemplateOverrideConfig(test_case.D,
			test_case.Template)
		switch {
		case overrideFile != test_case.OverrideFile:
			t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() created overrideFile"+
				" error."+
				"\n\rcreatedFile: \"%s\"\n\rexpected: \"%s\"",
				test_case.Id, overrideFile, test_case.OverrideFile)
		case err == nil || test_case.Error == nil:
			if !(err == nil && test_case.Error == nil) {
				t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err, test_case.Error)
			} else {
				err, jsonDiffs = CompareJsonAndMap(overrideFile,
					test_case.OverrideFileDataStruct)
				if err != nil {
					t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() "+
						" json file and test data struct failed."+
						"\n\rJson file error : \"%s",
						test_case.Id, err.Error())
				}
				if jsonDiffs != "" {
					t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() generated"+
						" json file is incorrect,"+
						"\n\rDiffs : \"%s",
						test_case.Id, jsonDiffs)
				}
			}
		case err != nil && test_case.Error != nil:
			switch {
			case err.Error() != test_case.Error.Error():
				t.Errorf("\n\nTC %d : CreateTemplateOverrideConfig() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					test_case.Id, err.Error(), test_case.Error.Error())
			default:
			}
		}
	}
}
