// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCKeyDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCKeyDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_key.testacc_ds_key", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCKeyDataSourceConfig() string {
	return fmt.Sprintf(`
data "ibm_ppc_key" "testacc_ds_key" {
    ppc_key_name = "%s"
    ppc_cloud_instance_id = "%s"
}`, acc.Ppc_key_name, acc.Ppc_cloud_instance_id)

}