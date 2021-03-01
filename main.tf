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
resource "homeauto_light" "main" {
    entity_id = "light.virtual_light_10"
    state     = "on"
    }