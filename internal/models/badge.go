package models

import (
	"encoding/json"
	"time"
)

type Badge struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IconURL     string    `json:"icon_url"`
	Points      int       `json:"points"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   int       `json:"created_by"`
	IsActive    bool      `json:"is_active"`
}

type AchievementRule struct {
	ID             int             `json:"id"`
	BadgeID        int             `json:"badge_id"`
	RuleType       string          `json:"rule_type"`
	ConditionValue json.RawMessage `json:"condition_value"`
	CreatedAt      time.Time       `json:"created_at"`
	CreatedBy      int             `json:"created_by"`
	IsActive       bool            `json:"is_active"`
}

type UserBadge struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BadgeID   int       `json:"badge_id"`
	AwardedAt time.Time `json:"awarded_at"`
	AwardedBy *int      `json:"awarded_by"` // Pointer to allow NULL for system-awarded badges
}

// Rule types and condition structures
const (
	RuleTypeTransactionCount = "transaction_count"
	RuleTypeTransferAmount   = "transfer_amount"
	RuleTypeConsecutiveDays  = "consecutive_days"
)

type TransactionCountCondition struct {
	MinTransactions int    `json:"min_transactions"`
	TimeFrame       string `json:"time_frame"` // e.g., "24h", "7d", "30d"
}

type TransferAmountCondition struct {
	MinAmount int64  `json:"min_amount"`
	TimeFrame string `json:"time_frame"`
	Currency  string `json:"currency"`
	Direction string `json:"direction"` // "sent", "received", or "total"
}

type ConsecutiveDaysCondition struct {
	Days          int    `json:"days"`
	ActivityType  string `json:"activity_type"` // e.g., "login", "transfer"
	MinActivities int    `json:"min_activities"`
}
