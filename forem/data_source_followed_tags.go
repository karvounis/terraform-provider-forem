package forem

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

const (
	FormatIntBase = 10
)

func dataSourceFollowedTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFollowedTagsRead,
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"points": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFollowedTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	tflog.Debug(ctx, "Getting followed tags")
	ftResp, err := client.GetFollowedTags()
	if err != nil {
		return diag.FromErr(err)
	}

	ftags := make([]interface{}, len(ftResp))
	for i, v := range ftResp {
		ft := make(map[string]interface{})

		ft["id"] = v.ID
		ft["name"] = v.Name
		ft["points"] = v.Points

		ftags[i] = ft
	}

	d.Set("tags", ftags)
	d.SetId(strconv.FormatInt(time.Now().Unix(), FormatIntBase))

	return nil
}
