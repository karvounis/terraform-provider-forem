package forem

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

const (
	DEV_TO_BASE_URL = "https://dev.to/api"
)

// Provider initialises the Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOREM_API_KEY", nil),
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOREM_HOST", DEV_TO_BASE_URL),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"forem_article": resourceArticle(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"forem_user":          dataSourceUser(),
			"forem_followed_tags": dataSourceFollowedTags(),
			"forem_listing":       dataSourceListing(),
			"forem_article":       dataSourceArticle(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	host := d.Get("host").(string)

	var diags diag.Diagnostics
	c, err := dev.NewClient(dev.Options{Token: apiKey, Host: host})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Forem client",
			Detail:   fmt.Sprintf("Unable to use apiKey for host `%s`", host),
		})
		return nil, diags
	}
	return c, diags
}
