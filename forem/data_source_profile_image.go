package forem

import (
	"fmt"
	"terraform-provider-forem/internal"

	dev "github.com/Mayowa-Ojo/dev-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProfileImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProfileImageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_image_90": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProfileImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	username := d.Get("username").(string)
	internal.LogDebug(fmt.Sprintf("Getting profile image for username: %s", username))
	profImageResp, err := client.GetProfileImage(username)
	if err != nil {
		return fmt.Errorf("error getting profile image: %w", err)
	}
	internal.LogDebug(fmt.Sprintf("Found profile image for username: %s", username))

	d.Set("type_of", profImageResp.TypeOf)
	d.Set("image_of", profImageResp.ImageOf)
	d.Set("profile_image", profImageResp.ProfileImage)
	d.Set("profile_image_90", profImageResp.ProfileImage90)
	d.SetId(username)

	return nil
}
