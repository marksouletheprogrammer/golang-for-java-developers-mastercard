package main

import (
	"golang-for-java-developers-training/common/external"
)

// Re-export functions from common/external for convenience
var FetchInventoryLevel = external.FetchInventoryLevel
var FetchDynamicPrice = external.FetchDynamicPrice
var FetchReviewSummary = external.FetchReviewSummary
