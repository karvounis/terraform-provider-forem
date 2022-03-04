package forem

import (
	"context"

	dev "github.com/Mayowa-Ojo/dev-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOREM_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			// "forem_articles": dataSourceArticles(),
			// "forem_article":  dataSourceArticle(),

			"forem_profile_image": dataSourceProfileImage(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	var diags diag.Diagnostics
	c, err := dev.NewClient(apiKey)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Forem client",
			Detail:   "Unable to use apiKey",
		})
		return nil, diags
	}

	return c, diags
}
