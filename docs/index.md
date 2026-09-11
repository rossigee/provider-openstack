# Provider OpenStack Documentation

A Crossplane provider for managing OpenStack cloud resources.

## Quick Links

- [README](README.md) — Basic overview

## Resource Documentation

Individual resource documentation will be added to the [resources/](resources/) folder.

### Compute

| Resource | API Group | Description |
|----------|-----------|-------------|
| Server | `server.openstack.m.crossplane.io/v1beta1` | Virtual servers |

### Block Storage

| Resource | API Group | Description |
|----------|-----------|-------------|
| Volume | `volume.openstack.m.crossplane.io/v1beta1` | Block volumes |

### Networking

| Resource | API Group | Description |
|----------|-----------|-------------|
| Network | `network.openstack.m.crossplane.io/v1beta1` | Networks |
| FloatingIP | `floatingip.openstack.m.crossplane.io/v1beta1` | Floating IPs |

### Identity

| Resource | API Group | Description |
|----------|-----------|-------------|
| User | `user.openstack.m.crossplane.io/v1beta1` | Keystone users |
| Project | `project.openstack.m.crossplane.io/v1beta1` | Keystone projects |

### DNS

| Resource | API Group | Description |
|----------|-----------|-------------|
| Record | `record.openstack.m.crossplane.io/v1beta1` | DNS records |

### Other Resources

See `apis/` directory for all available resources.

## Status

This is a grandfathered provider (exception to the no-Terraform rule).
