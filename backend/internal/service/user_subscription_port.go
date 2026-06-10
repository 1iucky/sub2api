package service

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// SubscriptionRef 用于缓存失效时按 plan 查询活跃订阅引用。
type SubscriptionRef struct {
	ID      int64
	UserID  int64
	GroupID int64
}

type UserSubscriptionRepository interface {
	Create(ctx context.Context, sub *UserSubscription) error
	GetByID(ctx context.Context, id int64) (*UserSubscription, error)
	GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error)
	GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*UserSubscription, error)
	Update(ctx context.Context, sub *UserSubscription) error
	Delete(ctx context.Context, id int64) error

	ListByUserID(ctx context.Context, userID int64) ([]UserSubscription, error)
	ListActiveByUserID(ctx context.Context, userID int64) ([]UserSubscription, error)
	ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]UserSubscription, *pagination.PaginationResult, error)
	List(ctx context.Context, params pagination.PaginationParams, userID, groupID *int64, status, platform, sortBy, sortOrder string) ([]UserSubscription, *pagination.PaginationResult, error)

	ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error)
	ExtendExpiry(ctx context.Context, subscriptionID int64, newExpiresAt time.Time) error
	UpdateStatus(ctx context.Context, subscriptionID int64, status string) error
	UpdateNotes(ctx context.Context, subscriptionID int64, notes string) error

	ActivateWindows(ctx context.Context, id int64, start time.Time) error
	ResetDailyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	ResetWeeklyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	ResetMonthlyUsage(ctx context.Context, id int64, newWindowStart time.Time) error
	IncrementUsage(ctx context.Context, id int64, costUSD float64) error
	IncrementRequestUsage(ctx context.Context, id int64, count int64) error

	// ListActiveRefsByPlanID 查询引用指定 plan 的所有活跃订阅引用，用于 plan 更新后缓存失效。
	ListActiveRefsByPlanID(ctx context.Context, planID int64) ([]SubscriptionRef, error)
	// CountActiveByPlanID 统计引用指定 plan 的活跃订阅数，用于 plan 删除保护。
	CountActiveByPlanID(ctx context.Context, planID int64) (int64, error)
	// CountByPlanID counts all non-deleted subscriptions referencing a plan, for plan deletion protection.
	CountByPlanID(ctx context.Context, planID int64) (int64, error)

	// UpdatePlanID updates the plan_id field of a subscription without touching other fields.
	UpdatePlanID(ctx context.Context, subscriptionID int64, planID *int64) error

	BatchUpdateExpiredStatus(ctx context.Context) (int64, error)
}
