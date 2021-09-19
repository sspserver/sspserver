package models

import "geniusrabbit.dev/sspserver/internal/billing"

// VirtualTarget it's a wrapper which don't have ID
type VirtualTarget struct {
	trg Target
}

// NewVirtualTarget wrapper of exists target
func NewVirtualTarget(trg Target) VirtualTarget {
	return VirtualTarget{trg: trg}
}

// Target object accessor
func (tw VirtualTarget) Target() Target { return tw.trg }

// ID of object (Zone OR SmartLink only)
func (tw VirtualTarget) ID() uint64 { return 0 }

// Codename of the target (equal to tagid)
func (tw VirtualTarget) Codename() string { return "" }

// PurchasePrice gives the price of view from external resource
func (tw VirtualTarget) PurchasePrice(a Action) billing.Money { return tw.trg.PurchasePrice(a) }

// RevenueShareFactor of current target
func (tw VirtualTarget) RevenueShareFactor() float64 { return tw.trg.RevenueShareFactor() }

// ComissionShareFactor of current target
func (tw VirtualTarget) ComissionShareFactor() float64 { return tw.trg.ComissionShareFactor() }

// Company object
func (tw VirtualTarget) Company() *Company { return tw.trg.Company() }

// CompanyID of current target
func (tw VirtualTarget) CompanyID() uint64 { return tw.trg.CompanyID() }
