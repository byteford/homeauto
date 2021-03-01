terraform {
  required_providers {
    homeauto = {
      version = "0.3.1"
      source  = "dgp.com/byteford/homeauto"
    }
  }
}
provider "homeauto"{
    host = "http://127.0.0.1:8123"
    beaver_token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiI4M2MyNDc0YTY3OWI0ZWRjOTQ2YjZjZTM4ZTlhZTNhNiIsImlhdCI6MTYxNDQxNzAxMSwiZXhwIjoxOTI5Nzc3MDExfQ.VSpG9ivML0bBIwG82j64ek2PbiuKhL7hTO4pJrLKSS8"
}
resource "homeauto_light" main{
  entity_id = "light.virtual_light_10"
  state = "off"
}