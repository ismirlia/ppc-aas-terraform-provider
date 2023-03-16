// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	//"fmt"
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCPublicNetwork() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPPCPublicNetworkRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPPCPublicNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	networkC := instance.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.GetAllPublic()
	if err != nil {
		return diag.FromErr(err)
	}
	if len(networkdata.Networks) < 1 {
		return diag.Errorf("error getting public network or no public network found in %s", cloudInstanceID)
	}

	d.SetId(*networkdata.Networks[0].NetworkID)
	if networkdata.Networks[0].Type != nil {
		d.Set("type", networkdata.Networks[0].Type)
	}
	if networkdata.Networks[0].Name != nil {
		d.Set("name", networkdata.Networks[0].Name)
	}
	if networkdata.Networks[0].VlanID != nil {
		d.Set("vlan_id", networkdata.Networks[0].VlanID)
	}

	return nil
}
