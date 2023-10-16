terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">=2.4.0"
    }
  }
}

provider "local" {
}

variable "kubeconfig" {
  description = "Kubeconfig to use"
  type        = string
  sensitive = true
}

resource "local_sensitive_file" "kubeconfig" {
    content  = "${var.kubeconfig}"
    filename = "/tmp/config"
}