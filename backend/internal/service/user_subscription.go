package service

import "time"

// SubscriptionPlanLimitData 轻量 Plan 模型，仅包含热路径需要的限额字段。
type SubscriptionPlanLimitData struct {
	ID                  int64
	DailyRequestLimit   *int64
	WeeklyRequestLimit  *int64
	MonthlyRequestLimit *int64
	UpdatedAt           time.Time
}

type UserSubscription struct {
	ID      int64
	UserID  int64
	GroupID int64

	StartsAt  time.Time
	ExpiresAt time.Time
	Status    string

	DailyWindowStart   *time.Time
	WeeklyWindowStart  *time.Time
	MonthlyWindowStart *time.Time

	DailyUsageUSD   float64
	WeeklyUsageUSD  float64
	MonthlyUsageUSD float64

	// Plan 实时引用
	PlanID *int64

	// 请求次数用量跟踪（与 USD 共享 window_start）
	DailyUsageRequests   int64
	WeeklyUsageRequests  int64
	MonthlyUsageRequests int64

	// Plan 关联（由 repository WithPlan() 加载）
	Plan *SubscriptionPlanLimitData

	AssignedBy *int64
	AssignedAt time.Time
	Notes      string

	CreatedAt time.Time
	UpdatedAt time.Time

	User           *User
	Group          *Group
	AssignedByUser *User
}

func (s *UserSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.ExpiresAt)
}

func (s *UserSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *UserSubscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	return int(time.Until(s.ExpiresAt).Hours() / 24)
}

func (s *UserSubscription) IsWindowActivated() bool {
	return s.DailyWindowStart != nil || s.WeeklyWindowStart != nil || s.MonthlyWindowStart != nil
}

func (s *UserSubscription) HasOneTimeDailyQuota() bool {
	if s == nil || s.StartsAt.IsZero() || s.ExpiresAt.IsZero() {
		return false
	}
	return !s.ExpiresAt.After(s.StartsAt.AddDate(0, 0, 1))
}

func (s *UserSubscription) NeedsDailyReset() bool {
	return s.NeedsDailyResetAt(time.Now())
}

func (s *UserSubscription) NeedsDailyResetAt(now time.Time) bool {
	if s.DailyWindowStart == nil {
		return false
	}
	if s.HasOneTimeDailyQuota() {
		return false
	}
	return !now.Before(s.DailyWindowStart.Add(24 * time.Hour))
}

func (s *UserSubscription) NeedsWeeklyReset() bool {
	if s.WeeklyWindowStart == nil {
		return false
	}
	return time.Since(*s.WeeklyWindowStart) >= 7*24*time.Hour
}

func (s *UserSubscription) NeedsMonthlyReset() bool {
	if s.MonthlyWindowStart == nil {
		return false
	}
	return time.Since(*s.MonthlyWindowStart) >= 30*24*time.Hour
}

func (s *UserSubscription) DailyResetTime() *time.Time {
	if s.DailyWindowStart == nil {
		return nil
	}
	if s.HasOneTimeDailyQuota() {
		t := s.ExpiresAt
		return &t
	}
	t := s.DailyWindowStart.Add(24 * time.Hour)
	return &t
}

func (s *UserSubscription) WeeklyResetTime() *time.Time {
	if s.WeeklyWindowStart == nil {
		return nil
	}
	t := s.WeeklyWindowStart.Add(7 * 24 * time.Hour)
	return &t
}

func (s *UserSubscription) MonthlyResetTime() *time.Time {
	if s.MonthlyWindowStart == nil {
		return nil
	}
	t := s.MonthlyWindowStart.Add(30 * 24 * time.Hour)
	return &t
}

// USD 限额检查（limit 来自 Group）

func (s *UserSubscription) CheckDailyLimit(group *Group, additionalCost float64) bool {
	if !group.HasDailyLimit() {
		return true
	}
	return s.DailyUsageUSD+additionalCost <= *group.DailyLimitUSD
}

func (s *UserSubscription) CheckWeeklyLimit(group *Group, additionalCost float64) bool {
	if !group.HasWeeklyLimit() {
		return true
	}
	return s.WeeklyUsageUSD+additionalCost <= *group.WeeklyLimitUSD
}

func (s *UserSubscription) CheckMonthlyLimit(group *Group, additionalCost float64) bool {
	if !group.HasMonthlyLimit() {
		return true
	}
	return s.MonthlyUsageUSD+additionalCost <= *group.MonthlyLimitUSD
}

func (s *UserSubscription) CheckAllLimits(group *Group, additionalCost float64) (daily, weekly, monthly bool) {
	daily = s.CheckDailyLimit(group, additionalCost)
	weekly = s.CheckWeeklyLimit(group, additionalCost)
	monthly = s.CheckMonthlyLimit(group, additionalCost)
	return
}

// 请求次数限额检查（limit 来自 Plan，三态语义：nil=不限, 0=禁止, >0=限额）

func (s *UserSubscription) CheckDailyRequestLimit(limit *int64, additionalCount int64) bool {
	if limit == nil {
		return true
	}
	return s.DailyUsageRequests+additionalCount <= *limit
}

func (s *UserSubscription) CheckWeeklyRequestLimit(limit *int64, additionalCount int64) bool {
	if limit == nil {
		return true
	}
	return s.WeeklyUsageRequests+additionalCount <= *limit
}

func (s *UserSubscription) CheckMonthlyRequestLimit(limit *int64, additionalCount int64) bool {
	if limit == nil {
		return true
	}
	return s.MonthlyUsageRequests+additionalCount <= *limit
}

// limitsFromSubscriptionPlan 从 subscription 的 Plan 中提取请求次数限额。
// 如果 PlanID 非 nil 但 Plan 未加载，返回 -1 作为哨兵值，表示计费不可用。
// 调用方应检查返回值是否为 -1 并返回错误。
func limitsFromSubscriptionPlan(sub *UserSubscription) (daily, weekly, monthly *int64) {
	if sub.PlanID == nil {
		return nil, nil, nil
	}
	if sub.Plan == nil {
		// Plan not loaded: return sentinel value -1 to signal error
		negOne := int64(-1)
		return &negOne, &negOne, &negOne
	}
	return sub.Plan.DailyRequestLimit, sub.Plan.WeeklyRequestLimit, sub.Plan.MonthlyRequestLimit
}

// isPlanNotLoadedSentinel checks if limitsFromSubscriptionPlan returned the -1 sentinel
// indicating PlanID is set but Plan was not loaded. Returns true if the caller should reject.
func isPlanNotLoadedSentinel(limit *int64) bool {
	return limit != nil && *limit == -1
}

