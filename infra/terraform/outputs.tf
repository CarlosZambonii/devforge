output "acr_login_server" {
  description = "URL do registry para docker login/push"
  value       = azurerm_container_registry.acr.login_server
}

output "acr_admin_username" {
  description = "Usuario admin do ACR"
  value       = azurerm_container_registry.acr.admin_username
}

output "acr_admin_password" {
  description = "Senha admin do ACR"
  value       = azurerm_container_registry.acr.admin_password
  sensitive   = true
}
