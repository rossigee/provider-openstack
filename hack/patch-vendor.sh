#!/usr/bin/env bash
# Patches vendored github.com/crossplane/crossplane-runtime (v1) so its
# unstructured.WrapperClient satisfies sigs.k8s.io/controller-runtime v0.24+
# client interfaces, which require Apply methods on Writer/SubResourceWriter.
#
# Upstream gap: crossplane-runtime v1 (<= v1.21.0-rc.0) predates those
# interface changes; only the /v2 module line is patched upstream.
# Remove this script once a v1 release carries the fix.
set -euo pipefail

f="vendor/github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/client.go"

if [[ ! -f "$f" ]]; then
  echo "patch-vendor: $f not found; run 'go mod vendor' first" >&2
  exit 1
fi

if grep -q "func (c \*WrapperClient) Apply" "$f"; then
  echo "patch-vendor: already applied"
  exit 0
fi

python3 - "$f" <<'PY'
import sys

path = sys.argv[1]
with open(path) as fh:
    src = fh.read()

client_anchor = """// DeleteAllOf deletes all objects of the given type matching the given options.
func (c *WrapperClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	if u, ok := obj.(Wrapper); ok {
		return c.kube.DeleteAllOf(ctx, u.GetUnstructured(), opts...)
	}
	return c.kube.DeleteAllOf(ctx, obj, opts...)
}
"""

client_patch = client_anchor + """
// Apply applies the given apply configuration.
func (c *WrapperClient) Apply(ctx context.Context, config runtime.ApplyConfiguration, opts ...client.ApplyOption) error {
	return c.kube.Apply(ctx, config, opts...)
}
"""

status_anchor = """// Patch patches the given object's subresource. obj must be a struct
// pointer so that obj can be updated with the content returned by the
// Server.
func (c *wrapperStatusClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	if u, ok := obj.(Wrapper); ok {
		return c.kube.Patch(ctx, u.GetUnstructured(), patch, opts...)
	}
	return c.kube.Patch(ctx, obj, patch, opts...)
}
"""

status_patch = status_anchor + """
// Apply applies the given apply configuration subresource.
func (c *wrapperStatusClient) Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.SubResourceApplyOption) error {
	return c.kube.Apply(ctx, obj, opts...)
}
"""

for anchor, patched in ((client_anchor, client_patch), (status_anchor, status_patch)):
    if anchor not in src:
        sys.exit("patch-vendor: anchor block not found in " + path)
    src = src.replace(anchor, patched, 1)

with open(path, "w") as fh:
    fh.write(src)
print("patch-vendor: applied")
PY

gofmt -l -w "$f"
