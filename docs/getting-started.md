# Getting Started with Provider OpenStack

Complete guide to installing and using provider-openstack v1beta1 with Crossplane.

## Prerequisites

- **Kubernetes cluster** with Crossplane v2.5+ installed
- **OpenStack** cloud environment (Nova, Neutron, Cinder, etc.)
- OpenStack credentials with appropriate permissions
- `kubectl` configured to access your cluster

## Installation

### 1. Install Crossplane

```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm install crossplane \
  crossplane-stable/crossplane \
  -n crossplane-system \
  --create-namespace
```

Wait for Crossplane to be ready:
```bash
kubectl wait -n crossplane-system --for=condition=Ready pods -l app.kubernetes.io/instance=crossplane --timeout=300s
```

### 2. Install Provider OpenStack

```bash
kubectl crossplane install provider ghcr.io/rossigee/provider-openstack:v1.3.2
```

Verify installation:
```bash
kubectl get providers
kubectl describe provider provider-openstack
```

## Configuration

### 1. Create OpenStack Credentials Secret

Create a secret containing OpenStack authentication details:

```bash
kubectl create secret generic openstack-creds \
  -n default \
  --from-literal=credentials='
{
  "auth_url": "https://openstack.example.com:5000/v3",
  "username": "crossplane-user",
  "password": "your-secure-password",
  "project_name": "admin",
  "user_domain_name": "Default",
  "project_domain_name": "Default",
  "region": "RegionOne"
}'
```

Or using a file:

```bash
cat > credentials.json <<EOF
{
  "auth_url": "https://openstack.example.com:5000/v3",
  "username": "crossplane-user",
  "password": "your-secure-password",
  "project_name": "admin",
  "user_domain_name": "Default",
  "project_domain_name": "Default",
  "region": "RegionOne"
}
EOF

kubectl create secret generic openstack-creds \
  -n default \
  --from-file=credentials=credentials.json
```

### 2. Create ProviderConfig

Apply the ProviderConfig to connect to OpenStack:

```yaml
apiVersion: openstack.m.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
  namespace: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: openstack-creds
      namespace: default
      key: credentials
```

Save as `provider-config.yaml` and apply:
```bash
kubectl apply -f provider-config.yaml
```

Verify:
```bash
kubectl get providerconfig
kubectl describe providerconfig default
```

## First Resource: Create a Virtual Machine

### 1. Create a KeyPair (SSH access)

```yaml
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: KeyPair
metadata:
  name: demo-key
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: demo-key
    type: ssh
```

Apply:
```bash
kubectl apply -f keypair.yaml
```

### 2. Create a Server (VM)

```yaml
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: Server
metadata:
  name: demo-server
  namespace: default
spec:
  providerConfigRef:
    name: default
  
  forProvider:
    name: demo-web-server
    imageRef: ubuntu-22.04      # Use your image name/ID
    flavorRef: m1.small         # Use your flavor name/ID
    keyName: demo-key
    networks:
    - name: private
    securityGroups:
    - default
    tags:
    - demo=true
    
  managementPolicies:
  - Create
  - Update
  
  deletionPolicy: Delete
```

Apply:
```bash
kubectl apply -f server.yaml
```

### 3. Monitor Resource Status

Watch the server creation:
```bash
kubectl describe server demo-server -n default
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-openstack -f
```

Check conditions:
```bash
kubectl get server demo-server -o jsonpath='{.status.conditions}'
```

## Common Tasks

### List All Resources

```bash
# List all servers
kubectl get servers -n default

# List all networks
kubectl get networks -n default

# List all volumes
kubectl get volumes -n default
```

### Describe a Resource

```bash
kubectl describe server demo-server -n default
kubectl describe network private -n default
```

### Update a Resource

Edit the resource:
```bash
kubectl edit server demo-server -n default
```

Or patch it:
```bash
kubectl patch server demo-server \
  -p '{"spec":{"forProvider":{"tags":["demo=true","updated=true"]}}}' \
  --type merge
```

### Delete a Resource

```bash
# Delete will remove both the CR and OpenStack resource (unless deletionPolicy: Orphan)
kubectl delete server demo-server -n default
```

## Network Setup Example

Create a complete network stack:

```yaml
---
# Network
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Network
metadata:
  name: demo-network
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: demo-network

---
# Subnet
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: Subnet
metadata:
  name: demo-subnet
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: demo-subnet
    networkName: demo-network
    cidr: 192.168.10.0/24
    gatewayIp: 192.168.10.1
    enableDhcp: true

---
# Security Group
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: SecurityGroup
metadata:
  name: demo-sg
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    name: demo-sg
    description: "Demo security group"

---
# Allow SSH
apiVersion: networking.openstack.m.crossplane.io/v1beta1
kind: SecurityGroupRule
metadata:
  name: allow-ssh
  namespace: default
spec:
  providerConfigRef:
    name: default
  forProvider:
    securityGroupName: demo-sg
    direction: ingress
    protocol: tcp
    portRangeMin: 22
    portRangeMax: 22
    cidr: 0.0.0.0/0
```

Apply:
```bash
kubectl apply -f network-stack.yaml
```

## Troubleshooting

### Provider not Ready

```bash
# Check provider status
kubectl get provider provider-openstack
kubectl describe provider provider-openstack

# Check logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-openstack
```

### Resource Stuck in Creating

```bash
# Check resource status
kubectl describe server demo-server

# Check conditions
kubectl get server demo-server -o jsonpath='{.status.conditions}' | jq

# Check provider logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-openstack -f
```

### Namespace Issues

Resources and ProviderConfig must be in the same namespace:

```bash
# Check both are in same namespace
kubectl get servers -A
kubectl get providerconfig -A

# Or create in crossplane-system
kubectl apply -f resource.yaml -n crossplane-system
```

### Authentication Errors

Verify credentials:
```bash
# Check secret exists
kubectl get secret openstack-creds -n default

# Verify format
kubectl get secret openstack-creds -n default -o jsonpath='{.data.credentials}' | base64 -d | jq
```

## Next Steps

1. **Explore more resources**: See [API Reference](API_REFERENCE.md) for all available resources
2. **Complex examples**: Check [examples/](../examples/) directory
3. **Migration**: If upgrading from v1alpha1, see [Migration Guide](MIGRATION_V1ALPHA1_TO_V1BETA1.md)
4. **Best practices**: Read [Crossplane documentation](https://docs.crossplane.io)

## Getting Help

- **GitHub Issues**: https://github.com/rossigee/provider-openstack/issues
- **Crossplane Slack**: https://slack.crossplane.io
- **OpenStack Documentation**: https://docs.openstack.org
