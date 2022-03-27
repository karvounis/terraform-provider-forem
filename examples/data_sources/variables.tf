variable "api_key" {
  type        = string
  description = "Forem API key"
}

variable "host" {
  type        = string
  description = "Forem Host"
  default     = "http://localhost:3000/api"
}

variable "article_id" {
  type        = string
  description = "ID of the article to use for the data source"
}

variable "listing_id" {
  type        = string
  description = "ID of the listing to use for the data source"
}

variable "user_id" {
  type        = string
  description = "ID of the user to use for the data source"
  default     = "1"
}

variable "user_username" {
  type        = string
  description = "Username of the user to use for the data source"
  default     = "admin_mcadmin"
}
