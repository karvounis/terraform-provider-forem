package forem

import (
	"fmt"
	"log"
	"strconv"
	"terraform-provider-forem/internal"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dev "github.com/karvounis/dev-client-go"
)

const (
	MaxArticleTags = 4
)

func resourceArticle() *schema.Resource {
	return &schema.Resource{
		Read:   resourceArticleRead,
		Create: resourceArticleCreate,
		Update: resourceArticleUpdate,
		Delete: resourceArticleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type_of": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body_markdown": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The body of the listing in Markdown format.",
			},
			"tag_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Comma separated list of tags.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    MaxArticleTags,
				Description: "Tags related to the listing.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"published": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True to create a published article, false otherwise. Defaults to false",
			},
			"series": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Article series name. All articles belonging to the same series need to have the same name in this parameter.",
			},
			"main_image": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"canonical_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"cover_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"social_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"readable_publish_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crossposted_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edited_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_comment_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Only users belonging to an organization can assign the article to it",
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
			"reading_time_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"page_views_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceArticleDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArticleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	title := d.Get("title").(string)

	var ab dev.ArticleBodySchema
	ab.Article.Title = title
	ab.Article.BodyMarkdown = d.Get("body_markdown").(string)

	if v, ok := d.GetOk("published"); ok {
		ab.Article.Published = v.(bool)
	}
	if v, ok := d.GetOk("series"); ok {
		ab.Article.Series = v.(string)
	}
	if v, ok := d.GetOk("main_image"); ok {
		ab.Article.MainImage = v.(string)
	}
	if v, ok := d.GetOk("canonical_url"); ok {
		ab.Article.CanonicalURL = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		ab.Article.Description = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.([]interface{})
		tagsList := []string{}
		for _, t := range tags {
			tagsList = append(tagsList, t.(string))
		}
		ab.Article.Tags = tagsList
	}
	if v, ok := d.GetOk("organization_id"); ok {
		ab.Article.OrganizationID = v.(int32)
	}

	internal.LogDebug(fmt.Sprintf("Creating article with title `%s`", title))

	resp, err := client.CreateArticle(ab, nil)
	if err != nil {
		return fmt.Errorf("error creating listing: %w", err)
	}
	internal.LogDebug(fmt.Sprintf("Created article ID: `%d`", resp.ID))

	d.SetId(strconv.Itoa(int(resp.ID)))

	return resourceArticleRead(d, meta)
}

func resourceArticleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	if d.HasChanges("title", "body_markdown", "published", "series", "main_image", "canonical_url", "description", "tags", "organization_id") {
		var ab dev.ArticleBodySchema
		ab.Article.Title = d.Get("title").(string)
		ab.Article.BodyMarkdown = d.Get("body_markdown").(string)

		if v, ok := d.GetOk("published"); ok {
			ab.Article.Published = v.(bool)
		}
		if v, ok := d.GetOk("series"); ok {
			ab.Article.Series = v.(string)
		}
		if v, ok := d.GetOk("main_image"); ok {
			ab.Article.MainImage = v.(string)
		}
		if v, ok := d.GetOk("canonical_url"); ok {
			ab.Article.CanonicalURL = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			ab.Article.Description = v.(string)
		}
		if v, ok := d.GetOk("tags"); ok {
			tags := v.([]interface{})
			tagsList := []string{}
			for _, t := range tags {
				tagsList = append(tagsList, t.(string))
			}
			ab.Article.Tags = tagsList
		}
		if v, ok := d.GetOk("organization_id"); ok {
			ab.Article.OrganizationID = v.(int32)
		}
		log.Println(ab)
		_, err := client.UpdateArticle(d.Id(), ab, nil)
		if err != nil {
			return fmt.Errorf("error creating listing: %w", err)
		}
	}
	return resourceArticleRead(d, meta)
}

func resourceArticleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*dev.Client)

	id := d.Get("id").(string)
	internal.LogDebug(fmt.Sprintf("Getting article: %s", id))

	page := int32(1)
	perPage := int32(1)
	missing := true
	var article dev.Article

	for missing {
		internal.LogDebug(fmt.Sprintf("Looking for article: %s with page: %d and perPage: %d", id, page, perPage))
		articleResp, err := client.GetUserArticles(dev.ArticleQueryParams{Page: page, PerPage: perPage})
		if err != nil {
			return fmt.Errorf("error getting article: %w", err)
		}
		log.Println(articleResp)

		if len(articleResp) == 0 {
			return fmt.Errorf("no more articles")
		}
		for _, v := range articleResp {
			if strconv.Itoa(int(v.ID)) == id {
				internal.LogDebug(fmt.Sprintf("Found article: %s", id))
				missing = false
				article = v
				break
			}
		}
		page++
	}

	d.Set("type_of", article.TypeOf)
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

	d.Set("tags", article.Tags)
	d.Set("tag_list", article.TagList)

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
