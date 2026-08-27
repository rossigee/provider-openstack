// Package v1alpha1 contains the identity v1alpha1 API types for the OpenStack provider.
// This file contains methods to satisfy the resource.Managed interface from crossplane-runtime v2.
package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

func (m *Project) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Project) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Project) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Project) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Project) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Project) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Project) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Project) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

func (m *User) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *User) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *User) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *User) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *User) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *User) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *User) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *User) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

func (m *Role) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Role) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Role) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Role) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Role) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Role) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Role) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Role) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }
