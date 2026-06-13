terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.100"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id                 = var.subscription_id
  skip_provider_registration = true
}

resource "azurerm_resource_group" "devforge" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_container_registry" "acr" {
  name                = var.acr_name
  resource_group_name = azurerm_resource_group.devforge.name
  location            = azurerm_resource_group.devforge.location
  sku                 = "Basic"
  admin_enabled       = true
}
