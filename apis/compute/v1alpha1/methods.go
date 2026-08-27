// Package v1alpha1 contains the compute v1alpha1 API types for the OpenStack provider.
// This file contains methods to satisfy the resource.Managed interface from crossplane-runtime v2.
package v1alpha1

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

func (m *Server) GetCondition(ct xpv2.ConditionType) xpv2.Condition { return m.Status.GetCondition(ct) }
func (m *Server) SetConditions(c ...xpv2.Condition)                 { m.Status.SetConditions(c...) }
func (m *Server) GetManagementPolicies() xpv2.ManagementPolicies    { return m.Spec.ManagementPolicies }
func (m *Server) SetManagementPolicies(p xpv2.ManagementPolicies)   { m.Spec.ManagementPolicies = p }
func (m *Server) GetDeletionPolicy() xpv2.DeletionPolicy            { return m.Spec.DeletionPolicy }
func (m *Server) SetDeletionPolicy(p xpv2.DeletionPolicy)           { m.Spec.DeletionPolicy = p }
func (m *Server) GetProviderConfigReference() *xpv2.Reference       { return m.Spec.ProviderConfigReference }
func (m *Server) SetProviderConfigReference(p *xpv2.Reference)      { m.Spec.ProviderConfigReference = p }

func (m *KeyPair) GetCondition(ct xpv2.ConditionType) xpv2.Condition {
	return m.Status.GetCondition(ct)
}
func (m *KeyPair) SetConditions(c ...xpv2.Condition)               { m.Status.SetConditions(c...) }
func (m *KeyPair) GetManagementPolicies() xpv2.ManagementPolicies  { return m.Spec.ManagementPolicies }
func (m *KeyPair) SetManagementPolicies(p xpv2.ManagementPolicies) { m.Spec.ManagementPolicies = p }
func (m *KeyPair) GetDeletionPolicy() xpv2.DeletionPolicy          { return m.Spec.DeletionPolicy }
func (m *KeyPair) SetDeletionPolicy(p xpv2.DeletionPolicy)         { m.Spec.DeletionPolicy = p }
func (m *KeyPair) GetProviderConfigReference() *xpv2.Reference     { return m.Spec.ProviderConfigReference }
func (m *KeyPair) SetProviderConfigReference(p *xpv2.Reference)    { m.Spec.ProviderConfigReference = p }
