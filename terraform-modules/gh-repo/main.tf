resource "github_repository" "main" {
  name        = var.name
  description = var.description

  visibility = "private"
}
