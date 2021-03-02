# Steps

anywhere you see `URL` change to your instance url (or localhost if not remote)
anywhere you see `NAME` change to a username (no spaces)

## 0. Pull down the `start` branch which has the boiler plate code in it

## 1. Set up the home-assistant instance

- Open Terminal
- `docker-compose up --detach` (--detach means we can still use the same terminal)
- Go to `URL:8123`
- Make account - not https so don't use an important password
(location and name doesn't matter)
- Click finish

## 2. get an api key

- Click on name (bottom left)
- Scroll to bottom of the page
- Create a token under Long-Lived Access Tokens
- Give it a name and click ok
- Save the token as we will use it later
- Click ok

## 3. Save api key

- Rename terraform.tfvars.example to terraform.tfvars
- In the file replace YOUR-TOKEN in bearer_token= "YOUR-TOKEN" with the token you just made (If you have lost your token do step 2 again)

## 4. Connect terraform to the provider

- run `sh build.sh NAME 0.0.1`
        - This with compile the provider and save it to go build -o ~/.terraform.d/plugins/github.com/byteford/homeauto/0.0.1/darwin_amd64/
        - The script will then run terraform plan and terraform apply
        - We get errors about use not doing anything with the `bearer_token` variable but it works other than that
- Go in to main.tf and set with provider up

```HCL
terraform {
    required_providers {
        homeauto = {
            version = "0.0.1"
            source  = "github.com/NAME/homeauto"
        }
    }
}
provider "homeauto" {
    host         = "http://127.0.0.1:8123"
    bearer_token = var.bearer_token
}  
```

- In `variables.tf` paste

```HCL
variable bearer_token{
    sensitive = true
}
```

- Now run `terraform init` and `terraform plan`
  - This should have run though and not given any errors, Though it isn't doing anything at the moment

## 5. Setting up the provider

- Open up `./provider/homeauto/provider.go`
- We have a skeliten of the Code to save time
- paste the following code in to the `return &schema.Provider block`

```go
Schema: map[string]*schema.Schema{
    "host": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_HOST", nil),
    },
    "bearer_token": &schema.Schema{
        Type:        schema.TypeString,
        Required:    true,
        Sensitive:   true,
        DefaultFunc: schema.EnvDefaultFunc("HOMEAUTO_BEARER_TOKEN", nil),
    },
},
ResourcesMap: map[string]*schema.Resource{
    "homeauto_light": resourceLight(),
},
ConfigureContextFunc: providerConfigure,
```

- schemas are how we layout the  data between go and terraform:
  - `Schema: map[string]*schema.Schema{` Is used to set the providers inputs
  - `ResourcesMap: map[string]*schema.Resource{` Is used to define what resources can be called
  - `ConfigureContextFunc: providerConfigure` says what function should be used to configure the provider

- In the `providerConfigure` replace its content with

```go
var diags diag.Diagnostics
bearerToken := d.Get("bearer_token").(string)

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

```

- This is now the Provider set up, but it wont do anything until we create a resource

## 6. The resource

- Return to the `main.tf` file
- Under the provider lets make a basic light. Paste the following

```HCL
resource "homeauto_light" "main" {
    entity_id = "light.virtual_light_10"
    state     = "on"
}
```

- run `sh build.sh NAME 0.0.1` to build and run everything again
  - You should get a lovely error message saying not set up, but if you scroll up you can see the plan trying to make the lights.
  - It knows what a light is as the schema is already made in `./provider/homeauto/resource_Light.go`

## 7. Getting creating

- Head to `./provider/homeauto/resource_Light.go`
  - You can see the schema has already been made but some things to point out

```go
    CreateContext: resourceLightCreate,
    ReadContext:   resourceLightRead,
    UpdateContext: resourceLightUpdate,
    DeleteContext: resourceLightDelete,
```

- This sets the methods to be run by terraform on different steps

- Replace the content of `resourceLightCreate` with:

```go
c := m.(*Client)
item := LightItem{
    EntityID: d.Get("entity_id").(string),
    State:    d.Get("state").(string),
}
o, err := c.StartLight(item)
if err != nil {
    return diag.FromErr(err)
}
d.SetId(o.EntityID)
return resourceLightRead(ctx, d, m)
```

- run `sh build.sh NAME 0.0.1`
- If you now go back to `URL:8123`, you should see a light appear

## 8. Saving the state

- Having a way of creating light but not saving anything about them its that useful
- still in `./provider/homeauto/resource_Light.go` lets add to the `resourceLightRead` function that will let terraform ask the api what states its in
- replace the body of the `resourceLightRead` function with

```go
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

```

- We still can't update the lights we make (if you get issues running terraform apply try deleting the `terraform.tfstate` file)

## 9. Updating the light

- One of the key features of terraform is being able to change the infrastructure by changing the .tf file. lets implement that now.
- still in the `./provider/homeauto/resource_Light.go` file add the following code to the `resourceLightUpdate` function:

```go
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

```

- If you go back to the `main.tf` and change the state to "off" and run `terraform apply` in the terminal if should turn the light "off"

## 10. cleaning it up

- The last part we need to implement is `terraform destroy`
- To do this the following code should be put in to `resourceLightDelete`:

```go
c := m.(*Client)
var diags diag.Diagnostics
lightID := d.Id()

err := c.DelLight(lightID)
if err != nil {
    return diag.FromErr(err)
}
d.SetId("")
return diags
```

- If you now run `sh build.sh NAME 0.0.1` to build the latest code
- Then run `terraform destroy` to destroy the light
- Check `URL:8123` and you should see the light no longer there

## 11. Extra fun

- If you have been looking around the tfstate and schemas you might have see there is information about the light we haven't used. Below it the code to get lights that can do more than just be on and off, replace all the code in `./provider/homeauto/resource_Light.go` with it. Have a read, have a play and ask if you have any questions:

```go
package homeauto

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
                Default:  "on",
            },
            "friendly_name": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Default:  "Light",
            },
            "color_mode": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Default:  "hs",
            },
            "brightness": &schema.Schema{
                Type:     schema.TypeInt,
                Optional: true,
                Default:  255,
            },
            "white_value": &schema.Schema{
                Type:     schema.TypeInt,
                Optional: true,
            },
            "supported_features": &schema.Schema{
                Type:     schema.TypeInt,
                Optional: true,
                Default:  147,
            },
            "hs_color": &schema.Schema{
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeFloat,
                },
            },
            "rgb_color": &schema.Schema{
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeInt,
                },
            },
            "xy_color": &schema.Schema{
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeFloat,
                },
            },
        },
    }
}
func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    item := LightItem{
        EntityID: d.Get("entity_id").(string),
        State:    d.Get("state").(string),
        Attr: Attributes{
            Brightness:        d.Get("brightness").(int),
            WhiteValue:        d.Get("white_value").(int),
            Name:              d.Get("friendly_name").(string),
            ColorMode:         d.Get("color_mode").(string),
            SupportedFeatures: d.Get("supported_features").(int),
        },
    }
    if d.Get("hs_color.#").(int) != 0 {
        item.Attr.HsColor = []float64{d.Get("hs_color.0").(float64), d.Get("hs_color.1").(float64)}
    }
    if d.Get("rgb_color.#").(int) != 0 {
        item.Attr.RgbColor = []int{d.Get("rgb_color.0").(int), d.Get("rgb_color.1").(int), d.Get("rgb_color.2").(int)}
    }
    if d.Get("xy_color.#").(int) != 0 {
        item.Attr.XyColor = []float64{d.Get("xy_color.0").(float64), d.Get("xy_color.1").(float64)}
    }
    o, err := c.StartLight(item)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(o.EntityID)
    return resourceLightRead(ctx, d, m)

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
    if err := d.Set("hs_color", light.Attr.HsColor); err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("rgb_color", light.Attr.RgbColor); err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("xy_color", light.Attr.XyColor); err != nil {
        return diag.FromErr(err)
    }
    return diags
}
func resourceLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    item := LightItem{
        EntityID: d.Get("entity_id").(string),
        State:    d.Get("state").(string),
        Attr: Attributes{
            Brightness:        d.Get("brightness").(int),
            WhiteValue:        d.Get("white_value").(int),
            Name:              d.Get("friendly_name").(string),
            ColorMode:         d.Get("color_mode").(string),
            SupportedFeatures: d.Get("supported_features").(int),
        },
    }
    if d.Get("hs_color.#").(int) != 0 {
        item.Attr.HsColor = []float64{d.Get("hs_color.0").(float64), d.Get("hs_color.1").(float64)}
    }
    if d.Get("rgb_color.#").(int) != 0 {
        item.Attr.RgbColor = []int{d.Get("rgb_color.0").(int), d.Get("rgb_color.1").(int), d.Get("rgb_color.2").(int)}
    }
    if d.Get("xy_color.#").(int) != 0 {
        item.Attr.XyColor = []float64{d.Get("xy_color.0").(float64), d.Get("xy_color.1").(float64)}
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
```

-A light that makes use of the new code

```hcl
resource "homeauto_light" "colour" {
    entity_id     = "light.virtual_light_12"
    state         = "on"
    brightness    = 100
    hs_color      = [300.0, 71.0]
    rgb_color     = [255, 72, 255]
    xy_color      = [0.38, 0.17]
    white_value   = 240
    friendly_name = "Light 14"
    color_mode    = "hs"
}

```
