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
		Description: "`forem_listing` data source that fetches information about a specific listing.\n\n" +
			"## API Docs\n\nhttps://developers.forem.com/api#operation/getListingById",
		ReadContext: dataSourceListingRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the listing.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"title": {
				Description: "Title of the listing.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"slug": {
				Description: "Slug of the listing.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"body_markdown": {
				Description: "The body of the listing in Markdown format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"category": {
				Description: "Category of the listing.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"published": {
				Description: "Whether the listing is published or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"tags": {
				Description: "List of tags related to the listing.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user": {
				Description: "User that has created this listing.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization": {
				Description: "Organization related to this listing.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

	d.SetId(id)
	d.Set("title", listingResp.Title)
	d.Set("slug", listingResp.Slug)
	d.Set("body_markdown", listingResp.BodyMarkdown)
	d.Set("category", listingResp.Category)
	d.Set("published", listingResp.Published)
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
