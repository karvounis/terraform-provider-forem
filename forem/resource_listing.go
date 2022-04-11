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
	maxListingTags = 8
)

var (
	allowedListingCategories = []string{"cfp", "forhire", "collabs", "education", "jobs", "mentors", "products", "mentees", "forsale", "events", "misc"}
	allowedListingActions    = []string{"draft", "bump", "publish", "unpublish"}
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
				ValidateFunc: validation.StringInSlice(allowedListingCategories, false),
			},
			"tags": {
				Description: "List of tags related to the listing.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    maxListingTags,
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
				ValidateFunc: validation.StringInSlice(allowedListingActions, false),
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

	var lbSchema dev.ListingBodySchema
	lbSchema.Listing.Title = title
	lbSchema.Listing.BodyMarkdown = d.Get("body_markdown").(string)
	lbSchema.Listing.Category = d.Get("category").(dev.ListingCategory)

	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
		lbSchema.Listing.Tags = tags
	}

	if v, ok := d.GetOk("contact_via_connect"); ok {
		lbSchema.Listing.ContactViaConnect = v.(bool)
	}
	if v, ok := d.GetOk("expires_at"); ok {
		lbSchema.Listing.ExpiresAt = v.(string)
	}
	if v, ok := d.GetOk("location"); ok {
		lbSchema.Listing.Location = v.(string)
	}
	if v, ok := d.GetOk("action"); ok {
		lbSchema.Listing.Action = v.(dev.Action)
	}
	if v, ok := d.GetOk("organization_id"); ok {
		lbSchema.Listing.OrganizationID = v.(int64)
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating listing with title `%s`", title))
	resp, err := client.CreateListing(lbSchema, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Created listing ID: `%d`", resp.ID))

	d.SetId(strconv.Itoa(int(resp.ID)))

	return resourceListingRead(ctx, d, meta)
}

func resourceListingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceListingRead(ctx, d, meta)
}

func resourceListingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Getting listing with ID: %s", id))
	resp, err := client.GetListingByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Found listing with ID: %s", id))

	d.SetId(id)
	d.Set("title", resp.Title)
	d.Set("body_markdown", resp.BodyMarkdown)
	d.Set("slug", resp.Slug)
	d.Set("category", resp.Category)
	d.Set("published", resp.Published)
	d.Set("tags", resp.Tags)

	if resp.User != nil {
		d.Set("user", map[string]interface{}{
			"name":             resp.User.Name,
			"username":         resp.User.Username,
			"twitter_username": resp.User.TwitterUsername,
			"github_username":  resp.User.GithubUsername,
			"website_url":      resp.User.WebsiteURL,
			"profile_image":    resp.User.ProfileImage,
		})
	} else {
		d.Set("user", map[string]interface{}{})
	}

	if resp.Organization != nil {
		d.Set("organization", map[string]interface{}{
			"name":             resp.Organization.Name,
			"username":         resp.Organization.Username,
			"slug":             resp.Organization.Slug,
			"profile_image":    resp.Organization.ProfileImage,
			"profile_image_90": resp.Organization.ProfileImage90,
		})
	} else {
		d.Set("organization", map[string]interface{}{})
	}

	return nil
}
