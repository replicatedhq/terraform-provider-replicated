output "kind_kubeconfig" {
  value = replicated_cluster.tf_cluster.kubeconfig
  sensitive = true
}