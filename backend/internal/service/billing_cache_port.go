package service

import (
	"time"
)

// SubscriptionCacheData represents cached subscription data
type SubscriptionCacheData struct {
	Status       string
	ExpiresAt    time.Time
	DailyUsage   float64
	WeeklyUsage  float64
	MonthlyUsage float64
	Version      int64

	// Plan live reference for request count limits
	PlanID *int64

	// Request count usage (shared window_start with USD)
	DailyUsageRequests   int64
	WeeklyUsageRequests  int64
	MonthlyUsageRequests int64

	// Plan request limits (live from current plan config, nil = unlimited)
	DailyRequestLimit    *int64
	WeeklyRequestLimit   *int64
	MonthlyRequestLimit  *int64

	// PlanVersion for detecting plan config changes
	PlanVersion int64
}
