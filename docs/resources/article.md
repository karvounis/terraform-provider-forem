---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "forem_article Resource - terraform-provider-forem"
subcategory: ""
description: |-
  forem_article resource creates and updates a particular article.
  API Docs
  https://developers.forem.com/api#operation/createArticlehttps://developers.forem.com/api#operation/updateArticle
---

# forem_article (Resource)

`forem_article` resource creates and updates a particular article.

## API Docs

- https://developers.forem.com/api#operation/createArticle
- https://developers.forem.com/api#operation/updateArticle

## Example Usage

```terraform
locals {
  series = "My TF provider - articles examples"
  tags   = ["golang", "terraform", "forem", "md"]
}

# Minimum required values set
resource "forem_article" "example_basic" {
  title         = "This a basic article"
  body_markdown = "A simple single-line body that is very basic..."
}

# Article generated from a file
resource "forem_article" "example_file" {
  title         = "Article from a file!"
  body_markdown = file("${path.module}/files/example.md")
  series        = local.series

  tags = local.tags
}

# Full Article example
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `body_markdown` (String) The body of the article in Markdown format.
- `title` (String) Title of the article.

### Optional

- `canonical_url` (String) Canonical URL of the article.
- `cover_image` (String) URL of the cover image of the article.
- `description` (String) Article description.
- `organization_id` (Number) Only users belonging to an organization can assign the article to it.
- `published` (Boolean) Set to `true` to create a published article. Defaults to: `false`.
- `series` (String) Article series name. All articles belonging to the same series need to have the same name in this parameter.
- `tags` (List of String) List of tags related to the article. Maximum items: `4`.

### Read-Only

- `comments_count` (Number) Number of comments.
- `created_at` (String) When the listing was created.
- `flare_tag` (Map of String) Flare tag object of the article.
- `id` (String) ID of the article.
- `organization` (Map of String) Organization object of the article.
- `page_views_count` (Number) Number of views.
- `path` (String) Path of the article URL.
- `positive_reactions_count` (Number) Number of positive reactions.
- `public_reactions_count` (Number) Number of public reactions.
- `published_at` (String) When the article was published.
- `published_timestamp` (String) When the article was published.
- `reading_time_minutes` (Number) Article reading time in minutes.
- `slug` (String) Slug of the article.
- `updated_at` (String) When the listing was updated.
- `url` (String) Full article URL.
- `user` (Map of String) User object of the article.


