terraform {
  required_providers {
    forem = {
      version = "0.2"
      source  = "hashicorp.com/edu/forem"
    }
  }
}

provider "forem" {}

output "fa" {
  value = data.forem_profile_image.karvounis
}

data "forem_profile_image" "karvounis" {
  username = "karvounis"
}
