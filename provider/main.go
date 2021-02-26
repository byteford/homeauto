package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	})
}

/*
// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_HOST", nil),
			},
			"beaverToken": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_BEAVER_TOKEN", nil),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to create homeauto client",
		Detail:   "Unable to authenticate user for authenticated home auto",
	})
	return nil, diags
	bearerToken := d.Get("beaverToken").(string)

	var host *string
	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	c, err := NewClient(host, &bearerToken)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create homeauto client",
			Detail:   "Unable to authenticate user for authenticated home auto",
		})
		return nil, diags
	}
	return c, diags
}

//NewClient -
func NewClient(host, token *string) (*Client, error) {
	if *token == "" {
		return nil, fmt.Errorf("no token")
	}
	c := Client{
		HTTPClient: &http.Client{},
		HostURL:    "",
	}

	if host != nil {
		c.HostURL = *host
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api", c.HostURL), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	_, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
*/
