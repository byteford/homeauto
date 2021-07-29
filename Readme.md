# Steps

anywhere you see `URL` change to your instance url (or localhost if not remote)
anywhere you see `NAME` change to a username (no spaces)

## 0. Pull down the `start` branch which has the boiler plate code in it

open


- Terminal

`<url>/wetty` 
 run `cd ~/workdir/homeauto`


- IDE if dont want to use vim to edit files


`<url>:8000`

## 1. Set up the home-assistant instance

If using the Playgrounds infrastruture skip to step 3.

1. Open Terminal
2. `docker-compose up --detach` (--detach means we can still use the same terminal)
3. Go to `URL:8123`
4. Make account - not https so don't use an important password
   (location and name doesn't matter)
5. Click finish

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

`If not running on linux inside the build.sh linux_amd64 need to change to your system `

- run `sh build.sh NAME 0.0.1` remember to change `NAME` to your panda name
  - This with compile the provider and save it to go build -o ~/.terraform.d/plugins/github.com/NAME/homeauto/0.0.1/linux_amd64/terraform-provider-homeauto_v0.0.1
  - The script will then run terraform plan and terraform apply
  - We get errors about use not doing anything with the `bearer_token` variable but it works other than that
- Go in to `main.tf` and set with provider up remember to change `NAME`

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

- run `sh build.sh NAME 0.0.1` remember to change `NAME` to your panda name
  - If you want to clean up the formatting in any of the terraform files just run `terraform fmt`
  - This should have run though and not given any errors, Though it isn't doing anything at the moment

## 5. Setting up the provider

- Open up `./provider/homeauto/provider.go`
- We have a skeleton of the Code to save time
- At the top of the file add `"net/http"` under `"context"`
- Replace the code for the `Provider()` function with the following

```go
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
```

- schemas are how we layout the  data between go and terraform:

  - `Schema: map[string]*schema.Schema{` Is used to set the providers inputs
  - `ResourcesMap: map[string]*schema.Resource{` Is used to define what resources can be called
  - `ConfigureContextFunc: providerConfigure` says what function should be used to configure the provider
- Replace the `providerConfigure()` function with

```go
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

```

- This is now the Provider set up, but it wont do anything until we create a resource

## 6. The resource

- Return to the `main.tf` file
- At the bottom of the file lets make a basic light. Paste the following

```HCL
resource "homeauto_light" "main" {
    entity_id = "light.virtual_light_10"
    state     = "on"
}
```

- run `sh build.sh NAME 0.0.1` remember to change `NAME` to your panda name
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
- Replace the `resourceLightCreate()` function with:

```go
func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    item := LightItem{
        EntityID: d.Get("entity_id").(string),
        State:    d.Get("state").(string),
    }
    err := StartLight(item,*c)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(item.EntityID)
    return resourceLightRead(ctx, d, m)
}
```

- run `sh build.sh NAME 0.0.1` remember to change `NAME` to your panda name 
    - Will through an error, but the light will be created
- If you now go back to `URL:8123`, you should see a light appear

## 8. Saving the state

- Having a way of creating light but not saving anything about them isn't that useful
- still in `./provider/homeauto/resource_Light.go` lets add to the `resourceLightRead` function that will let terraform ask the api what states its in
- replace the `resourceLightRead()` function with

```go
func resourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    var diags diag.Diagnostics
    lightID := d.Id()

    light, err := GetLight(lightID, *c)
    if err != nil {
        return diag.FromErr(err)
    }
    if err := d.Set("state", light.State); err != nil {
        return diag.FromErr(err)
    }
    return diags
}

```

- Have a look in to the `terraform.tfstate` that was made when we did the `terraform apply` in the last step, the `status = tanted` meaning terraform doesn't trust the state
- Delete the `terraform.tfstate` and `terraform.tfstate.backup` files if they exsist
- run `sh build.sh NAME 0.0.1` remember to change `NAME` to your panda name
    -and look at the `terraform.tfstate` file again, you will see the status value isnt there this time because the state can be tracked now

## 9. Updating the light

- One of the key features of terraform is being able to change the infrastructure by changing the .tf file. Lets implement that now.
- Still in the `./provider/homeauto/resource_Light.go` file replace the `resourceLightUpdate` function with:

```go
func resourceLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    item := LightItem{
        EntityID: d.Get("entity_id").(string),
        State:    d.Get("state").(string),
    }

    err := StartLight(item, *c)
    if err != nil {
        return diag.FromErr(err)
    }
    return resourceLightRead(ctx, d, m)
}

```

- If you go back to the `main.tf` and change the state to "off" and run `sh build.sh NAME 0.0.1` in the terminal if should turn the light "off"

## 10. Cleaning it up

- The last part we need to implement is `terraform destroy`
- To do this the following code should be put in to `resourceLightDelete`:

```go
func resourceLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    var diags diag.Diagnostics
    lightID := d.Id()

    err := DeleteLight(lightID, *c)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId("")
    return diags
}
```
- You might have to delete `"fmt"` at the top of the file as we are no longer using the go model
- If you now run `sh build.sh NAME 0.0.1` to build the latest code
- Then run `terraform destroy` to destroy the light
- Check `URL:8123` and you should see the light no longer there

## 11. Extra fun

- If you have been looking around the tfstate and schemas you might have see there is information about the light we haven't used. Below it the code to get lights that can do more than just be on and off, replace all the code in `./provider/homeauto/resource_Light.go` with it. Have a read, have a play and ask if you have any questions:

```go
package homeauto

import (
    "context"

"   github.com/hashicorp/terraform-plugin-sdk/v2/diag"
"g  ithub.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceLight contains the definition of the recourse that terraform creates
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

//resourceLightCreate is run when terraform apply creates a new resource
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
    err := StartLight(item, *c)
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(item.EntityID)
return resourceLightRead(ctx, d, m)
}

//resourceLightRead is used to get the state of a light from the API
func resourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    var diags diag.Diagnostics
    lightID := d.Id()

    light, err := GetLight(lightID, *c)
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

//resourceLightUpdate is called when a light already exists but the state needs to change
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
    err := StartLight(item, *c)
    if err != nil {
        return diag.FromErr(err)
    }
    return resourceLightRead(ctx, d, m)
}

//resourceLightDelete is called when you run terraform delete or remove a resource from the state
func resourceLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    c := m.(*Client)
    var diags diag.Diagnostics
    lightID := d.Id()

    err := DeleteLight(lightID, *c)
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
