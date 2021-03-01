terraform {
  required_providers {
    homeauto = {
      version = "0.3.1"
      source  = "dgp.com/byteford/homeauto"
    }
  }
}
provider "homeauto" {
  host         = "http://127.0.0.1:8123"
  beaver_token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiI4M2MyNDc0YTY3OWI0ZWRjOTQ2YjZjZTM4ZTlhZTNhNiIsImlhdCI6MTYxNDQxNzAxMSwiZXhwIjoxOTI5Nzc3MDExfQ.VSpG9ivML0bBIwG82j64ek2PbiuKhL7hTO4pJrLKSS8"
}
resource "homeauto_light" "main" {
  entity_id = "light.virtual_light_10"
  state     = "on"
}
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