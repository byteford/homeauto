package homeauto

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider defines the data terraform uses to make the provider and resources
func Provider() *schema.Provider {
	return &schema.Provider{
	
	}
}
// providerConfigure is used to set up the Client object which is used when calling the API
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	var c = Client{}
	return c, diags
}
