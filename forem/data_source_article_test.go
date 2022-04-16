package forem_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArticleDataSource_queryByID(t *testing.T) {
	articleID := os.Getenv("TEST_DATA_FOREM_ARTICLE_ID")
	resourceName := "data.forem_article.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleDataSourceConfig_queryByID(articleID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", articleID),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "body_markdown"),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckNoResourceAttr(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckResourceAttrSet(resourceName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceName, "published_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "comments_count"),
					resource.TestCheckResourceAttrSet(resourceName, "positive_reactions_count"),
					resource.TestCheckResourceAttrSet(resourceName, "public_reactions_count"),
					resource.TestCheckResourceAttrSet(resourceName, "reading_time_minutes"),
				),
			},
		},
	})
}

func TestAccArticleDataSource_queryByUsernameAndSlug(t *testing.T) {
	articleUsername := os.Getenv("TEST_DATA_FOREM_ARTICLE_USERNAME")
	articleSlug := os.Getenv("TEST_DATA_FOREM_ARTICLE_SLUG")
	resourceName := "data.forem_article.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleDataSourceConfig_queryByUsernameAndSlug(articleUsername, articleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "body_markdown"),
					resource.TestCheckResourceAttr(resourceName, "slug", articleSlug),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckResourceAttrSet(resourceName, "published_at"),
					resource.TestCheckResourceAttrSet(resourceName, "published_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "comments_count"),
					resource.TestCheckResourceAttrSet(resourceName, "positive_reactions_count"),
					resource.TestCheckResourceAttrSet(resourceName, "public_reactions_count"),
					resource.TestCheckResourceAttrSet(resourceName, "reading_time_minutes"),
				),
			},
		},
	})
}

func testAccArticleDataSourceConfig_queryByID(id string) string {
	return fmt.Sprintf(`
data "forem_article" "test" {
	id = "%s"
}
`, id)
}

func testAccArticleDataSourceConfig_queryByUsernameAndSlug(username, slug string) string {
	return fmt.Sprintf(`
data "forem_article" "test" {
	username = "%s"
	slug = "%s"
}
`, username, slug)
}
