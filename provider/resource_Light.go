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
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"color_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"brightness": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"white_value": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"supported_features": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hs_color": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
			},
			"rgb_color": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"xy_color": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
			},
		},
	}
}
func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//var diags diag.Diagnostics
	c := m.(*Client)
	item := LightItem{
		EntityID:          d.Get("entity_id").(string),
		State:             d.Get("state").(string),
		Brightness:        d.Get("brightness").(string),
		HsColor:           d.Get("hs_color").(string),
		RgbColor:          d.Get("rgb_color").(string),
		XyColor:           d.Get("xy_color").(string),
		WhiteValue:        d.Get("white_value").(string),
		Name:              d.Get("friendly_name").(string),
		ColorMode:         d.Get("color_mode").(string),
		SupportedFeatures: d.Get("supported_features").(string),
	}
	if item.State == "" {
		item.State = "on"
	}
	if item.SupportedFeatures == nil {
		item.SupportedFeatures = 147
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
