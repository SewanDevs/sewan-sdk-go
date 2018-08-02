package sewan_go_sdk

import (
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func TestDeleteTerraformResource(t *testing.T) {
	d := CreateTestResourceSchema("resource to delete")
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
		Id         int
		VmMap      map[string]interface{}
		VmIdString string
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
		d = CreateTestResourceSchema(testCase.VmIdString)
		schemaTooler.SchemaTools.UpdateLocalResourceState(testCase.VmMap,
			d,
			&schemaTooler)
		for key, value := range testCase.VmMap {
			diffs = cmp.Diff(d.Get(key), value)
			switch {
			case key != ID_FIELD:
				if diffs != "" {
					t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
						testCase.Id, key, diffs)
				}
			default:
				if d.Id() != testCase.VmIdString {
					t.Errorf("\n\nTC %d : Update of Id reserved field failed "+
						ERROR_TEST_RESULT_DIFFS,
						testCase.Id, d.Id(), testCase.VmIdString)
				}
			}
		}
	}
}

func TestUpdateVdcResourcesNames(t *testing.T) {
	testCases := []struct {
		Id        int
		DInit     *schema.ResourceData
		DFinalMap map[string]interface{}
		Err       error
	}{
		{
			1,
			vdcSchemaInit(VDC_RESOURCES_NAMES_PRE_UPDATE_MAP),
			VDC_RESOURCES_NAMES_UPDATED_MAP,
			nil,
		},
	}
	var (
		diffs string
	)
	schemaTooler := SchemaTooler{
		SchemaTools: Schema_Schemaer{},
	}
	for _, testCase := range testCases {
		schemaTooler.SchemaTools.UpdateVdcResourcesNames(testCase.DInit)
		for key, value := range testCase.DFinalMap {
			diffs = cmp.Diff(testCase.DInit.Get(key), value)
			if diffs != "" {
				t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
					testCase.Id, key, diffs)
			}
		}
	}
}
