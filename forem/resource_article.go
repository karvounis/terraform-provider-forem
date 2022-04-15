package forem

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-forem/internal/article"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dev "github.com/karvounis/dev-client-go"
)

const (
	maxArticleTags      = 4
	readArticlesPerPage = 25
)

func resourceArticle() *schema.Resource {
	return &schema.Resource{
		Description: "`forem_article` resource creates and updates a particular article." +
			"\n\n## API Docs\n\n" +
			"- https://developers.forem.com/api#operation/createArticle\n" +
			"- https://developers.forem.com/api#operation/updateArticle",
		ReadContext:   resourceArticleRead,
		CreateContext: resourceArticleCreate,
		UpdateContext: resourceArticleUpdate,
		DeleteContext: resourceArticleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "Title of the article.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "Slug of the article.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"body_markdown": {
				Description: "The body of the article in Markdown format.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tags": {
				Description: "List of tags related to the article.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    maxArticleTags,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"published": {
				Description: "Set to `true` to create a published article.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"series": {
				Description: "Article series name. All articles belonging to the same series need to have the same name in this parameter.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cover_image": {
				Description:  "URL of the cover image of the article.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"canonical_url": {
				Description:  "Canonical URL of the article.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"url": {
				Description: "Full article URL.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Article description.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"published_at": {
				Description: "When the article was published.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"published_timestamp": {
				Description: "When the article was published.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"path": {
				Description: "Path of the article URL.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_id": {
				Description: "Only users belonging to an organization can assign the article to it.",
				Type:        schema.TypeInt,
				Optional:    true,
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
			"reading_time_minutes": {
				Description: "Article reading time in minutes.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"page_views_count": {
				Description: "Number of views.",
				Type:        schema.TypeInt,
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
			"created_at": {
				Description: "When the listing was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "When the listing was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// TODO: Waiting for API to allow deletion of an article
func resourceArticleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceArticleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	abc := article.GetArticleBodySchemaFromResourceData(d)
	tflog.Debug(ctx, fmt.Sprintf("Creating article with title: `%s`", abc.Article.Title))

	resp, err := client.CreateArticle(abc, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Created article ID: `%d`", resp.ID))

	d.SetId(strconv.Itoa(int(resp.ID)))
	d.Set("created_at", time.Now().Format(time.RFC3339))

	return resourceArticleRead(ctx, d, meta)
}

func resourceArticleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	tflog.Debug(ctx, fmt.Sprintf("Updating article with ID: %s", d.Id()))
	abc := article.GetArticleBodySchemaFromResourceData(d)
	if _, err := client.UpdateArticle(d.Id(), abc, nil); err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Updated listing with ID: %s", d.Id()))

	d.Set("updated_at", time.Now().Format(time.RFC3339))

	return resourceArticleRead(ctx, d, meta)
}

func resourceArticleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Getting article with ID: %s", id))

	page := int32(1)
	perPage := int32(readArticlesPerPage)
	missing := true
	var article dev.Article

	for missing {
		tflog.Debug(ctx, fmt.Sprintf("Looking for article: %s with page: %d and perPage: %d", id, page, perPage))
		articleResp, err := client.GetUserArticles(dev.ArticleQueryParams{Page: page, PerPage: perPage})
		if err != nil {
			return diag.FromErr(err)
		}
		if len(articleResp) == 0 {
			return diag.Errorf("no more articles")
		}

		for _, v := range articleResp {
			if strconv.Itoa(int(v.ID)) == id {
				tflog.Debug(ctx, fmt.Sprintf("Found article with ID: %s", id))
				missing = false
				article = v
				break
			}
		}
		page++
	}

	d.SetId(id)
	d.Set("title", article.Title)
	d.Set("description", article.Description)
	d.Set("body_markdown", article.BodyMarkdown)
	d.Set("slug", article.Slug)
	d.Set("path", article.Path)
	d.Set("url", article.URL)
	d.Set("canonical_url", article.CanonicalURL)
	d.Set("cover_image", article.CoverImage)

	d.Set("published", article.Published)
	d.Set("published_at", article.PublishedAt)
	d.Set("published_timestamp", article.PublishedTimestamp)

	d.Set("comments_count", article.CommentsCount)
	d.Set("positive_reactions_count", article.PositiveReactionsCount)
	d.Set("public_reactions_count", article.PublicReactionsCount)
	d.Set("reading_time_minutes", article.ReadingTimeMinutes)
	d.Set("page_views_count", article.PageViewsCount)

	d.Set("tags", article.TagList)

	if article.User != nil {
		d.Set("user", map[string]interface{}{
			"name":             article.User.Name,
			"username":         article.User.Username,
			"twitter_username": article.User.TwitterUsername,
			"github_username":  article.User.GithubUsername,
			"website_url":      article.User.WebsiteURL,
			"profile_image":    article.User.ProfileImage,
		})
	} else {
		d.Set("user", map[string]interface{}{})
	}

	if article.Organization != nil {
		d.Set("organization", map[string]interface{}{
			"name":             article.Organization.Name,
			"username":         article.Organization.Username,
			"slug":             article.Organization.Slug,
			"profile_image":    article.Organization.ProfileImage,
			"profile_image_90": article.Organization.ProfileImage90,
		})
	} else {
		d.Set("organization", map[string]interface{}{})
	}

	if article.FlareTag != nil {
		d.Set("flare_tag", map[string]interface{}{
			"name":           article.FlareTag.Name,
			"bg_color_hex":   article.FlareTag.BGColorHEX,
			"text_color_hex": article.FlareTag.TextColorHEX,
		})
	} else {
		d.Set("flare_tag", map[string]interface{}{})
	}

	return nil
}
