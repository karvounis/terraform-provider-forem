package forem_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArticle_basic(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	title := gofakeit.HipsterSentence(5)
	body_markdown := gofakeit.HipsterParagraph(2, 5, 10, "\n")
	resourceName := "forem_article.test"

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
				),
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
