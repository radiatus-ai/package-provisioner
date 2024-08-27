resource "google_cloud_run_v2_service" "main" {
  name     = var.name
  location = "us-central1"
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
}

# resource "google_compute_subnetwork" "custom_test" {
#   name          = "run-subnetwork"
#   ip_cidr_range = "10.2.0.0/28"
#   region        = "us-central1"
#   network       = var.network.id
# }

module "glb-0" {
  source              = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-lb-app-ext?ref=v32.0.0"
  project_id          = "rad-dev-canvas-kwm6"
  name                = "glb-test-0"
  use_classic_version = false
  backend_service_configs = {
    default = {
      backends = [
        { backend = "neg-0" }
      ]
      health_checks = []
    }
  }
  # with a single serverless NEG the implied default health check is not needed
  health_check_configs = {}
  neg_configs = {
    neg-0 = {
      cloudrun = {
        region = var.region
        target_service = {
          name = google_cloud_run_v2_service.main.name
        }
      }
    }
  }
}
