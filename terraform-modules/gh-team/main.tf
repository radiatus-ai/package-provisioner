resource "github_team" "main" {
  name        = var.name
  description = var.description
  privacy     = var.privacy
}

resource "github_team_members" "members" {
  team_id = github_team.main.id

  dynamic "members" {
    for_each = toset(var.members)
    content {
      username = members.value.github
      role     = "member"
    }
  }
}
