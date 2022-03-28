package forem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

func dataSourceListing() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListingRead,
		Schema: map[string]*schema.Schema{
			"type_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body_markdown": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"processed_html": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tag_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"organization": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceListingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Getting listing: %s", id))
	listingResp, err := client.GetListingByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Found listing: %s", id))

	d.Set("type_of", listingResp.TypeOf)
	d.SetId(id)
	d.Set("title", listingResp.Title)
	d.Set("slug", listingResp.Slug)
	d.Set("body_markdown", listingResp.BodyMarkdown)
	d.Set("category", listingResp.Category)
	d.Set("processed_html", listingResp.ProcessedHTML)
	d.Set("published", listingResp.Published)
	d.Set("tag_list", listingResp.TagList)
	d.Set("tags", listingResp.Tags)

	if listingResp.User != nil {
		d.Set("user", map[string]interface{}{
			"name":             listingResp.User.Name,
			"username":         listingResp.User.Username,
			"twitter_username": listingResp.User.TwitterUsername,
			"github_username":  listingResp.User.GithubUsername,
			"website_url":      listingResp.User.WebsiteURL,
			"profile_image":    listingResp.User.ProfileImage,
		})
	}

	if listingResp.Organization != nil {
		d.Set("organization", map[string]interface{}{
			"name":             listingResp.Organization.Name,
			"username":         listingResp.Organization.Username,
			"slug":             listingResp.Organization.Slug,
			"profile_image":    listingResp.Organization.ProfileImage,
			"profile_image_90": listingResp.Organization.ProfileImage90,
		})
	}

	return nil
}
