locals {
  tags = ["golang", "terraform", "forem", "md", "listings"]
}

# Minimum required values set
resource "forem_listing" "example_basic" {
  title         = "This a basic listing"
  body_markdown = "A simple single-line body that is very basic even for a listing..."
}

# Listing generated from a file
resource "forem_listing" "example_file" {
  title         = "How to create a Listing from a file!"
  body_markdown = file("${path.module}/files/example.md")

  tags = local.tags
}

# Full Listing example
resource "forem_listing" "example_full" {
  title               = "My first listing using TF Forem provider!"
  action              = "publish"
  category            = "cfp"
  expires_at          = "09/06/2025"
  contact_via_connect = false
  location            = "Amsterdam"

  body_markdown = <<-EOT
    This is the markdown for the listing
    # First header
    - bullet1
    - bullet2
    - bullet3
    - bullet4
EOT

  tags = local.tags
}
