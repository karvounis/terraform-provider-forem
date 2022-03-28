package forem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dev "github.com/karvounis/dev-client-go"
)

func dataSourceArticle() *schema.Resource {
	return &schema.Resource{
		Description: "Article Data source.",
		ReadContext: dataSourceArticleRead,
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
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cover_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"readable_publish_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"social_image": {
				Type:     schema.TypeString,
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
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"canonical_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comments_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"positive_reactions_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_reactions_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edited_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crossposted_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_comment_at": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"published_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body_html": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body_markdown": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reading_time_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flare_tag": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArticleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)

	tflog.Debug(ctx, fmt.Sprintf("Getting article: %s", id))
	articlesResp, err := client.GetPublishedArticleByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Found article: %s", id))

	d.Set("type_of", articlesResp.Article.TypeOf)
	d.SetId(id)
	d.Set("title", articlesResp.Article.Title)
	d.Set("description", articlesResp.Article.Description)
	d.Set("cover_image", articlesResp.Article.CoverImage)
	d.Set("readable_publish_date", articlesResp.Article.ReadablePublishDate)
	d.Set("social_image", articlesResp.Article.SocialImage)
	d.Set("tag_list", articlesResp.TagList)
	d.Set("tags", articlesResp.Tags)
	d.Set("slug", articlesResp.Article.Slug)
	d.Set("path", articlesResp.Article.Path)
	d.Set("url", articlesResp.Article.URL)
	d.Set("canonical_url", articlesResp.Article.CanonicalURL)
	d.Set("comments_count", articlesResp.Article.CommentsCount)
	d.Set("positive_reactions_count", articlesResp.Article.PositiveReactionsCount)
	d.Set("public_reactions_count", articlesResp.Article.PublicReactionsCount)
	d.Set("created_at", articlesResp.Article.CreatedAt)
	d.Set("edited_at", articlesResp.Article.EditedAt)
	d.Set("crossposted_at", articlesResp.Article.CrosspostedAt)
	d.Set("published_at", articlesResp.Article.PublishedAt)
	d.Set("last_comment_at", articlesResp.Article.LastCommentAt)
	d.Set("published_timestamp", articlesResp.Article.PublishedTimestamp)
	d.Set("body_html", articlesResp.Article.BodyHTML)
	d.Set("body_markdown", articlesResp.Article.BodyMarkdown)
	d.Set("reading_time_minutes", articlesResp.Article.ReadingTimeMinutes)

	if articlesResp.User != nil {
		d.Set("user", map[string]interface{}{
			"name":             articlesResp.User.Name,
			"username":         articlesResp.User.Username,
			"twitter_username": articlesResp.User.TwitterUsername,
			"github_username":  articlesResp.User.GithubUsername,
			"website_url":      articlesResp.User.WebsiteURL,
			"profile_image":    articlesResp.User.ProfileImage,
		})
	}

	if articlesResp.Organization != nil {
		d.Set("user", map[string]interface{}{
			"name":             articlesResp.Organization.Name,
			"username":         articlesResp.Organization.Username,
			"slug":             articlesResp.Organization.Slug,
			"profile_image":    articlesResp.Organization.ProfileImage,
			"profile_image_90": articlesResp.Organization.ProfileImage90,
		})
	}

	if articlesResp.FlareTag != nil {
		d.Set("flare_tag", map[string]interface{}{
			"name":           articlesResp.FlareTag.Name,
			"bg_color_hex":   articlesResp.FlareTag.BGColorHEX,
			"text_color_hex": articlesResp.FlareTag.TextColorHEX,
		})
	}

	return nil
}
