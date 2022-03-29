locals {
  series = "My TF provider - articles examples"
  tags   = ["golang", "terraform", "forem", "md"]
}

resource "forem_article" "example_basic" {
  title         = "This a basic article"
  body_markdown = "A simple single-line body that is very basic..."
  published     = true
  series        = local.series

  tags = local.tags
}

resource "forem_article" "example_full" {
  title         = "My first article using TF Forem provider!"
  published     = true
  cover_image   = "https://www.hdnicewallpapers.com/Walls/Big/Rainbow/Rainbow_on_Mountain_HD_Image.jpg"
  description   = "That is a description!"
  canonical_url = "https://github.com/karvounis"
  series        = local.series

  body_markdown = <<-EOT
    This is the markdown for the article
    # First header
    - bullet1
    - bullet2
    - bullet3
    - bullet4

    ## Second header
    ```yaml
    key: value
    key2: value2
    key3: value3
    ```

    | Fruit | price |
    | --- | --- |
    | banana | 2 |
    | watermelon | 4 |
    | apple | 1.2 |
  EOT

  tags = local.tags
}


resource "forem_article" "example_file" {
  title         = "Article from a file!"
  body_markdown = file("${path.module}/files/example.md")
}
