package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// LoadBalancer methods
func (m *LoadBalancer) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *LoadBalancer) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *LoadBalancer) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *LoadBalancer) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *LoadBalancer) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *LoadBalancer) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *LoadBalancer) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *LoadBalancer) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}

// Listener methods
func (m *Listener) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Listener) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Listener) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Listener) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Listener) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Listener) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Listener) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *Listener) SetProviderConfigReference(p *xpv2.Reference) { m.Spec.ProviderConfigReference = p }

// Pool methods
func (m *Pool) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Pool) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Pool) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Pool) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Pool) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Pool) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Pool) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Pool) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

// Member methods
func (m *Member) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Member) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Member) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Member) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Member) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Member) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Member) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Member) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

// HealthMonitor methods
func (m *HealthMonitor) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *HealthMonitor) SetConditions(c ...xpv2.Condition) {
	m.Status.SetConditions(c...)
}
func (m *HealthMonitor) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *HealthMonitor) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *HealthMonitor) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *HealthMonitor) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *HealthMonitor) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *HealthMonitor) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}
