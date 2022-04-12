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
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
				),
			},
		},
	})
}

func TestAccListing_publish(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	resourceName := "forem_listing.test"

	lbc := getListingBodySchemaToPublish()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccListingPublish(lbc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", lbc.Listing.Title),
					resource.TestCheckResourceAttr(resourceName, "body_markdown", lbc.Listing.BodyMarkdown),
					resource.TestCheckResourceAttr(resourceName, "category", string(lbc.Listing.Category)),
					resource.TestCheckResourceAttr(resourceName, "published", "true"),
					resource.TestCheckResourceAttr(resourceName, "contact_via_connect", strconv.FormatBool(lbc.Listing.ContactViaConnect)),
					resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(lbc.Listing.Tags))),
					resource.TestCheckResourceAttr(resourceName, "location", lbc.Listing.Location),
					resource.TestCheckResourceAttr(resourceName, "expires_at", lbc.Listing.ExpiresAt),
					resource.TestCheckResourceAttrSet(resourceName, "user.username"),
					resource.TestCheckNoResourceAttr(resourceName, "action"),
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

func testAccListingPublish(lbc dev.ListingBodySchema) string {
	return fmt.Sprintf(`
resource "forem_listing" "test" {
	title         		= %q
	body_markdown		= %q
	category      		= %q
	expires_at    		= %q
	contact_via_connect = %t
	location 			= %q

	tags = %s
}`,
		lbc.Listing.Title,
		lbc.Listing.BodyMarkdown,
		lbc.Listing.Category,
		lbc.Listing.ExpiresAt,
		lbc.Listing.ContactViaConnect,
		lbc.Listing.Location,
		strings.Join(strings.Split(fmt.Sprintf("%q\n", lbc.Listing.Tags), " "), ", "),
	)
}

func getListingBodySchemaToPublish() dev.ListingBodySchema {
	return dev.ListingBodySchema{
		Listing: struct {
			Title             string              `json:"title"`
			BodyMarkdown      string              `json:"body_markdown"`
			Category          dev.ListingCategory `json:"category"`
			Tags              []string            `json:"tags"`
			TagList           string              `json:"tag_list"`
			ExpiresAt         string              `json:"expires_at"`
			ContactViaConnect bool                `json:"contact_via_connect"`
			Location          string              `json:"location"`
			OrganizationID    int64               `json:"organization_id,omitempty"`
			Action            dev.Action          `json:"action"`
		}{
			Title:             gofakeit.Sentence(5),
			BodyMarkdown:      gofakeit.Paragraph(1, 2, 5, "\n"),
			Category:          dev.ListingCategory(gofakeit.RandomString([]string{string(dev.ListingCategoryCfp), string(dev.ListingCategoryEvents), string(dev.ListingCategoryMisc)})),
			Tags:              []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
			ExpiresAt:         gofakeit.DateRange(time.Now(), time.Now().AddDate(0, 0, gofakeit.IntRange(1, 25))).Format("2006-01-02"),
			ContactViaConnect: gofakeit.Bool(),
			Location:          gofakeit.City(),
		},
	}
}
