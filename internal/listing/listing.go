package listing

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karvounis/dev-client-go"
)

// GetListingBodySchemaFromResourceData generates a dev.ListingBodySchema object based on the values in the schema.ResourceData object
func GetListingBodySchemaFromResourceData(d *schema.ResourceData) dev.ListingBodySchema {
	var lbc dev.ListingBodySchema
	lbc.Listing.Title = d.Get("title").(string)
	lbc.Listing.BodyMarkdown = d.Get("body_markdown").(string)
	lbc.Listing.Category = dev.ListingCategory(d.Get("category").(string))
	if v, ok := d.GetOk("expires_at"); ok {
		lbc.Listing.ExpiresAt = v.(string)
	}
	if v, ok := d.GetOk("contact_via_connect"); ok {
		lbc.Listing.ContactViaConnect = v.(bool)
	}
	if v, ok := d.GetOk("location"); ok {
		lbc.Listing.Location = v.(string)
	}
	if v, ok := d.GetOk("action"); ok {
		lbc.Listing.Action = dev.Action(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
		lbc.Listing.Tags = tags
	}
	return lbc
}
