terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
}

# expects a `GITHUB_TOKEN` env var
# this is a PAT token for the user radiatus-gh-bot
# this user was created because PATs are user-specific
# and we shouldn't tie this to a human user
provider "github" {
  owner = var.organization.name
}
