package forem_test

import (
	"terraform-provider-forem/forem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

func init() {
	testAccProviders = map[string]*schema.Provider{
		"forem": forem.Provider(),
	}
}
