package repository

import "verve/internal/models"

type BadgeRepository interface {
	Create(badge *models.Badge) error
	Update(badge *models.Badge) error
	FindByID(id int) (*models.Badge, error)
	FindAll(includeInactive bool) ([]*models.Badge, error)
	Delete(id int) error
}

type AchievementRuleRepository interface {
	Create(rule *models.AchievementRule) error
	Update(rule *models.AchievementRule) error
	FindByBadgeID(badgeID int) ([]*models.AchievementRule, error)
	FindAll(includeInactive bool) ([]*models.AchievementRule, error)
	Delete(id int) error
}

type UserBadgeRepository interface {
	Award(userBadge *models.UserBadge) error
	FindByUserID(userID int) ([]*models.UserBadge, error)
	FindByBadgeID(badgeID int) ([]*models.UserBadge, error)
	HasBadge(userID, badgeID int) (bool, error)
}
