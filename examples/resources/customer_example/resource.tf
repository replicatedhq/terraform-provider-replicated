resource "replicated_customer" "tf_customer" {
  name       = "terraform_customer"
  app_id     = "app_id"
  channel_id = "channel_id"

  expires_at                   = "2028-01-30T15:04:05Z"
  is_kots_install_enabled      = true
  is_gitops_supported          = true
  is_installer_support_enabled = true

  entitlement_values = {
    environment = "test"
  }
}
