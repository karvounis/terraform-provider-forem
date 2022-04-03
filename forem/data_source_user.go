package forem

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "``forem_user` fetches information about a particular user. You can either use the user's ID or the user's username." +
			"\n\n## API Docs\n\n" +
			"https://developers.forem.com/api#operation/getUser",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "ID of the user. Please specify the `id` or the `username` of the desired user.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "username"},
			},
			"username": {
				Description:  "Username of the user. Please specify the `id` or the `username` of the desired user.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "username"},
			},
			"name": {
				Description: "Name of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"summary": {
				Description: "Summary of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"twitter_username": {
				Description: "User's twitter username. Can be null.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"github_username": {
				Description: "User's github username. Can be null.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"website_url": {
				Description: "User's website URL. Can be null.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"location": {
				Description: "User's location. Can be null.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"joined_at": {
				Description: "Date of joining (formatted with strftime '%b %e, %Y').",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"profile_image": {
				Description: "Profile image (320x320).",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	var userResp *dev.User
	var err error
	if v, ok := d.GetOk("id"); ok {
		id := v.(string)
		tflog.Debug(ctx, fmt.Sprintf("Getting user with id: %s", id))
		userResp, err = client.GetUserByID(id)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Found user with id: %s", id))
	} else {
		username := d.Get("username").(string)
		tflog.Debug(ctx, fmt.Sprintf("Getting user with username: %s", username))
		userResp, err = client.GetUserByUsername(dev.UserQueryParams{URL: username})
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Found user with username: %s", username))
	}

	d.SetId(strconv.Itoa(int(userResp.ID)))
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
