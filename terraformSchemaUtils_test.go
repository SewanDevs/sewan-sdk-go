package sewan_go_sdk

import (
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

func TestDeleteTerraformResource(t *testing.T) {
	d := CreateTestResourceSchema("resource to delete")
	schemaTooler := SchemaTooler{
		SchemaTools: SchemaSchemaer{},
	}
	schemaTooler.SchemaTools.DeleteTerraformResource(d)
	if d.Id() != "" {
		t.Errorf("Deletion of unit test resource failed.")
	}
}

func TestUpdateLocalResourceState_AND_ReadElement(t *testing.T) {
	testCases := []struct {
		ID         int
		VMMap      map[string]interface{}
		VMIDString string
	}{
		{
			1,
			testUpdateVMMap,
			"unit test vm",
		},
		{
			2,
			testUpdateVMMapFloatID,
			"121212.12",
		},
		{
			3,
			testUpdateVMMapIntID,
			"1212",
		},
	}
	var (
		d     *schema.ResourceData
		diffs string
	)
	schemaTooler := SchemaTooler{
		SchemaTools: SchemaSchemaer{},
	}
	for _, testCase := range testCases {
		d = CreateTestResourceSchema(testCase.VMIDString)
		schemaTooler.SchemaTools.UpdateLocalResourceState(testCase.VMMap,
			d,
			&schemaTooler)
		for key, value := range testCase.VMMap {
			diffs = cmp.Diff(d.Get(key), value)
			switch {
			case key != IDField:
				if diffs != "" {
					t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
						testCase.ID, key, diffs)
				}
			default:
				if d.Id() != testCase.VMIDString {
					t.Errorf("\n\nTC %d : Update of ID reserved field failed "+
						errTestResultDiffs,
						testCase.ID, d.Id(), testCase.VMIDString)
				}
			}
		}
	}
}

func TestUpdateVdcResourcesNames(t *testing.T) {
	testCases := []struct {
		ID        int
		DInit     *schema.ResourceData
		DFinalMap map[string]interface{}
		Err       error
	}{
		{
			1,
			vdcSchemaInit(vdcResourcesNamesPreUpdateMap),
			vdcResourcesNamesUpdatedMap,
			nil,
		},
	}
	var (
		diffs string
	)
	schemaTooler := SchemaTooler{
		SchemaTools: SchemaSchemaer{},
	}
	for _, testCase := range testCases {
		schemaTooler.SchemaTools.UpdateVdcResourcesNames(testCase.DInit)
		for key, value := range testCase.DFinalMap {
			diffs = cmp.Diff(testCase.DInit.Get(key), value)
			if diffs != "" {
				t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
					testCase.ID, key, diffs)
			}
		}
	}
}
