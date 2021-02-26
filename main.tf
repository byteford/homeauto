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
    beaver_token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIwZDczODhjMGM2OWU0OWRmYjg3NjI5YmIyMjRkYzEyNSIsImlhdCI6MTYxNDM0ODQ0OCwiZXhwIjoxOTI5NzA4NDQ4fQ.p7IQqfwuUN0_L8-saNhiIIo-MqE0bFq-kTKOkdiqYqg"
}
resource "homeauto_light" main{}