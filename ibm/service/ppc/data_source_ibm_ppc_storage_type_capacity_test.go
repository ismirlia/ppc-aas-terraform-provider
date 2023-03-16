// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCStorageTypeCapacityDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCStorageTypeCapacityDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_storage_type_capacity.type", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCStorageTypeCapacityDataSourceConfig() string {
	return fmt.Sprintf(`
		data "ibm_ppc_storage_type_capacity" "type" {
			ppc_cloud_instance_id = "%s"
			ppc_storage_type = "%s"
		}
	`, acc.Ppc_cloud_instance_id, acc.PpcStorageType)
}
