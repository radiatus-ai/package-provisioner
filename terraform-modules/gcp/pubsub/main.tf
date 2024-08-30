resource "google_pubsub_topic" "main" {
  name = var.name

  # labels = {
  #   foo = "bar"
  # }

  message_retention_duration = "86600s"
}

# todo: manage this elsewhere, like a pre-provision step. 
# it's awk to ship it here IMO
# locals {
#   apis_to_enable = ["pubsub.googleapis.com"]
# }

# resource "google_project_service" "enable_apis" {
#   for_each = toset(local.apis_to_enable)
#   project  = var.gcp_authentication.project_id
#   service  = each.value

#   disable_on_destroy = false
# }
