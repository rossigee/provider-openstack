// Package v1alpha1 contains the networking v1alpha1 API types for the OpenStack provider.
// This file contains methods to satisfy the resource.Managed interface from crossplane-runtime v2.
package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// ConditionedStatus methods - explicit wrappers to ensure interface satisfaction
// across vendored module boundaries.
func (m *Network) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *Network) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *Network) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *Network) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *Network) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *Network) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *Network) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *Network) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }

func (m *Subnet) GetCondition(ct xpv2.ConditionType) xpv2.Condition { return m.Status.GetCondition(ct) }
func (m *Subnet) SetConditions(c ...xpv2.Condition)                 { m.Status.SetConditions(c...) }
func (m *Subnet) GetManagementPolicies() xpv2.ManagementPolicies    { return m.Spec.ManagementPolicies }
func (m *Subnet) SetManagementPolicies(p xpv2.ManagementPolicies)   { m.Spec.ManagementPolicies = p }
func (m *Subnet) GetDeletionPolicy() xpv2.DeletionPolicy            { return m.Spec.DeletionPolicy }
func (m *Subnet) SetDeletionPolicy(p xpv2.DeletionPolicy)           { m.Spec.DeletionPolicy = p }
func (m *Subnet) GetProviderConfigReference() *xpv2.Reference       { return m.Spec.ProviderConfigReference }
func (m *Subnet) SetProviderConfigReference(p *xpv2.Reference)      { m.Spec.ProviderConfigReference = p }

func (m *Router) GetCondition(ct xpv2.ConditionType) xpv2.Condition { return m.Status.GetCondition(ct) }
func (m *Router) SetConditions(c ...xpv2.Condition)                 { m.Status.SetConditions(c...) }
func (m *Router) GetManagementPolicies() xpv2.ManagementPolicies    { return m.Spec.ManagementPolicies }
func (m *Router) SetManagementPolicies(p xpv2.ManagementPolicies)   { m.Spec.ManagementPolicies = p }
func (m *Router) GetDeletionPolicy() xpv2.DeletionPolicy            { return m.Spec.DeletionPolicy }
func (m *Router) SetDeletionPolicy(p xpv2.DeletionPolicy)           { m.Spec.DeletionPolicy = p }
func (m *Router) GetProviderConfigReference() *xpv2.Reference       { return m.Spec.ProviderConfigReference }
func (m *Router) SetProviderConfigReference(p *xpv2.Reference)      { m.Spec.ProviderConfigReference = p }

func (m *RouterInterface) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *RouterInterface) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *RouterInterface) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *RouterInterface) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *RouterInterface) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *RouterInterface) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *RouterInterface) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *RouterInterface) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}

func (m *SecurityGroup) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *SecurityGroup) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *SecurityGroup) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *SecurityGroup) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *SecurityGroup) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *SecurityGroup) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *SecurityGroup) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *SecurityGroup) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}

func (m *SecurityGroupRule) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *SecurityGroupRule) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *SecurityGroupRule) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *SecurityGroupRule) SetManagementPolicies(p xpv2.ManagementPolicies) {
	m.Spec.ManagementPolicies = p
}
func (m *SecurityGroupRule) GetDeletionPolicy() xpv2.DeletionPolicy  { return m.Spec.DeletionPolicy }
func (m *SecurityGroupRule) SetDeletionPolicy(p xpv2.DeletionPolicy) { m.Spec.DeletionPolicy = p }
func (m *SecurityGroupRule) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *SecurityGroupRule) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}

func (m *FloatingIP) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *FloatingIP) SetConditions(c ...xpv2.Condition) { m.Status.SetConditions(c...) }
func (m *FloatingIP) GetManagementPolicies() xpv2.ManagementPolicies {
	return m.Spec.ManagementPolicies
}
func (m *FloatingIP) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *FloatingIP) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *FloatingIP) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *FloatingIP) GetProviderConfigReference() *xpv2.Reference {
	return m.Spec.ProviderConfigReference
}
func (m *FloatingIP) SetProviderConfigReference(p *xpv2.Reference) {
	m.Spec.ProviderConfigReference = p
}
