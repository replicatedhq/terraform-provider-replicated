output "kubeconfig_location" {
  value = local_sensitive_file.kubeconfig.filename
  sensitive = true
}