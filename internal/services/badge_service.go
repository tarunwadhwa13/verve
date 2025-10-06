package services

import (
	"encoding/json"
	"errors"
	"verve/internal/models"
	"verve/internal/repository"
)

type BadgeService struct {
	badgeRepo     repository.BadgeRepository
	ruleRepo      repository.AchievementRuleRepository
	userBadgeRepo repository.UserBadgeRepository
}

// GetBadgeHolders returns all users who have been awarded a specific badge
func (s *BadgeService) GetBadgeHolders(badgeID int) ([]*models.UserBadge, error) {
	return s.userBadgeRepo.FindByBadgeID(badgeID)
}

func NewBadgeService(
	badgeRepo repository.BadgeRepository,
	ruleRepo repository.AchievementRuleRepository,
	userBadgeRepo repository.UserBadgeRepository,
) *BadgeService {
	return &BadgeService{
		badgeRepo:     badgeRepo,
		ruleRepo:      ruleRepo,
		userBadgeRepo: userBadgeRepo,
	}
}

// CreateBadge creates a new badge with optional achievement rules
func (s *BadgeService) CreateBadge(
	name, description, iconURL string,
	points int,
	createdBy int,
	rules []models.AchievementRule,
) (*models.Badge, error) {
	badge := &models.Badge{
		Name:        name,
		Description: description,
		IconURL:     iconURL,
		Points:      points,
		CreatedBy:   createdBy,
		IsActive:    true,
	}

	if err := s.badgeRepo.Create(badge); err != nil {
		return nil, err
	}

	// Create achievement rules if provided
	for i := range rules {
		rules[i].BadgeID = badge.ID
		rules[i].CreatedBy = createdBy
		rules[i].IsActive = true
		if err := s.ruleRepo.Create(&rules[i]); err != nil {
			return nil, err
		}
	}

	return badge, nil
}

// UpdateBadge updates an existing badge and its rules
func (s *BadgeService) UpdateBadge(
	id int,
	name, description, iconURL *string,
	points *int,
	isActive *bool,
) (*models.Badge, error) {
	badge, err := s.badgeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		badge.Name = *name
	}
	if description != nil {
		badge.Description = *description
	}
	if iconURL != nil {
		badge.IconURL = *iconURL
	}
	if points != nil {
		badge.Points = *points
	}
	if isActive != nil {
		badge.IsActive = *isActive
	}

	if err := s.badgeRepo.Update(badge); err != nil {
		return nil, err
	}

	return badge, nil
}

// GetBadge retrieves a badge by ID along with its rules
func (s *BadgeService) GetBadge(id int) (*models.Badge, []*models.AchievementRule, error) {
	badge, err := s.badgeRepo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	rules, err := s.ruleRepo.FindByBadgeID(id)
	if err != nil {
		return nil, nil, err
	}

	return badge, rules, nil
}

// ListBadges retrieves all badges, optionally including inactive ones
func (s *BadgeService) ListBadges(includeInactive bool) ([]*models.Badge, error) {
	return s.badgeRepo.FindAll(includeInactive)
}

// AwardBadge manually awards a badge to a user
func (s *BadgeService) AwardBadge(userID, badgeID int, awardedBy *int) error {
	// Check if badge exists and is active
	badge, err := s.badgeRepo.FindByID(badgeID)
	if err != nil {
		return err
	}
	if !badge.IsActive {
		return errors.New("badge is not active")
	}

	// Check if user already has the badge
	hasBadge, err := s.userBadgeRepo.HasBadge(userID, badgeID)
	if err != nil {
		return err
	}
	if hasBadge {
		return errors.New("user already has this badge")
	}

	userBadge := &models.UserBadge{
		UserID:    userID,
		BadgeID:   badgeID,
		AwardedBy: awardedBy,
	}

	return s.userBadgeRepo.Award(userBadge)
}

// GetUserBadges retrieves all badges awarded to a user
func (s *BadgeService) GetUserBadges(userID int) ([]*models.UserBadge, error) {
	return s.userBadgeRepo.FindByUserID(userID)
}

// ValidateAchievementRule validates the rule condition format based on rule type
func (s *BadgeService) ValidateAchievementRule(rule *models.AchievementRule) error {
	switch rule.RuleType {
	case models.RuleTypeTransactionCount:
		var condition models.TransactionCountCondition
		if err := json.Unmarshal(rule.ConditionValue, &condition); err != nil {
			return errors.New("invalid transaction count condition format")
		}
		if condition.MinTransactions <= 0 {
			return errors.New("min_transactions must be positive")
		}

	case models.RuleTypeTransferAmount:
		var condition models.TransferAmountCondition
		if err := json.Unmarshal(rule.ConditionValue, &condition); err != nil {
			return errors.New("invalid transfer amount condition format")
		}
		if condition.MinAmount <= 0 {
			return errors.New("min_amount must be positive")
		}

	case models.RuleTypeConsecutiveDays:
		var condition models.ConsecutiveDaysCondition
		if err := json.Unmarshal(rule.ConditionValue, &condition); err != nil {
			return errors.New("invalid consecutive days condition format")
		}
		if condition.Days <= 0 {
			return errors.New("days must be positive")
		}

	default:
		return errors.New("unsupported rule type")
	}

	return nil
}
