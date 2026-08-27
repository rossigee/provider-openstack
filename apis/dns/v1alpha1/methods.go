// Package v1alpha1 contains the dns v1alpha1 API types for the OpenStack provider.
// This file contains methods to satisfy the resource.Managed interface from crossplane-runtime v2.
package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

func (m *Zone) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Zone) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Zone) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Zone) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Zone) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Zone) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Zone) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Zone) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

func (m *RecordSet) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *RecordSet) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *RecordSet) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *RecordSet) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *RecordSet) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *RecordSet) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *RecordSet) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *RecordSet) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}
