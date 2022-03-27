terraform {
  required_version = ">= 0.12"

  required_providers {
    forem = {
      version = "0.2"
      source  = "github.com/karvounis/forem"
    }
  }
}

provider "forem" {
  api_key = var.api_key
  host    = var.host
}
