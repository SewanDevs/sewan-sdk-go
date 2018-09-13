package sewansdk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func unitTestVDCSchema(resourceName string) *schema.ResourceData {
	resourceResponse := resource(VdcResourceType)
	d := resourceResponse.TestResourceData()
	d.SetId("UnitTest vdc")
	d.Set(NameField, resourceName)
	d.Set(DataCenterField, rightDatacenter)
	return d
}

func unitTestVDCWrongDataCenterSchema(resourceName string) *schema.ResourceData {
	resourceResponse := resource(VdcResourceType)
	d := resourceResponse.TestResourceData()
	d.SetId("UnitTest vdc")
	d.Set(NameField, resourceName)
	d.Set(DataCenterField, wrongDatacenter)
	return d
}

func unitTestVMSchema(resourceName string) *schema.ResourceData {
	resourceResponse := resource(VMResourceType)
	d := resourceResponse.TestResourceData()
	d.SetId("UnitTest vm")
	d.Set(NameField, resourceName)
	return d
}
