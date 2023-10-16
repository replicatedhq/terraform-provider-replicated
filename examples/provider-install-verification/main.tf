module "kind" {
  source = "./kind"
}

module "kubeconfig" {
    source = "./kubeconfig"
    kubeconfig = module.kind.kind_kubeconfig
}

module "app" {
    source = "./app"
    kubeconfig_location = module.kubeconfig.kubeconfig_location
}