# Jumpserver Provider

The Jumpserver provider allows you to manage Jumpserver resources including users, assets, system users, and asset permissions using Terraform.

## Example Usage

```hcl
provider "jumpserver" {
  base_url = "https://jumpserver.example.com"
  username = "admin"
  password = "adminpass"
}
```

## Argument Reference

* `base_url` (Required) - The base URL of your Jumpserver instance.
* `username` (Required) - The username used to authenticate with Jumpserver.
* `password` (Required) - The password used to authenticate with Jumpserver.