terraform {
  required_version = ">= 0.12"
}

provider "forem" {
  api_key = var.api_key
  host    = var.host
}
