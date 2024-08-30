terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.36.0"
    }
  }
}


provider "google" {
  project     = var.gcp_authentication.project_id
  credentials = jsonencode(var.gcp_authentication)
}

