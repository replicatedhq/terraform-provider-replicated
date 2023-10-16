terraform {
  required_providers {
    replicated = {
      source = "registry.terraform.io/replicated/replicated"
    }
  }
}

provider "replicated" {
}

resource "replicated_cluster" "tf_cluster" {
  name          = "terraformCLuster"
  distribution  = "kind"
  wait_duration = "20m"
  ttl           = "10m"
  instance_type = "r1.large"
  nodes         = 1
  disk          = 100
}

