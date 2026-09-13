# API Reference - v1beta1

## Overview

All OpenStack resources follow the Crossplane v2 pattern with namespaced scope. Resources are organized by service with their own API groups.

## API Groups and Versions

- **Compute**: `compute.openstack.m.crossplane.io/v1beta1`
- **Networking**: `networking.openstack.m.crossplane.io/v1beta1`
- **BlockStorage**: `blockstorage.openstack.m.crossplane.io/v1beta1`
- **Image**: `image.openstack.m.crossplane.io/v1beta1`
- **Identity**: `identity.openstack.m.crossplane.io/v1beta1`
- **LoadBalancing**: `loadbalancing.openstack.m.crossplane.io/v1beta1`
- **DNS**: `dns.openstack.m.crossplane.io/v1beta1`
- **ProviderConfig**: `openstack.m.crossplane.io/v1beta1`

## Common Spec Fields

All resources support these fields in `spec`:

```yaml
spec:
  # Reference to the ProviderConfig (required)
  providerConfigRef:
    name: provider-name

  # Crossplane v2 management policies (optional, default: [Create, Update, Delete])
  managementPolicies:
  - Create
  - Update
  - Delete
  - Observe

  # Deletion policy (optional, default: Delete)
  # Delete: delete resource when CR is deleted
  # Orphan: leave resource when CR is deleted
  deletionPolicy: Delete

  # Resource parameters for the cloud provider
  forProvider:
    # Provider-specific fields
```

## Common Status Fields

All resources include status conditions:

```yaml
status:
  # Conditions track resource state
  conditions:
  - type: Ready
    status: "True"
    lastTransitionTime: 2024-01-15T10:30:00Z
    reason: Available
    message: "Resource is ready"
  
  # Provider-specific observation
  atProvider:
    # Cloud provider resource details
```

## Compute Service Resources

### Server

Virtual machine instance in Nova.

```yaml
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: Server
metadata:
  name: web-server
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: web-server            # Required: server name
    imageRef: ubuntu-22.04      # Required: image ID or name
    flavorRef: m1.small         # Required: flavor ID or name
    keyName: my-key             # Optional: key pair for SSH
    networks:                   # Optional: networks to attach
    - name: private-network
    securityGroups:             # Optional: security groups
    - web-servers
    userData: "#!/bin/bash..."  # Optional: cloud-init script
    availabilityZone: nova      # Optional: AZ preference
    tags:
    - environment=prod
```

### KeyPair

SSH key pair for server access.

```yaml
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: KeyPair
metadata:
  name: my-keypair
spec:
  forProvider:
    name: my-keypair                    # Required: key pair name
    publicKey: "ssh-rsa AAAA..."        # Optional: public key material
    type: ssh                           # Optional: ssh (default) or x509
```

## Networking Service Resources

### Network

Virtual network in Neutron.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Network
metadata:
  name: private-network
spec:
  forProvider:
    name: private-network   # Required: network name
    # adminStateUp: true    # Optional: enable/disable (default: true)
```

### Subnet

IP subnet within a network.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Subnet
metadata:
  name: private-subnet
spec:
  forProvider:
    name: private-subnet              # Required: subnet name
    networkName: private-network      # Required: parent network
    cidr: 192.168.1.0/24              # Required: CIDR block
    ipVersion: 4                      # Optional: 4 or 6 (default: 4)
    enableDhcp: true                  # Optional: DHCP (default: true)
    gatewayIp: 192.168.1.1            # Optional: gateway IP
    dnsNameservers:                   # Optional: DNS servers
    - 8.8.8.8
    - 8.8.4.4
    allocationPools:                  # Optional: DHCP ranges
    - start: 192.168.1.10
      end: 192.168.1.250
```

### Router

Network router for inter-subnet routing.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Router
metadata:
  name: main-router
spec:
  forProvider:
    name: main-router   # Required: router name
    # adminStateUp: true # Optional: enable/disable (default: true)
```

### RouterInterface

Attachment between router and subnet.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: RouterInterface
metadata:
  name: router-subnet-attach
spec:
  forProvider:
    routerName: main-router      # Required: router name
    subnetName: private-subnet   # Required: subnet name
```

### SecurityGroup

Firewall rules group.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: SecurityGroup
metadata:
  name: web-sg
spec:
  forProvider:
    name: web-sg                      # Required: group name
    description: "Web server rules"   # Optional: description
```

### SecurityGroupRule

Individual security group rule.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: SecurityGroupRule
metadata:
  name: allow-http
spec:
  forProvider:
    securityGroupName: web-sg  # Required: security group
    direction: ingress         # Required: ingress or egress
    protocol: tcp              # Optional: tcp, udp, icmp
    portRangeMin: 80           # Optional: minimum port
    portRangeMax: 80           # Optional: maximum port
    cidr: 0.0.0.0/0            # Optional: IP range
```

### Port

Network interface (network port).

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Port
metadata:
  name: web-port
spec:
  forProvider:
    name: web-port                  # Required: port name
    networkName: private-network    # Required: network
    subnetName: private-subnet      # Optional: subnet
```

### FloatingIP

Public/floating IP address.

```yaml
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: FloatingIP
metadata:
  name: web-floating-ip
spec:
  forProvider:
    floatingNetworkName: public  # Required: external network
    portName: web-port           # Optional: port to associate
    description: "Web server IP" # Optional: description
```

## BlockStorage Service Resources

### Volume

Persistent block storage volume.

```yaml
apiVersion: blockstorage.openstack.m.crossplane.io/v1beta1
kind: Volume
metadata:
  name: data-volume
spec:
  forProvider:
    name: data-volume                   # Required: volume name
    size: 100                           # Required: size in GB
    volumeType: ssd                     # Optional: volume type
    availabilityZone: nova              # Optional: AZ preference
    description: "Data volume"          # Optional: description
```

### VolumeSnapshot

Point-in-time snapshot of a volume.

```yaml
apiVersion: blockstorage.openstack.m.crossplane.io/v1beta1
kind: VolumeSnapshot
metadata:
  name: data-backup
spec:
  forProvider:
    name: data-backup                  # Required: snapshot name
    volumeName: data-volume            # Required: source volume
    description: "Daily backup"        # Optional: description
    forceCreate: false                 # Optional: force if volume busy
```

## Identity Service Resources

### User

Identity service user account.

```yaml
apiVersion: identity.openstack.m.crossplane.io/v1beta1
kind: User
metadata:
  name: app-user
spec:
  forProvider:
    name: app-user                          # Required: username
    password: "password123"                 # Optional: password
    email: app@example.com                  # Optional: email
    enabled: true                           # Optional: enable/disable
    description: "Application service user" # Optional: description
```

### Project

Identity service project (tenant).

```yaml
apiVersion: identity.openstack.m.crossplane.io/v1beta1
kind: Project
metadata:
  name: app-project
spec:
  forProvider:
    name: app-project                  # Required: project name
    description: "Application project" # Optional: description
    enabled: true                      # Optional: enable/disable
```

### Role

RBAC role for authorization.

```yaml
apiVersion: identity.openstack.m.crossplane.io/v1beta1
kind: Role
metadata:
  name: custom-role
spec:
  forProvider:
    name: custom-role              # Required: role name
    description: "Custom role"     # Optional: description
```

## Image Service Resources

### Image

Virtual machine image/snapshot.

```yaml
apiVersion: image.openstack.m.crossplane.io/v1beta1
kind: Image
metadata:
  name: ubuntu-22-04
spec:
  forProvider:
    name: ubuntu-22-04              # Required: image name
    containerFormat: bare           # Required: bare, ami, ari, aki, ovf
    diskFormat: qcow2               # Required: qcow2, raw, vhd, vmdk, iso
    minDiskGb: 10                   # Optional: minimum disk GB
    minRamMb: 512                   # Optional: minimum RAM MB
    isPublic: false                 # Optional: public (default: false)
    tags:
    - ubuntu
    - v22.04
```

## LoadBalancing Service Resources

### LoadBalancer

Octavia load balancer.

```yaml
apiVersion: loadbalancing.openstack.m.crossplane.io/v1beta1
kind: LoadBalancer
metadata:
  name: web-lb
spec:
  forProvider:
    name: web-lb                    # Required: LB name
    description: "Web LB"           # Optional: description
    vipSubnetName: private-subnet   # Required: subnet for VIP
    adminStateUp: true              # Optional: enable/disable
```

### Listener

Load balancer listener (port configuration).

```yaml
apiVersion: loadbalancing.openstack.m.crossplane.io/v1beta1
kind: Listener
metadata:
  name: http-listener
spec:
  forProvider:
    name: http-listener           # Required: listener name
    loadBalancerName: web-lb      # Required: parent LB
    protocol: HTTP                # Required: HTTP, HTTPS, TCP, UDP
    protocolPort: 80              # Required: port number
    adminStateUp: true            # Optional: enable/disable
```

### Pool

Backend pool for load balancer.

```yaml
apiVersion: loadbalancing.openstack.m.crossplane.io/v1beta1
kind: Pool
metadata:
  name: web-pool
spec:
  forProvider:
    name: web-pool                # Required: pool name
    loadBalancerName: web-lb      # Required: parent LB
    listenerName: http-listener   # Required: parent listener
    protocol: HTTP                # Required: protocol
    lbAlgorithm: ROUND_ROBIN      # Required: load balancing algorithm
    adminStateUp: true            # Optional: enable/disable
```

## DNS Service Resources

### Zone

DNS zone (domain).

```yaml
apiVersion: dns.openstack.m.crossplane.io/v1beta1
kind: Zone
metadata:
  name: example-zone
spec:
  forProvider:
    name: example.com              # Required: zone name (FQDN)
    type: PRIMARY                  # Optional: PRIMARY or SECONDARY
    email: admin@example.com        # Optional: admin email
    description: "Main zone"       # Optional: description
    ttl: 3600                      # Optional: default TTL
```

## ProviderConfig

Configuration for provider authentication.

```yaml
apiVersion: openstack.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: default
spec:
  credentials:
    source: Secret                  # Required: Secret source
    secretRef:
      name: openstack-creds        # Required: secret name
      namespace: default           # Required: secret namespace
      key: credentials             # Required: data key in secret
```

## Credential Format

Credentials secret should contain authentication details in JSON or HCL format:

```json
{
  "auth_url": "https://openstack.example.com:5000/v3",
  "username": "crossplane-user",
  "password": "user-password",
  "project_name": "admin",
  "user_domain_name": "Default",
  "project_domain_name": "Default",
  "region": "RegionOne"
}
```

## Management Policy Values

Control resource lifecycle:

- `Create`: Allow creation of resources
- `Update`: Allow updating resource parameters
- `Delete`: Allow deletion of resources
- `Observe`: Allow reading resource status

Default: `[Create, Update, Delete]`

## Status Conditions

Resources track state via conditions:

- `Ready`: Resource is ready for use (status = True/False)
- `Synced`: Resource spec is synced with provider (status = True/False)
- `Creating`: Resource creation in progress
- `Deleting`: Resource deletion in progress

