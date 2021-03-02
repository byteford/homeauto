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