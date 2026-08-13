# Crossplane Provider OpenStack

[![CI](https://img.shields.io/github/actions/workflow/status/rossigee/provider-openstack/ci.yml?branch=master)][build]
[![Version](https://img.shields.io/github/v/release/rossigee/provider-openstack)][releases]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[build]: https://github.com/rossigee/provider-openstack/actions/workflows/ci.yml
[releases]: https://github.com/rossigee/provider-openstack/releases

`provider-openstack` is a [Crossplane](https://crossplane.io/) provider for managing OpenStack resources.

## Container Registry

- **Primary**: `ghcr.io/rossigee/provider-openstack:v0.10.0`

## Overview

A Crossplane provider for managing OpenStack resources.

## Features

- **Compute**: Manage OpenStack virtual machines and flavors
- **Networking**: Virtual networks, routers, and security groups
- **Storage**: Block and object storage management
- **Identity**: Keystone authentication and authorization

## Resource Types

This is an [upjet](https://github.com/crossplane/upjet)-generated provider covering the following OpenStack service categories, each with multiple resource kinds:

| Category | API Group prefix | Examples |
|----------|-------------------|----------|
| Compute | `compute.openstack.crossplane.io` | InstanceV2, FlavorV2, AggregateV2 |
| Networking | `networking.openstack.crossplane.io` | Network, Subnet, Port |
| Firewall (FWaaS) | `fw.openstack.crossplane.io` | Firewall groups, policies, rules |
| VPN-as-a-Service | `vpnaas.openstack.crossplane.io` | IPSec/IKE policies, VPN services |
| Load Balancer (Octavia) | `lb.openstack.crossplane.io` | Load balancers, listeners, pools, monitors |
| Block Storage (Cinder) | `blockstorage.openstack.crossplane.io` | Volumes, volume types, snapshots |
| Object Storage (Swift) | `objectstorage.openstack.crossplane.io` | Containers |
| Shared File System (Manila) | `sharedfilesystem.openstack.crossplane.io` | Shares, share networks |
| Images (Glance) | `images.openstack.crossplane.io` | Images |
| Identity (Keystone) | `identity.openstack.crossplane.io` | Projects, users, roles |
| Key Manager (Barbican) | `keymanager.openstack.crossplane.io` | Secrets, containers |
| DNS (Designate) | `dns.openstack.crossplane.io` | Zones, recordsets |
| Orchestration (Heat) | `orchestration.openstack.crossplane.io` | Stacks |
| Database (Trove) | `db.openstack.crossplane.io` | Database instances |
| Container Infrastructure (Magnum) | `containerinfra.openstack.crossplane.io` | Clusters, cluster templates |

See `apis/` for the full generated CRD list and `examples/` for per-resource usage examples.

## Getting Started

### Installation

Install the provider:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-openstack
spec:
  package: ghcr.io/rossigee/provider-openstack:v0.10.0
```

### Configuration

```yaml
---
# Providerconfig that referers to the secret
apiVersion: openstack.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: provider-openstack-config
spec:
  credentials:
    source: Secret
    secretRef:
      key: config
      name: provider-openstack-config
      namespace: crossplane

---
# Secret that stores credentials and other configuration
apiVersion: v1
kind: Secret
metadata:
  name: provider-openstack-config
  namespace: crossplane
type: Opaque
data:
  config: <see below>
```

The secret key must contain a json dictionary that provides the authentication data.
You can create the secret via this command:

```bash
kubectl create secret generic provider-openstack-config --from-file=config=config.json --namespace crossplane
```

```json
// config.json
{
  "auth_url": "https://auth.openstack.example/",
  "application_credential_id": "123456789",
  "application_credential_secret": "secret-key"
}
```

Check [Terraform OpenStack provider docs](https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs#configuration-reference) to see available configuration settings. Currently not all options of the upstream provider are supported. Check [client code](https://github.com/crossplane-contrib/provider-openstack/blob/main/internal/clients/openstack.go#L66) to see if your option is supported. If something is missing, please open a new issue.


### Deployment Customization

You can use a `DeploymentRuntimeConfig` to provide custom arguments or otherwise modify the provider deployment

Available command line arguments can be found [here](cmd/provider/main.go)

```yaml
---
# Create a DeploymentRuntimeConfig to customize the provider deployment
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: provider-openstack
spec:
  deploymentTemplate:
    spec:
      # Control replica count to temporary disable deployment. Do not scale more than 1 replica.
      replicas: 1
      selector: {}
      template:
        metadata:
          annotations:
            # Add annotations, e.g. to enable metrics scraping
            prometheus.io/path: /metrics
            prometheus.io/port: "8080"
            prometheus.io/scrape: "true"
        spec:
          containers:
          - args:
            # Add command line arguments, e.g. to enable management policies
            - --enable-management-policies
            name: package-runtime

---
# Add this to your provider resource to reference the DeploymentRuntimeConfig
spec:
  runtimeConfigRef:
    apiVersion: pkg.crossplane.io/v1beta1
    kind: DeploymentRuntimeConfig
    name: provider-openstack
```

## Development

Install the required submodules to build and run:

```bash
make submodules
```

Apply the Current CRDs and a providerConfig:

```bash
kubectl apply -f package/crds
kubectl apply -f examples/providerconfig/providerconfig.yaml
```

Run against a Kubernetes cluster: (make sure to apply CRDs and providerConfig)

```bash
make run
```

Run a testbuild with linting:

```bash
make reviewable
```

Build binary:

```bash
make build
```

### Release a new version (Maintainer)

- Update Changelog (Add new Version & Date)
- Create or merge to existing release branch (release-v(major).(minor))
- Run Release pipeline on release branch, using specific version as parameter

## Contributing

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/crossplane-contrib/provider-openstack/issues).

## License

provider-openstack is under the Apache 2.0 license.
