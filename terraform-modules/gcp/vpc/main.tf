resource "google_compute_network" "main" {
  name                    = var.name
  auto_create_subnetworks = false
}
