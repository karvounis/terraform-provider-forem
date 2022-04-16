package forem_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/karvounis/dev-client-go"
)

func TestAccArticle_createDraft(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	resourceName := "forem_article.test"
	title := gofakeit.SentenceSimple()
	bodyMarkdown := gofakeit.HipsterParagraph(2, 5, 10, "\n")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleBasicConfig(title, bodyMarkdown),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", bodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "canonical_url"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckNoResourceAttr(resourceName, "series"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckNoResourceAttr(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckResourceAttr(resourceName, "published_at", ""),
					resource.TestCheckResourceAttr(resourceName, "published_timestamp", ""),
					resource.TestCheckResourceAttr(resourceName, "comments_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "positive_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "public_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "page_views_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttrSet(resourceName, "reading_time_minutes"),
				),
			},
		},
	})
}

func TestAccArticle_publishAndEdit(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	tags := []string{strings.ToLower(gofakeit.Word()), strings.ToLower(gofakeit.Word()), strings.ToLower(gofakeit.Word())}
	series := gofakeit.LoremIpsumSentence(3)
	published := true

	article := &dev.Article{
		Title:        gofakeit.HipsterSentence(5),
		BodyMarkdown: gofakeit.HipsterParagraph(2, 5, 10, "\n"),
		Published:    published,
		Description:  gofakeit.LoremIpsumSentence(3),
		CanonicalURL: gofakeit.URL(),
		CoverImage:   gofakeit.URL(),
		TagList:      tags,
	}

	articleUpd := &dev.Article{
		Title:        gofakeit.HipsterSentence(5),
		BodyMarkdown: gofakeit.HipsterParagraph(2, 5, 10, "\n"),
		Published:    published,
		Description:  gofakeit.LoremIpsumSentence(3),
		CanonicalURL: gofakeit.URL(),
		CoverImage:   gofakeit.URL(),
		TagList:      append(tags, strings.ToLower(gofakeit.Word())),
	}

	resourceName := "forem_article.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleFullConfig(article, series),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", article.Title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", article.BodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "published", strconv.FormatBool(published)),
					resource.TestCheckResourceAttr(resourceName, "description", article.Description),
					resource.TestCheckResourceAttr(resourceName, "cover_image", article.CoverImage),
					resource.TestCheckResourceAttr(resourceName, "canonical_url", article.CanonicalURL),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(tags))),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttr(resourceName, "series", series),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckNoResourceAttr(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckResourceAttrSet(resourceName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceName, "published_timestamp"),
					resource.TestCheckResourceAttr(resourceName, "comments_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "positive_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "public_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "page_views_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttrSet(resourceName, "reading_time_minutes"),
				),
			},
			{
				Config: testAccArticleFullConfig(articleUpd, series),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", articleUpd.Title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", articleUpd.BodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "published", strconv.FormatBool(published)),
					resource.TestCheckResourceAttr(resourceName, "description", articleUpd.Description),
					resource.TestCheckResourceAttr(resourceName, "cover_image", articleUpd.CoverImage),
					resource.TestCheckResourceAttr(resourceName, "canonical_url", articleUpd.CanonicalURL),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(articleUpd.TagList))),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttr(resourceName, "series", series),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckResourceAttrSet(resourceName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceName, "published_timestamp"),
					resource.TestCheckResourceAttr(resourceName, "comments_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "positive_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "public_reactions_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttr(resourceName, "page_views_count", strconv.Itoa(0)),
					resource.TestCheckResourceAttrSet(resourceName, "reading_time_minutes"),
				),
			},
		},
	})
}

func TestAccArticle_tooManyTags(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	article := &dev.Article{
		Title:        gofakeit.SentenceSimple(),
		BodyMarkdown: gofakeit.HipsterParagraph(2, 5, 10, "\n"),
		Published:    gofakeit.Bool(),
		Description:  gofakeit.LoremIpsumSentence(5),
		CanonicalURL: gofakeit.URL(),
		CoverImage:   gofakeit.URL(),
		TagList:      []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccArticleFullConfig(article, gofakeit.SentenceSimple()),
				ExpectError: regexp.MustCompile(`Error: Too many list items`),
			},
		},
	})
}

func testAccArticleBasicConfig(title, bodyMarkdown string) string {
	return fmt.Sprintf(`
resource "forem_article" "test" {
	title         = %q
	body_markdown = %q
}`, title, bodyMarkdown)
}

func testAccArticleFullConfig(article *dev.Article, series string) string {
	return fmt.Sprintf(`
resource "forem_article" "test" {
	title         = %q
	body_markdown = %q
	published     = %v
	description   = %q
	canonical_url = %q
	series        = %q
	cover_image   = %q
	tags          = %s
}`, article.Title, article.BodyMarkdown, article.Published, article.Description, article.CanonicalURL, series,
		article.CoverImage, strings.Join(strings.Split(fmt.Sprintf("%q\n", article.TagList), " "), ", "))
}
