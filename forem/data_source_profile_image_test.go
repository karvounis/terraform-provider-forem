package forem_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProfileImageDataSource(t *testing.T) {
	username := "karvounis"
	dataSourceName := "data.forem_profile_image.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileImageDataSourceConfig_basic(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
					resource.TestCheckResourceAttr(dataSourceName, "id", username),
					resource.TestCheckResourceAttr(dataSourceName, "type_of", "profile_image"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_of"),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile_image"),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile_image_90"),
					testCheckResourceAttrExistsInArray(dataSourceName, "image_of", []string{"user", "organization"}),
				),
			},
		},
	})
}

// testCheckResourceAttrExistsInArray is a TestCheckFunc which ensures that
// the value in state for the given name/key combination exists in the provided array.
func testCheckResourceAttrExistsInArray(name, key string, arr []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := modulePrimaryInstanceState(s.RootModule(), name)
		if err != nil {
			return err
		}
		ak := is.Attributes[key]
		for _, v := range arr {
			if v == ak {
				return nil
			}
		}
		return fmt.Errorf("key '%s' does not exist in the provided array '%#v'", key, arr)
	}
}

// modulePrimaryInstanceState returns the instance state for the given resource
// name in a ModuleState
func modulePrimaryInstanceState(ms *terraform.ModuleState, name string) (*terraform.InstanceState, error) {
	rs, ok := ms.Resources[name]
	if !ok {
		return nil, fmt.Errorf("Not found: %s in %s", name, ms.Path)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("No primary instance: %s in %s", name, ms.Path)
	}

	return is, nil
}

func testAccProfileImageDataSourceConfig_basic(username string) string {
	return fmt.Sprintf(`
data "forem_profile_image" "test" {
	username = "%s"
}
`, username)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FOREM_API_KEY"); v == "" {
		t.Fatal("FOREM_API_KEY must be set for acceptance tests")
	}
}
