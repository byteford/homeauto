package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLightCreate,
		ReadContext:   resourceLightRead,
		UpdateContext: resourceLightUpdate,
		DeleteContext: resourceLightDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//var diags diag.Diagnostics
	c := m.(*Client)
	item := LightItem{
		EntityID: d.Get("entity_id").(string),
		State:    d.Get("state").(string),
	}
	if item.State == "" {
		item.State = "on"
	}
	o, err := c.StartLight(item)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(o.EntityID)
	return resourceLightRead(ctx, d, m)
	//return diags
}
func resourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics
	lightID := d.Id()

	light, err := c.GetLight(lightID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", light.State); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func resourceLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	item := LightItem{
		EntityID: d.Get("entity_id").(string),
		State:    d.Get("state").(string),
	}
	_, err := c.StartLight(item)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceLightRead(ctx, d, m)
}
func resourceLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var diags diag.Diagnostics
	lightID := d.Id()

	err := c.DelLight(lightID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
