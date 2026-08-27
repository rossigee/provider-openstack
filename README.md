# Crossplane Provider OpenStack

[![CI](https://img.shields.io/github/actions/workflow/status/rossigee/provider-openstack/ci.yml?branch=master)][build]
[![Version](https://img.shields.io/github/v/release/rossigee/provider-openstack)][releases]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[build]: https://github.com/rossigee/provider-openstack/actions/workflows/ci.yml
[releases]: https://github.com/rossigee/provider-openstack/releases

A native [Crossplane](https://crossplane.io/) provider for managing OpenStack resources.

## Container Registry

- **Primary**: `ghcr.io/rossigee/provider-openstack:v0.10.0`

## Overview

A hand-written Crossplane provider for managing OpenStack resources across compute, networking, storage, identity, DNS, load balancing, and image services.

## Resource Types

All resources use the `openstack.crossplane.io/v1alpha1` API group.

### Networking (11 resources)

| Kind | Description |
|------|-------------|
| Network | Neutron virtual network |
| Subnet | Neutron subnet with CIDR and DHCP settings |
| Router | Neutron router for inter-network routing |
| RouterInterface | Attachment between a router and subnet |
| SecurityGroup | Neutron security group |
| SecurityGroupRule | Individual rule within a security group |
| FloatingIP | Neutron floating IP address |
| Port | Neutron network port |
| SubnetPool | Pool of address prefixes for subnet allocation |
| Trunk | Neutron trunk for VLAN tagging |
| RBACPolicy | RBAC policy for cross-project network sharing |

### Compute (2 resources)

| Kind | Description |
|------|-------------|
| Server | Nova virtual machine instance |
| KeyPair | Nova key pair for SSH access |

### Block Storage (3 resources)

| Kind | Description |
|------|-------------|
| Volume | Cinder block storage volume |
| VolumeType | Cinder volume type with extra specs |
| VolumeSnapshot | Cinder volume snapshot |

### Image (1 resource)

| Kind | Description |
|------|-------------|
| Image | Glance virtual machine image |

### Identity (3 resources)

| Kind | Description |
|------|-------------|
| Project | Keystone project |
| User | Keystone user |
| Role | Keystone role |

### DNS (2 resources)

| Kind | Description |
|------|-------------|
| Zone | Designate DNS zone |
| RecordSet | Designate DNS record set |

### Load Balancing (5 resources)

| Kind | Description |
|------|-------------|
| LoadBalancer | Octavia load balancer |
| Listener | Octavia listener for frontend connections |
| Pool | Octavia backend server pool |
| Member | Octavia pool member (backend server) |
| HealthMonitor | Octavia health check monitor |

**Total: 27 managed resource types**

## Getting Started

### Installation

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
apiVersion: v1
kind: Secret
metadata:
  name: openstack-credentials
  namespace: crossplane-system
type: Opaque
data:
  credentials: <base64-encoded JSON>
---
apiVersion: openstack.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: openstack-credentials
      namespace: crossplane-system
      key: credentials
```

The secret must contain a JSON dictionary with authentication data:

```bash
kubectl create secret generic openstack-credentials \
  --from-file=credentials=config.json \
  --namespace crossplane-system
```

```json
{
  "auth_url": "https://auth.openstack.example/",
  "application_credential_id": "123456789",
  "application_credential_secret": "secret-key"
}
```

### Quick Start

Create a network, subnet, and router:

```yaml
apiVersion: openstack.crossplane.io/v1alpha1
kind: Network
metadata:
  name: my-network
spec:
  forProvider:
    name: my-network
    adminStateUp: true
  providerConfigRef:
    name: default
  deletionPolicy: Delete
---
apiVersion: openstack.crossplane.io/v1alpha1
kind: Subnet
metadata:
  name: my-subnet
spec:
  forProvider:
    name: my-subnet
    networkId: my-network
    cidr: "10.0.0.0/24"
    dnsNameservers:
      - "8.8.8.8"
  providerConfigRef:
    name: default
  deletionPolicy: Delete
```

### Deployment Customization

Use a `DeploymentRuntimeConfig` to customize the provider deployment:

```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: provider-openstack
spec:
  deploymentTemplate:
    spec:
      replicas: 1
      selector: {}
      template:
        spec:
          containers:
          - args:
            - --enable-management-policies
            name: package-runtime
```

## Development

Install required submodules:

```bash
make submodules
```

Apply CRDs and ProviderConfig:

```bash
kubectl apply -f package/crds
kubectl apply -f examples/providerconfig/providerconfig.yaml
```

Run against a Kubernetes cluster:

```bash
make run
```

Build and test:

```bash
make reviewable
make build
```

## License

provider-openstack is under the Apache 2.0 license.
