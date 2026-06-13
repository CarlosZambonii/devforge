variable "subscription_id" {
  description = "ID da subscription Azure"
  type        = string
}

variable "resource_group_name" {
  description = "Nome do resource group"
  type        = string
  default     = "rg-devforge"
}

variable "location" {
  description = "Regiao do Azure"
  type        = string
  default     = "eastus2"
}

variable "acr_name" {
  description = "Nome do Azure Container Registry (deve ser globalmente unico, so letras e numeros)"
  type        = string
}
