package homeauto

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider defines the data terraform uses to make the provider and resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_HOST", nil),
				Description: "There URL of the server: eg. `http://127.0.0.1:8123`",
			},
			"bearer_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_BEARER_TOKEN", nil),
				Description: "There bearer Token of the server",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"homeauto_light": resourceLight(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure is used to set up the Client object which is used when calling the API
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	bearerToken := d.Get("bearer_token").(string)

	var host string
	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = tempHost
	}

	c := NewClient(host, bearerToken, &http.Client{})
	return c, diags
}
