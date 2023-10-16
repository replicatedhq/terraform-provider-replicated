terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">=2.23.0"
    }
  }
}

variable "kubeconfig_location" {
  description = "Kubeconfig location to use"
  type        = string
  sensitive = true
}

provider "kubernetes" {
  config_path    = "${var.kubeconfig_location}"
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "my-first-namespace"
  }
}