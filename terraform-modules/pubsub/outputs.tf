output "topic" {
  value = {
    id = google_pubsub_topic.main.id
  }
}
