package main

import (
	"golang-for-java-developers-training/common/domain"
)

// Re-export types from common/domain for convenience
type Merchant = domain.Merchant
type MerchantStore = domain.MerchantStore

// Re-export functions from common/domain
var NewMerchantStore = domain.NewMerchantStore
var ValidateMerchant = domain.ValidateMerchant
