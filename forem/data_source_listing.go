package forem

import (
	"fmt"
	"terraform-provider-forem/internal"

	dev "github.com/Mayowa-Ojo/dev-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceListing() *schema.Resource {
	return &schema.Resource{
		Description: "Listing data source",
		Read:        dataSourceListingRead,
		Schema: map[string]*schema.Schema{
			"type_of": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"listing"}, false),
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"cfp", "forhire", "collabs", "education", "jobs", "mentors", "products", "mentees", "forsale", "events", "misc"}, false),
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
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceListingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)
	internal.LogDebug(fmt.Sprintf("Getting listing: %s", id))
	listingResp, err := client.GetListingByID(id)
	if err != nil {
		return fmt.Errorf("error getting listing: %w", err)
	}
	internal.LogDebug(fmt.Sprintf("Found listing: %s", id))

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

	return nil
}
