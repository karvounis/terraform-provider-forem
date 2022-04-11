package forem_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/karvounis/dev-client-go"
)

func TestAccListing_draft(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	resourceName := "forem_listing.test"
	title := gofakeit.HipsterSentence(2)
	bodyMarkdown := gofakeit.HipsterParagraph(2, 3, 5, "\n")
	category := gofakeit.RandomString([]string{string(dev.ListingCategoryCfp), string(dev.ListingCategoryEvents), string(dev.ListingCategoryMisc)})
	tags := []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccListingDraft(title, bodyMarkdown, category, tags),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", bodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "category", category),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "contact_via_connect", "false"),
					resource.TestCheckResourceAttr(resourceName, "action", "draft"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(tags))),
					resource.TestCheckNoResourceAttr(resourceName, "expires_at"),
					resource.TestCheckNoResourceAttr(resourceName, "location"),
				),
			},
		},
	})
}

func testAccListingDraft(title, bodyMarkdown, category string, tags []string) string {
	return fmt.Sprintf(`
resource "forem_listing" "test" {
	title         = %q
	body_markdown = %q
	category      = %q
	action        = "draft"

	tags = %s
}`, title, bodyMarkdown, category, strings.Join(strings.Split(fmt.Sprintf("%q\n", tags), " "), ", "))
}
