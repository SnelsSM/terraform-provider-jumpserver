# Jumpserver Terraform Provider

The Jumpserver Terraform Provider allows you to manage [Jumpserver](https://www.jumpserver.org/) resources using Terraform. [Jumpserver](https://www.jumpserver.org/) is an open-source bastion host that helps manage assets, users, and permissions.

## Installation

To use this provider, you need to install it. You can do this by adding it to your Terraform configuration.

### Terraform Configuration

Add the following to your `main.tf` file to use the Jumpserver provider:

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "atwatanmalikm/jumpserver"
      version = "~> 1.0.0"
    }
  }
}

provider "jumpserver" {
  base_url = "https://jumpserver.example.com"
  username = "admin"
  password = "adminpass"
}
```

## Resources

This provider supports the following resources:

* `jumpserver_user`
* `jumpserver_asset`
* `jumpserver_system_user`
* `jumpserver_asset_permission`

## Resource Definitions

For detailed information on each resource, see the following documentation:

* [User Resource](docs/resources/user.md)
* [Asset Resource](docs/resources/asset.md)
* [System User Resource](docs/resources/system_user.md)
* [Asset Permission Resource](docs/resources/asset_permission.md)