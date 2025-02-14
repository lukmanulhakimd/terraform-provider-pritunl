package provider

import (
	"context"

	"github.com/disc/terraform-provider-pritunl/internal/pritunl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_URL", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_PASSWORD", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_INSECURE", false),
			},
			"connection_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_CONNECTION_CHECK", true),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pritunl_organization": resourceOrganization(),
			"pritunl_server":       resourceServer(),
			"pritunl_user":         resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pritunl_host": dataSourceHost(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)
	connectionCheck := d.Get("connection_check").(bool)

	apiClient := pritunl.NewClient(url, username, password, insecure)

	if connectionCheck {
		// execute test api call to ensure that provided credentials are valid and pritunl api works
		err := apiClient.TestApiCall()
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return apiClient, nil
}
