package forem

import (
	"fmt"
	"log"
	"strconv"
	"time"

	dev "github.com/Mayowa-Ojo/dev-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	FormatIntBase = 10
)

func dataSourceFollowedTags() *schema.Resource {
	return &schema.Resource{
		Description: "Followed tags data source",
		Read:        dataSourceFollowedTagsRead,
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

func dataSourceFollowedTagsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	log.Print("[DEBUG] Getting followed tags")
	ftResp, err := client.GetFollowedTags()
	if err != nil {
		return fmt.Errorf("error getting followed tags: %w", err)
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
