package forem_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	username := os.Getenv("TEST_DATA_FOREM_USER_USERNAME")
	userID := os.Getenv("TEST_DATA_FOREM_USER_ID")
	dataSourceName := "data.forem_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_username(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "type_of", "user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "joined_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile_image"),
				),
			},
			{
				Config: testAccUserDataSourceConfig_id(userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "username"),
					resource.TestCheckResourceAttr(dataSourceName, "id", userID),
					resource.TestCheckResourceAttr(dataSourceName, "type_of", "user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "joined_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile_image"),
				),
			},
			{
				Config:      testAccUserDataSourceConfig_username(acctest.RandString(20)),
				ExpectError: regexp.MustCompile(`Error: not found: 404`),
			},
		},
	})
}

func testAccUserDataSourceConfig_username(username string) string {
	return fmt.Sprintf(`
data "forem_user" "test" {
	username = "%s"
}
`, username)
}

func testAccUserDataSourceConfig_id(id string) string {
	return fmt.Sprintf(`
data "forem_user" "test" {
	id = "%s"
}
`, id)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FOREM_API_KEY"); v == "" {
		t.Fatal("FOREM_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("FOREM_HOST"); v == "" {
		t.Fatal("FOREM_HOST must be set for acceptance tests")
	}
}
