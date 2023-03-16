// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCVolumeGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCVolumeGroupRead,
		Schema: map[string]*schema.Schema{
			PPCVolumeGroupID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID or Name of the volume group",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"volume_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume group name",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consistency_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_description_errors": vgStatusDescriptionErrors(),
		},
	}
}

func vgStatusDescriptionErrors() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"message": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"volume_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceIBMPPCVolumeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	vgClient := instance.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.Get(d.Get(PPCVolumeGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*vgData.ID)
	d.Set("status", vgData.Status)
	d.Set("volume_group_name", vgData.Name)
	d.Set("consistency_group_name", vgData.ConsistencyGroupName)
	d.Set("replication_status", vgData.ReplicationStatus)
	if vgData.StatusDescription != nil {
		d.Set("status_description_errors", flattenVolumeGroupStatusDescription(vgData.StatusDescription.Errors))
	}

	return nil
}

func flattenVolumeGroupStatusDescription(list []*models.StatusDescriptionError) (errors []map[string]interface{}) {
	if list != nil {
		errors := make([]map[string]interface{}, len(list))
		for i, data := range list {
			l := map[string]interface{}{
				"key":        data.Key,
				"message":    data.Message,
				"volume_ids": data.VolIDs,
			}

			errors[i] = l
		}
		return errors
	}
	return
}
