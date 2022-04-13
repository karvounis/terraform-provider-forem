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
	allowedListingActions    = []string{string(dev.ActionDraft), string(dev.ActionBump), string(dev.ActionPublish), string(dev.ActionUnpublish)}
	allowedListingCategories = []string{
		string(dev.ListingCategoryCfp),
		string(dev.ListingCategoryForhire),
		string(dev.ListingCategoryCollabs),
		string(dev.ListingCategoryEducation),
		string(dev.ListingCategoryJobs),
		string(dev.ListingCategoryMentors),
		string(dev.ListingCategoryProducts),
		string(dev.ListingCategoryMentees),
		string(dev.ListingCategoryForsale),
		string(dev.ListingCategoryEvents),
		string(dev.ListingCategoryMisc),
	}
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
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization": {
				Description: "Organization object of the listing.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// TODO: Waiting for API to allow deletion of a listing
func resourceListingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceListingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	title := d.Get("title").(string)
	category := dev.ListingCategory(d.Get("category").(string))

	var lbSchema dev.ListingBodySchema
	lbSchema.Listing.Title = title
	lbSchema.Listing.BodyMarkdown = d.Get("body_markdown").(string)
	lbSchema.Listing.Category = category

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
		lbSchema.Listing.Action = dev.Action(v.(string))
	}
	if v, ok := d.GetOk("organization_id"); ok {
		lbSchema.Listing.OrganizationID = v.(int64)
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating listing with title `%s` and category `%s`", title, category))
	resp, err := client.CreateListing(lbSchema, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Created listing with ID: `%d`", resp.ID))

	d.SetId(strconv.Itoa(int(resp.ID)))

	return resourceListingRead(ctx, d, meta)
}

func resourceListingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	var lbc dev.ListingBodySchema

	if d.HasChange("title") {
		lbc.Listing.Title = d.Get("title").(string)
	}
	if d.HasChange("body_markdown") {
		lbc.Listing.BodyMarkdown = d.Get("body_markdown").(string)
	}
	if d.HasChange("category") {
		lbc.Listing.Category = dev.ListingCategory(d.Get("category").(string))
	}
	if d.HasChange("expires_at") {
		lbc.Listing.ExpiresAt = d.Get("expires_at").(string)
	}
	if d.HasChange("contact_via_connect") {
		lbc.Listing.ContactViaConnect = d.Get("contact_via_connect").(bool)
	}
	if d.HasChange("location") {
		lbc.Listing.Location = d.Get("location").(string)
	}
	if d.HasChange("action") {
		lbc.Listing.Action = dev.Action(d.Get("action").(string))
	}
	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := []string{}
			for _, t := range v.([]interface{}) {
				tags = append(tags, t.(string))
			}
			lbc.Listing.Tags = tags
		}
	}

	_, err := client.UpdateListing(d.Id(), lbc, nil)
	if err != nil {
		return diag.FromErr(err)
	}

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
