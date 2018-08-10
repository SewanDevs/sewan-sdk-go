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
		Id         int
		VmMap      map[string]interface{}
		VmIdString string
	}{
		{
			1,
			testUpdateVmMap,
			"unit test vm",
		},
		{
			2,
			testUpdateVmMapFloatId,
			"121212.12",
		},
		{
			3,
			testUpdateVmMapIntId,
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
		d = CreateTestResourceSchema(testCase.VmIdString)
		schemaTooler.SchemaTools.UpdateLocalResourceState(testCase.VmMap,
			d,
			&schemaTooler)
		for key, value := range testCase.VmMap {
			diffs = cmp.Diff(d.Get(key), value)
			switch {
			case key != IdField:
				if diffs != "" {
					t.Errorf("\n\nTC %d : Update of %s field failed (-got +want) :\n%s",
						testCase.Id, key, diffs)
				}
			default:
				if d.Id() != testCase.VmIdString {
					t.Errorf("\n\nTC %d : Update of Id reserved field failed "+
						errTestResultDiffs,
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
					testCase.Id, key, diffs)
			}
		}
	}
}
