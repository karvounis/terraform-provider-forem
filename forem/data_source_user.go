package forem

import (
	"fmt"
	"strconv"
	"terraform-provider-forem/internal"

	dev "github.com/Mayowa-Ojo/dev-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "User data source",
		Read:        dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"type_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "username"},
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "username"},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"summary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"twitter_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"github_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	var userResp *dev.User
	var err error
	if v, ok := d.GetOk("id"); ok {
		id := v.(string)
		internal.LogDebug(fmt.Sprintf("Getting user with id: %s", id))
		userResp, err = client.GetUserByID(id)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		internal.LogDebug(fmt.Sprintf("Found user with id: %s", id))
		d.Set("type_of", userResp.TypeOf)
	} else {
		username := d.Get("username").(string)
		internal.LogDebug(fmt.Sprintf("Getting user with username: %s", username))
		userResp, err = client.GetUserByUsername(dev.UserQueryParams{URL: username})
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		internal.LogDebug(fmt.Sprintf("Found user with username: %s", username))
	}

	d.SetId(strconv.Itoa(int(userResp.ID)))
	d.Set("type_of", userResp.TypeOf)
	d.Set("username", userResp.Username)
	d.Set("name", userResp.Name)
	d.Set("summary", userResp.Summary)
	d.Set("twitter_username", userResp.TwitterUsername)
	d.Set("github_username", userResp.GithubUsername)
	d.Set("website_url", userResp.WebsiteURL)
	d.Set("location", userResp.Location)
	d.Set("joined_at", userResp.JoinedAt)
	d.Set("profile_image", userResp.ProfileImage)

	return nil
}
