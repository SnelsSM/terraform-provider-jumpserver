# Jumpserver Provider

The Jumpserver provider allows you to manage Jumpserver resources including users, assets, system users, and asset permissions using Terraform.

## Example Usage

Via username and password:  
```hcl
provider "jumpserver" {
  base_url = "https://jumpserver.example.com"
  username = "admin"
  password = "adminpass"
}
```

Via access key:  
```hcl
provider "jumpserver" {
  base_url = "https://jumpserver.example.com"
  access_key = "XXXXXXX"
  secret_key = "YYYYYYY"
}
```

## Argument Reference

* `base_url` (Required) - The base URL of your Jumpserver instance. Can also be set via environment variable JUMPSERVER_BASE_URL;
* `username` (Optional) - The username used to authenticate with Jumpserver. Can also be set via environment variable JUMPSERVER_USERNAME;
* `password` (Optional) - The password used to authenticate with Jumpserver. Can also be set via environment variable JUMPSERVER_PASSWORD;
* `access_key` (Optional) - Jumpserver API Access Key. Can also be set via environment variable JUMPSERVER_ACCESS_KEY;
* `secret_key` (Optional) - Jumpserver API Secret Key. Can also be set via environment variable JUMPSERVER_SECRET_KEY;
* `skip_tls_verify` (Optional) - If true, skip SSL certificate validation (insecure). Can also be set via environment variable JUMPSERVER_SKIP_TLS_VERIFY. Default: false;