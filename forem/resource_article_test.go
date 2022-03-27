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

func TestAccArticle_basic(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	resourceName := "forem_article.test"
	title := gofakeit.SentenceSimple()
	body_markdown := gofakeit.HipsterParagraph(2, 5, 10, "\n")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleBasicConfig(title, body_markdown),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", body_markdown),
					resource.TestCheckResourceAttr(resourceName, "type_of", "article"),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckNoResourceAttr(resourceName, "series"),
				),
			},
		},
	})
}

func TestAccArticle_full(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	title := gofakeit.HipsterSentence(5)
	published := gofakeit.Bool()
	canonical_url := gofakeit.URL()
	tags := []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()}
	series := gofakeit.LoremIpsumSentence(3)
	article := &dev.Article{
		Title:        title,
		BodyMarkdown: gofakeit.HipsterParagraph(2, 5, 10, "\n"),
		Published:    published,
		Description:  gofakeit.LoremIpsumSentence(5),
		CanonicalURL: canonical_url,
		CoverImage:   gofakeit.URL(),
		TagList:      tags,
	}

	articleUpd := &dev.Article{
		Title:        title,
		BodyMarkdown: gofakeit.HipsterParagraph(2, 5, 10, "\n"),
		Published:    published,
		Description:  gofakeit.LoremIpsumSentence(5),
		CanonicalURL: canonical_url,
		CoverImage:   gofakeit.URL(),
		TagList:      append(tags, gofakeit.Word()),
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
					resource.TestCheckResourceAttr(resourceName, "type_of", "article"),
					resource.TestCheckResourceAttr(resourceName, "published", strconv.FormatBool(published)),
					resource.TestCheckResourceAttr(resourceName, "description", article.Description),
					resource.TestCheckResourceAttr(resourceName, "main_image", article.CoverImage),
					resource.TestCheckResourceAttr(resourceName, "canonical_url", article.CanonicalURL),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(tags))),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttr(resourceName, "series", series),
				),
			},
			{
				Config: testAccArticleFullConfig(articleUpd, series),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", articleUpd.Title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", articleUpd.BodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "type_of", "article"),
					resource.TestCheckResourceAttr(resourceName, "published", strconv.FormatBool(published)),
					resource.TestCheckResourceAttr(resourceName, "description", articleUpd.Description),
					resource.TestCheckResourceAttr(resourceName, "main_image", articleUpd.CoverImage),
					resource.TestCheckResourceAttr(resourceName, "canonical_url", articleUpd.CanonicalURL),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(articleUpd.TagList))),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttr(resourceName, "series", series),
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

func testAccArticleBasicConfig(title, body_markdown string) string {
	return fmt.Sprintf(`
resource "forem_article" "test" {
	title         = %[1]q
	body_markdown = %[2]q
	}`, title, body_markdown)
}

func testAccArticleFullConfig(article *dev.Article, series string) string {
	return fmt.Sprintf(`
resource "forem_article" "test" {
	title         = %[1]q
	body_markdown = %[2]q
	published     = %v
	description   = %q
	canonical_url = %q
	series        = %q
	main_image    = %q
	tags          = %s
}`, article.Title, article.BodyMarkdown, article.Published, article.Description, article.CanonicalURL, series,
		article.CoverImage, strings.Join(strings.Split(fmt.Sprintf("%q\n", article.TagList), " "), ", "))
}
