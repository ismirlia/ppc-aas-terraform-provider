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

func DataSourceIBMPPCNetwork() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPPCNetworkRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCNetworkName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Network Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"cidr": {
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
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ip_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used_ip_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used_ip_percent": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"name": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "This value is deprecated in favor of" + helpers.PPCNetworkName,
			},
			"dns": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"jumbo": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPPCNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	networkC := instance.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.Get(d.Get(helpers.PPCNetworkName).(string))
	if err != nil || networkdata == nil {
		return diag.FromErr(err)
	}

	d.SetId(*networkdata.NetworkID)
	if networkdata.Cidr != nil {
		d.Set("cidr", networkdata.Cidr)
	}
	if networkdata.Type != nil {
		d.Set("type", networkdata.Type)
	}
	d.Set("gateway", networkdata.Gateway)
	if networkdata.VlanID != nil {
		d.Set("vlan_id", networkdata.VlanID)
	}
	if networkdata.IPAddressMetrics.Available != nil {
		d.Set("available_ip_count", networkdata.IPAddressMetrics.Available)
	}
	if networkdata.IPAddressMetrics.Used != nil {
		d.Set("used_ip_count", networkdata.IPAddressMetrics.Used)
	}
	if networkdata.IPAddressMetrics.Utilization != nil {
		d.Set("used_ip_percent", networkdata.IPAddressMetrics.Utilization)
	}
	if networkdata.Name != nil {
		d.Set("name", networkdata.Name)
	}
	if len(networkdata.DNSServers) > 0 {
		d.Set("dns", networkdata.DNSServers)
	}
	d.Set("jumbo", networkdata.Jumbo)

	return nil

}
