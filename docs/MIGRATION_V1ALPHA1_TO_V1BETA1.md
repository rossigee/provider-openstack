# Migration Guide: v1alpha1 to v1beta1

## Overview

Provider OpenStack v1beta1 is a **breaking change** from v1alpha1. It implements Crossplane v2 patterns with namespaced resources and improved management policies.

**Key Changes:**
- Resources are now **namespaced** (not cluster-scoped)
- API groups changed to `*.openstack.m.crossplane.io` format
- `v1alpha1` API versions are **no longer supported**
- All resources use Crossplane v2 ManagedResourceSpec pattern

## Before (v1alpha1)

```yaml
apiVersion: compute.openstack.crossplane.io/v1alpha1
kind: Server
metadata:
  name: my-server
  namespace: default
spec:
  providerConfigRef:
    name: openstack-provider
  forProvider:
    name: web-server
    imageRef: ubuntu-20.04
    flavorRef: m1.small
```

## After (v1beta1)

```yaml
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: Server
metadata:
  name: my-server
  namespace: default  # Required - resources are namespaced
spec:
  providerConfigRef:
    name: openstack-provider
  forProvider:
    name: web-server
    imageRef: ubuntu-20.04
    flavorRef: m1.small
```

## Migration Steps

### 1. Update Provider Installation

```bash
# Install new v1beta1 provider
helm repo add crossplane-contrib https://charts.crossplane.io/contrib
helm install crossplane-openstack crossplane-contrib/provider-openstack \
  --version v0.11.0+ \
  -n crossplane-system
```

### 2. Backup Existing Resources

```bash
# Export all v1alpha1 resources
kubectl get all -A -o yaml > v1alpha1-backup.yaml
```

### 3. Create New v1beta1 Resources

Recreate each resource using v1beta1 API:

```bash
kubectl apply -f - <<EOF
apiVersion: compute.openstack.m.crossplane.io/v1beta1
kind: Server
metadata:
  name: my-server-v2
  namespace: default
spec:
  providerConfigRef:
    name: openstack-provider
  forProvider:
    name: web-server
    imageRef: ubuntu-20.04
    flavorRef: m1.small
EOF
```

### 4. Verify New Resources

```bash
kubectl get servers -n default
kubectl describe server my-server-v2 -n default
```

### 5. Delete v1alpha1 Resources

Once verified, clean up old resources:

```bash
kubectl delete -f v1alpha1-backup.yaml
```

## API Group Mapping

| Service | v1alpha1 | v1beta1 |
|---------|----------|---------|
| Compute | `compute.openstack.crossplane.io/v1alpha1` | `compute.openstack.m.crossplane.io/v1beta1` |
| Networking | `networking.openstack.crossplane.io/v1alpha1` | `networking.openstack.m.crossplane.io/v1beta1` |
| BlockStorage | `blockstorage.openstack.crossplane.io/v1alpha1` | `blockstorage.openstack.m.crossplane.io/v1beta1` |
| Image | `image.openstack.crossplane.io/v1alpha1` | `image.openstack.m.crossplane.io/v1beta1` |
| Identity | `identity.openstack.crossplane.io/v1alpha1` | `identity.openstack.m.crossplane.io/v1beta1` |
| LoadBalancing | `loadbalancing.openstack.crossplane.io/v1alpha1` | `loadbalancing.openstack.m.crossplane.io/v1beta1` |
| DNS | `dns.openstack.crossplane.io/v1alpha1` | `dns.openstack.m.crossplane.io/v1beta1` |

## Namespace Requirement

v1beta1 resources are **namespaced**, not cluster-scoped. When creating resources:

```yaml
metadata:
  namespace: default  # REQUIRED
```

If namespace is omitted, the resource will be created in the active namespace.

## Management Policies (v1beta1 New Feature)

Control whether Crossplane manages resource lifecycle:

```yaml
spec:
  managementPolicies:
  - Create
  - Update
  # - Delete  (omit to prevent deletion)
```

Valid policies: `Create`, `Update`, `Delete`, `Observe`

## Breaking Changes Summary

| Aspect | v1alpha1 | v1beta1 |
|--------|----------|---------|
| Scope | Cluster-scoped | Namespaced |
| API Group | `*.openstack.crossplane.io` | `*.openstack.m.crossplane.io` |
| Version | v1alpha1 | v1beta1 |
| ProviderConfig Scope | Cluster-scoped | Namespaced |
| Spec Structure | Legacy | xpv2.ManagedResourceSpec |
| Management Policies | Not supported | Supported |

## Troubleshooting

### "Resource not found" errors

Check API group and version:
```bash
kubectl api-resources | grep openstack
```

### Namespace issues

Ensure resource and ProviderConfig are in same namespace:
```bash
kubectl get server -A
kubectl get providerconfig -A
```

### Condition not progressing

Check resource conditions:
```bash
kubectl describe server my-server -n default
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-openstack
```

## Support

For issues or questions:
- GitHub Issues: https://github.com/rossigee/provider-openstack/issues
- Crossplane Slack: https://slack.crossplane.io
