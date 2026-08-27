#!/usr/bin/env bash
# The v1 crossplane-runtime patch is no longer needed - the provider now
# uses crossplane-runtime v2 which already has the required Apply methods.
set -euo pipefail
echo "patch-vendor: v1 runtime patch not needed (using crossplane-runtime v2)"
