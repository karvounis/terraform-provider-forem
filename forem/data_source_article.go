package forem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

func dataSourceArticle() *schema.Resource {
	return &schema.Resource{
		Description: "`forem_article` data source fetches information about a particular article." +
			"\n\n## API Docs\n\n" +
			"https://developers.forem.com/api#operation/getArticleById",
		ReadContext: dataSourceArticleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the article.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"title": {
				Description: "Title of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Article description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cover_image": {
				Description: "URL of the cover image of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"social_image": {
				Description: "Social image of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tags": {
				Description: "List of tags related to the article.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"slug": {
				Description: "Slug of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"path": {
				Description: "Path of the article URL.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"url": {
				Description: "Full article URL.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"canonical_url": {
				Description: "Canonical URL of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"comments_count": {
				Description: "Number of comments.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"positive_reactions_count": {
				Description: "Number of positive reactions.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"public_reactions_count": {
				Description: "Number of public reactions.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_at": {
				Description: "When the article was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"edited_at": {
				Description: "When the article was edited.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"crossposted_at": {
				Description: "When the article was crossposted.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"published_at": {
				Description: "When the article was published.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_comment_at": {
				Description: "When the article was last commented.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"published_timestamp": {
				Description: "When the article was published.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"reading_time_minutes": {
				Description: "Article reading time in minutes.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"body_markdown": {
				Description: "The body of the article in Markdown format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user": {
				Description: "User object of the article.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization": {
				Description: "Organization object of the article.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flare_tag": {
				Description: "Flare tag object of the article.",
				Type:        schema.TypeMap,
				Computed:    true,
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

	d.SetId(id)
	d.Set("title", articlesResp.Article.Title)
	d.Set("description", articlesResp.Article.Description)
	d.Set("body_markdown", articlesResp.Article.BodyMarkdown)
	d.Set("reading_time_minutes", articlesResp.Article.ReadingTimeMinutes)

	d.Set("url", articlesResp.Article.URL)
	d.Set("canonical_url", articlesResp.Article.CanonicalURL)
	d.Set("cover_image", articlesResp.Article.CoverImage)
	d.Set("social_image", articlesResp.Article.SocialImage)
	d.Set("slug", articlesResp.Article.Slug)
	d.Set("path", articlesResp.Article.Path)
	d.Set("tags", articlesResp.Tags)

	d.Set("comments_count", articlesResp.Article.CommentsCount)
	d.Set("positive_reactions_count", articlesResp.Article.PositiveReactionsCount)
	d.Set("public_reactions_count", articlesResp.Article.PublicReactionsCount)

	d.Set("created_at", articlesResp.Article.CreatedAt)
	d.Set("edited_at", articlesResp.Article.EditedAt)
	d.Set("crossposted_at", articlesResp.Article.CrosspostedAt)
	d.Set("published_at", articlesResp.Article.PublishedAt)
	d.Set("last_comment_at", articlesResp.Article.LastCommentAt)
	d.Set("published_timestamp", articlesResp.Article.PublishedTimestamp)

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
