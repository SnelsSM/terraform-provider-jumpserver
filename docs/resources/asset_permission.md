# `jumpserver_asset_permission` Resource

The jumpserver_asset_permission resource allows you to create and manage asset permissions in Jumpserver. Asset permissions define which users have access to which assets and with which system users.

## Example Usage

```hcl
resource "jumpserver_asset_permission" "example_permission" {
  name                  = "User 1"
  is_active             = true
  users_display         = ["user1"]
  assets_display        = ["server-vm1"]
  system_users_display  = ["student"]
}
```

## Argument Reference

* `name` - (Required) The name of the asset permission.
* `is_active` - (Optional) Whether the permission is active.
* `users_display` - (Optional) List of users the permission applies to.
* `assets_display` - (Optional) List of assets the permission applies to.
* `system_users_display` - (Optional) List of system users the permission applies to.

## Attribute Reference

* `id` - The ID of the asset permission.
* `name` - The name of the asset permission.
* `is_active` - Whether the permission is active.
* `users_display` - List of users the permission applies to.
* `assets_display` - List of assets the permission applies to.
* `system_users_display` - List of system users the permission applies to.