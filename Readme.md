Steps:

1. Set up the home-assistant instance

    - open Terminal
    - `docker-compose up`
    - go to `URL:8123`
    - Make account - not https so don't use an important password 
    (location and name doesn't matter)

    - click finish

2. get an api key

    - click on name (bottom left)
    - Scroll to bottom of the page
    - Create a token under Long-Lived Access Tokens
    - Give it a name and click ok
    - Save the token as we will use it later
    - Click ok

3. Save api key
    - Rename terraform.tfvars.example to terraform.tfvars
    - In the file replace YOUR-TOKEN in bearer_token= "YOUR-TOKEN" with the token you just made (If you have lost your token do step 2 again)

4. Connect terraform to the provider
    - run `sh build.sh 0.0.1`
        - This with compile the provider and save it to go build -o ~/.terraform.d/plugins/github.com/byteford/homeauto/0.0.1/darwin_amd64/
        - The script will then run terraform plan and terraform apply
        - We get errors about use not doing anything with the `bearer_token` variable but it works other than that
    - Go in to main.tf and set with provider up

        ```HCL
        terraform {
            required_providers {
                homeauto = {
                    version = "0.0.1"
                    source  = "github.com/byteford/homeauto"
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

5. Setting up the provider
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

    - schemas are how we layout the  data between go and terraform
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
