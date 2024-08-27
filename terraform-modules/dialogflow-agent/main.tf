resource "google_dialogflow_cx_agent" "main" {
  display_name               = var.name
  location                   = "global"
  default_language_code      = "en"
  supported_language_codes   = ["fr", "de", "es"]
  time_zone                  = "America/New_York"
  description                = var.description
  avatar_uri                 = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
  enable_stackdriver_logging = false
  enable_spell_correction    = true
  speech_to_text_settings {
    enable_speech_adaptation = true
  }
  project = "rad-dev-canvas-kwm6"
}
