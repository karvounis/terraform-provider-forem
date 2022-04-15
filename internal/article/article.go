package article

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/karvounis/dev-client-go"
)

// GetArticleBodySchemaFromResourceData generates a dev.ArticleBodySchema object based on the values in the schema.ResourceData object
func GetArticleBodySchemaFromResourceData(d *schema.ResourceData) dev.ArticleBodySchema {
	var abc dev.ArticleBodySchema
	abc.Article.Title = d.Get("title").(string)
	abc.Article.BodyMarkdown = d.Get("body_markdown").(string)

	if v, ok := d.GetOk("published"); ok {
		abc.Article.Published = v.(bool)
	}
	if v, ok := d.GetOk("series"); ok {
		abc.Article.Series = v.(string)
	}
	if v, ok := d.GetOk("cover_image"); ok {
		abc.Article.MainImage = v.(string)
	}
	if v, ok := d.GetOk("canonical_url"); ok {
		abc.Article.CanonicalURL = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		abc.Article.Description = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
		abc.Article.Tags = tags
	}
	if v, ok := d.GetOk("organization_id"); ok {
		abc.Article.OrganizationID = v.(int32)
	}
	return abc
}
