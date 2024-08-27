resource "google_pubsub_topic" "main" {
  name = var.name

  labels = {
    foo = "bar"
  }

  message_retention_duration = "86600s"
}
