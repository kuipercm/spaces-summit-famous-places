terraform {
  required_providers {
    google = "~> 3.63.0"
    google-beta = "~> 3.63.0"
  }
}

provider "google" {
  project = var.projectname
  region  = var.region
  zone    = var.zone
  credentials = var.google_credentials
}

provider "google-beta" {
  project = var.projectname
  region  = var.region
  zone    = var.zone
  credentials = var.google_credentials
}

resource "google_cloudbuild_trigger" "famous_places_build_trigger" {
  provider = "google"
  description = "Build famous places repo from github when push to master detected"
  github {
    owner = "kuipercm"
    name = "spaces-summit-famous-places"
    push {
      branch = "^master$"
    }
  }
  filename = "cloudbuild.yaml"
}