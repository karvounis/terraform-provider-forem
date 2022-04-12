package forem

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dev "github.com/karvounis/dev-client-go"
)

const (
	devToBaseURL = "https://dev.to/api"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		if len(s.AtLeastOneOf) > 0 {
			desc += fmt.Sprintf(" At least one of the following has to be added: `%s`.", strings.Join(s.AtLeastOneOf, ", "))
		}
		if len(s.ConflictsWith) > 0 {
			desc += fmt.Sprintf(" Conflicts with the following: `%s`.", strings.Join(s.ConflictsWith, ", "))
		}
		if len(s.RequiredWith) > 0 {
			desc += fmt.Sprintf(" Required to be set with the following: `%s`.", strings.Join(s.RequiredWith, ", "))
		}
		if s.MinItems > 0 {
			desc += fmt.Sprintf(" Minimum items: `%d`.", s.MinItems)
		}
		if s.MaxItems > 0 {
			desc += fmt.Sprintf(" Maximum items: `%d`.", s.MaxItems)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider initialises the Forem Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_key": {
				Description: "API key to be able to communicate with the FOREM API. Can be specified with the `FOREM_API_KEY` environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOREM_API_KEY", nil),
			},
			"host": {
				Description: "Host of the FOREM API. You can specify the `dev.to` or any other Forem installation. Can be specified with the `FOREM_HOST` environment variable. Default: `https://dev.to/api`.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOREM_HOST", devToBaseURL),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"forem_article": resourceArticle(),
			"forem_listing": resourceListing(),
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
