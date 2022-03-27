variable "api_key" {
  type        = string
  description = "Forem API key"
}

variable "host" {
  type        = string
  description = "Forem Host"
  default     = "http://localhost:3000/api"
}
