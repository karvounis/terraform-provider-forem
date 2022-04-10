package forem

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dev "github.com/karvounis/dev-client-go"
)

const (
	MaxListingTags = 8
)

var (
	AllowedListingCategories = []string{"cfp", "forhire", "collabs", "education", "jobs", "mentors", "products", "mentees", "forsale", "events", "misc"}
	AllowedListingActions    = []string{"draft", "bump", "publish", "unpublish"}
)

func resourceListing() *schema.Resource {
	return &schema.Resource{
		Description: "`forem_listing` resource creates and updates a particular listing. A listing is a classified ad that users create on Forem. They can be related to conference announcements, job offers, mentorships, upcoming events and more." +
			"\n\n## API Docs\n\n" +
			"- https://developers.forem.com/api#operation/createListing\n" +
			"- https://developers.forem.com/api#operation/updateListing",
		ReadContext:   resourceListingRead,
		CreateContext: resourceListingCreate,
		UpdateContext: resourceListingUpdate,
		DeleteContext: resourceListingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the listing.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "Title of the listing.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"body_markdown": {
				Description: "The body of the listing in Markdown format.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"category": {
				Description:  "The category that the listing belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(AllowedListingCategories, false),
			},
			"tags": {
				Description: "List of tags related to the listing.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    MaxListingTags,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"contact_via_connect": {
				Description: "True if users are allowed to contact the listing's owner via DEV connect, false otherwise.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"expires_at": {
				Description: "Date and time of expiration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location": {
				Description: "Geographical area or city for the listing.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"action": {
				Description:  "Set it to `draft` to create an unpublished listing.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(AllowedListingActions, false),
			},
			"organization_id": {
				Description: "The id of the organization the user is creating the listing for. Only users belonging to an organization can assign the listing to it.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"slug": {
				Description: "Slug of the listing.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"published": {
				Description: "Whether the listing is published or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "When the listing was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user": {
				Description: "User object of the listing.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization": {
				Description: "Organization object of the listing.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceListingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceListingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	title := d.Get("title").(string)

	var ab dev.ListingBodySchema
	ab.Listing.Title = title
	ab.Listing.BodyMarkdown = d.Get("body_markdown").(string)

	if v, ok := d.GetOk("tags"); ok {
		tags := v.([]interface{})
		tagsList := []string{}
		for _, t := range tags {
			tagsList = append(tagsList, t.(string))
		}
		ab.Listing.Tags = tagsList
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating listing with title `%s`", title))

	resp, err := client.CreateListing(ab, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Created listing ID: `%d`", resp.ID))

	d.SetId(strconv.Itoa(int(resp.ID)))

	return resourceListingRead(ctx, d, meta)
}

func resourceListingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	if d.HasChanges("title", "body_markdown", "published", "series", "cover_image", "description", "tags", "organization_id", "canonical_url") {
		var ab dev.ListingBodySchema
		ab.Listing.Title = d.Get("title").(string)
		ab.Listing.BodyMarkdown = d.Get("body_markdown").(string)

		if v, ok := d.GetOk("tags"); ok {
			tags := v.([]interface{})
			tagsList := []string{}
			for _, t := range tags {
				tagsList = append(tagsList, t.(string))
			}
			ab.Listing.Tags = tagsList
		}

		_, err := client.UpdateListing(d.Id(), ab, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceListingRead(ctx, d, meta)
}

func resourceListingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// client := meta.(*dev.Client)

	// id := d.Get("id").(string)
	// tflog.Debug(ctx, fmt.Sprintf("Getting listing: %s", id))

	// d.SetId(id)
	// d.Set("title", listing.Title)
	// d.Set("description", listing.Description)
	// d.Set("body_markdown", listing.BodyMarkdown)
	// d.Set("slug", listing.Slug)
	// d.Set("path", listing.Path)
	// d.Set("url", listing.URL)
	// d.Set("canonical_url", listing.CanonicalURL)
	// d.Set("cover_image", listing.CoverImage)

	// d.Set("published", listing.Published)
	// d.Set("published_at", listing.PublishedAt)
	// d.Set("published_timestamp", listing.PublishedTimestamp)

	// d.Set("comments_count", listing.CommentsCount)
	// d.Set("positive_reactions_count", listing.PositiveReactionsCount)
	// d.Set("public_reactions_count", listing.PublicReactionsCount)
	// d.Set("reading_time_minutes", listing.ReadingTimeMinutes)
	// d.Set("page_views_count", listing.PageViewsCount)

	// d.Set("tags", listing.TagList)

	// if listing.User != nil {
	// 	d.Set("user", map[string]interface{}{
	// 		"name":             listing.User.Name,
	// 		"username":         listing.User.Username,
	// 		"twitter_username": listing.User.TwitterUsername,
	// 		"github_username":  listing.User.GithubUsername,
	// 		"website_url":      listing.User.WebsiteURL,
	// 		"profile_image":    listing.User.ProfileImage,
	// 	})
	// } else {
	// 	d.Set("user", map[string]interface{}{})
	// }

	// if listing.Organization != nil {
	// 	d.Set("organization", map[string]interface{}{
	// 		"name":             listing.Organization.Name,
	// 		"username":         listing.Organization.Username,
	// 		"slug":             listing.Organization.Slug,
	// 		"profile_image":    listing.Organization.ProfileImage,
	// 		"profile_image_90": listing.Organization.ProfileImage90,
	// 	})
	// } else {
	// 	d.Set("organization", map[string]interface{}{})
	// }

	return nil
}
