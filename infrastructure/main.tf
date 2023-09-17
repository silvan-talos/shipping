provider "aws" {
  region = var.aws_region
}

terraform {
  cloud {
    organization = "tallarry"

    workspaces {
      name = "example-workspace"
    }
  }
}

locals {
  name = "shipping"
}
