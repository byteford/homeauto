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

    - Now run `terraform init`
        - This should have run though and not given any errors, Though it isn't doing anything at the moment

5. Setting up the provider
    - Open up `./provider/homeauto/provider.go`