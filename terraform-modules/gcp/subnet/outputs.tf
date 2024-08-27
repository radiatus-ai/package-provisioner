output "subnet" {
  value = {
    id = google_compute_subnetwork.main.id
  }
}
