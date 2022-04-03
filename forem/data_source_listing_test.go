package forem_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccListingDataSource(t *testing.T) {
	listingID := os.Getenv("TEST_DATA_FOREM_LISTING_ID")
	resourceName := "data.forem_listing.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccListingDataSourceConfig_id(listingID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", listingID),
					resource.TestCheckResourceAttrSet(resourceName, "title"),
					resource.TestCheckResourceAttrSet(resourceName, "slug"),
					resource.TestCheckResourceAttrSet(resourceName, "body_markdown"),
					resource.TestCheckResourceAttrSet(resourceName, "category"),
				),
			},
		},
	})
}

func testAccListingDataSourceConfig_id(id string) string {
	return fmt.Sprintf(`
data "forem_listing" "test" {
	id = "%s"
}
`, id)
}
