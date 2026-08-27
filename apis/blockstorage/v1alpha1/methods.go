// Package v1alpha1 contains the blockstorage v1alpha1 API types for the OpenStack provider.
// This file contains methods to satisfy the resource.Managed interface from crossplane-runtime v2.
package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

func (m *Volume) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Volume) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Volume) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Volume) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Volume) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Volume) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Volume) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Volume) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

func (m *VolumeType) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *VolumeType) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *VolumeType) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *VolumeType) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *VolumeType) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *VolumeType) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *VolumeType) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *VolumeType) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}

func (m *VolumeSnapshot) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *VolumeSnapshot) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *VolumeSnapshot) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *VolumeSnapshot) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *VolumeSnapshot) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *VolumeSnapshot) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *VolumeSnapshot) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *VolumeSnapshot) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}
