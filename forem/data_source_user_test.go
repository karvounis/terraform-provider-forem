package forem_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource_username(t *testing.T) {
	username := "karvounis"
	dataSourceName := "data.forem_user.test"
	randID := strconv.Itoa(acctest.RandIntRange(0, 500000))
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
				Config: testAccUserDataSourceConfig_id(randID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "username"),
					resource.TestCheckResourceAttr(dataSourceName, "id", randID),
					resource.TestCheckResourceAttr(dataSourceName, "type_of", "user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "joined_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile_image"),
				),
			},
			{
				Config:      testAccUserDataSourceConfig_username(acctest.RandString(20)),
				ExpectError: regexp.MustCompile(`error getting user: not found: 404`),
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
