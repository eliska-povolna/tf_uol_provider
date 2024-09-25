terraform {
  required_providers {
    uol = {
      source = "terraform.local/local/uol"
      version = "1.0.0"
    }
  }
}


provider "uol" {
  email = "demo@ucetnictvi-on-line.cz"
  token = "inZOJH2uLeaTCxU8aDm_7w"
}

resource "uol_contact" "example_contact" {
  name = "John Doe"
}