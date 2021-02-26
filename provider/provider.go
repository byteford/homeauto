package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_HOST", nil),
			},
			"beaver_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_BEAVER_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"homeauto_light": resourceLight(),
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
	bearerToken := d.Get("beaver_token").(string)

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
			Detail:   fmt.Sprintf("%v", err),
		})
		return nil, diags
	}
	return c, diags
}

//NewClient -
func NewClient(host, token *string) (*Client, error) {
	if token == nil {
		return nil, fmt.Errorf("no token")
	}
	c := Client{
		HTTPClient: &http.Client{},
		HostURL:    "",
	}

	if host != nil {
		c.HostURL = *host
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	_, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
