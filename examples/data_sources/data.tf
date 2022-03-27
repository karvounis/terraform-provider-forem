data "forem_article" "example" {
  id = var.article_id
}

data "forem_user" "example_id" {
  id = var.user_id
}

data "forem_user" "example_username" {
  username = var.user_username
}

data "forem_listing" "example" {
  id = var.listing_id
}

data "forem_followed_tags" "example" {}
