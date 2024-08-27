variable "name" {
  type = string
}

variable "network" {
  type    = any
  default = {}
}

variable "region" {
  type = string
}

variable "cidr" {
  type = string
}

variable "gcp_authentication" {
  type = any
}