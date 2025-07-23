# `jumpserver_node` Resource

The jumpserver_node resource allows you to create and manage nodes in Jumpserver.

## Example Usage

```hcl
resource "jumpserver_node" "example_node" {
  name      = "example"
  parent    = "DEFAULT"
}
```

## Argument Reference

* `name` - (Required) Node name.
* `parent` - (Required) Parent node name


## Attribute Reference

* `id` - The ID of the node.
* `name` - The name of the new node.
* `parent` - The name of the parent node.
