output "organization" {
  value = {
    id   = data.github_organization.main.id
    name = var.name
  }
}
