package forem_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArticleDataSource(t *testing.T) {
	articleID := os.Getenv("TEST_DATA_FOREM_ARTICLE_ID")
	resourceName := "data.forem_article.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccArticleDataSourceConfig_id(articleID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", articleID),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "body_markdown"),
					resource.TestCheckResourceAttr(resourceName, "type_of", "article"),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
				),
			},
		},
	})
}

func testAccArticleDataSourceConfig_id(id string) string {
	return fmt.Sprintf(`
data "forem_article" "test" {
	id = "%s"
}
`, id)
}
