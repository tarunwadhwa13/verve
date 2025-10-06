package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresBadgeRepository struct {
	DB *sql.DB
}

func NewPostgresBadgeRepository(db *sql.DB) repository.BadgeRepository {
	return &postgresBadgeRepository{DB: db}
}

func (r *postgresBadgeRepository) Create(badge *models.Badge) error {
	return r.DB.QueryRow(`
		INSERT INTO badges (name, description, icon_url, points, created_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`,
		badge.Name, badge.Description, badge.IconURL, badge.Points, badge.CreatedBy, badge.IsActive,
	).Scan(&badge.ID, &badge.CreatedAt)
}

func (r *postgresBadgeRepository) Update(badge *models.Badge) error {
	_, err := r.DB.Exec(`
		UPDATE badges
		SET name = $1, description = $2, icon_url = $3, points = $4, is_active = $5
		WHERE id = $6`,
		badge.Name, badge.Description, badge.IconURL, badge.Points, badge.IsActive, badge.ID,
	)
	return err
}

func (r *postgresBadgeRepository) FindByID(id int) (*models.Badge, error) {
	badge := &models.Badge{}
	err := r.DB.QueryRow(`
		SELECT id, name, description, icon_url, points, created_at, created_by, is_active
		FROM badges WHERE id = $1`,
		id,
	).Scan(
		&badge.ID, &badge.Name, &badge.Description, &badge.IconURL,
		&badge.Points, &badge.CreatedAt, &badge.CreatedBy, &badge.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return badge, nil
}

func (r *postgresBadgeRepository) FindAll(includeInactive bool) ([]*models.Badge, error) {
	query := `
		SELECT id, name, description, icon_url, points, created_at, created_by, is_active
		FROM badges`
	if !includeInactive {
		query += " WHERE is_active = true"
	}

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []*models.Badge
	for rows.Next() {
		badge := &models.Badge{}
		err := rows.Scan(
			&badge.ID, &badge.Name, &badge.Description, &badge.IconURL,
			&badge.Points, &badge.CreatedAt, &badge.CreatedBy, &badge.IsActive,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, nil
}

func (r *postgresBadgeRepository) Delete(id int) error {
	_, err := r.DB.Exec("UPDATE badges SET is_active = false WHERE id = $1", id)
	return err
}

type postgresAchievementRuleRepository struct {
	DB *sql.DB
}

func NewPostgresAchievementRuleRepository(db *sql.DB) repository.AchievementRuleRepository {
	return &postgresAchievementRuleRepository{DB: db}
}

func (r *postgresAchievementRuleRepository) Create(rule *models.AchievementRule) error {
	return r.DB.QueryRow(`
		INSERT INTO achievement_rules (badge_id, rule_type, condition_value, created_by, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`,
		rule.BadgeID, rule.RuleType, rule.ConditionValue, rule.CreatedBy, rule.IsActive,
	).Scan(&rule.ID, &rule.CreatedAt)
}

func (r *postgresAchievementRuleRepository) Update(rule *models.AchievementRule) error {
	_, err := r.DB.Exec(`
		UPDATE achievement_rules
		SET rule_type = $1, condition_value = $2, is_active = $3
		WHERE id = $4`,
		rule.RuleType, rule.ConditionValue, rule.IsActive, rule.ID,
	)
	return err
}

func (r *postgresAchievementRuleRepository) FindByBadgeID(badgeID int) ([]*models.AchievementRule, error) {
	rows, err := r.DB.Query(`
		SELECT id, badge_id, rule_type, condition_value, created_at, created_by, is_active
		FROM achievement_rules
		WHERE badge_id = $1`,
		badgeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.AchievementRule
	for rows.Next() {
		rule := &models.AchievementRule{}
		err := rows.Scan(
			&rule.ID, &rule.BadgeID, &rule.RuleType, &rule.ConditionValue,
			&rule.CreatedAt, &rule.CreatedBy, &rule.IsActive,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *postgresAchievementRuleRepository) FindAll(includeInactive bool) ([]*models.AchievementRule, error) {
	query := `
		SELECT id, badge_id, rule_type, condition_value, created_at, created_by, is_active
		FROM achievement_rules`
	if !includeInactive {
		query += " WHERE is_active = true"
	}

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.AchievementRule
	for rows.Next() {
		rule := &models.AchievementRule{}
		err := rows.Scan(
			&rule.ID, &rule.BadgeID, &rule.RuleType, &rule.ConditionValue,
			&rule.CreatedAt, &rule.CreatedBy, &rule.IsActive,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *postgresAchievementRuleRepository) Delete(id int) error {
	_, err := r.DB.Exec("UPDATE achievement_rules SET is_active = false WHERE id = $1", id)
	return err
}

type postgresUserBadgeRepository struct {
	DB *sql.DB
}

func NewPostgresUserBadgeRepository(db *sql.DB) repository.UserBadgeRepository {
	return &postgresUserBadgeRepository{DB: db}
}

func (r *postgresUserBadgeRepository) Award(userBadge *models.UserBadge) error {
	return r.DB.QueryRow(`
		INSERT INTO user_badges (user_id, badge_id, awarded_by)
		VALUES ($1, $2, $3)
		RETURNING id, awarded_at`,
		userBadge.UserID, userBadge.BadgeID, userBadge.AwardedBy,
	).Scan(&userBadge.ID, &userBadge.AwardedAt)
}

func (r *postgresUserBadgeRepository) FindByUserID(userID int) ([]*models.UserBadge, error) {
	rows, err := r.DB.Query(`
		SELECT id, user_id, badge_id, awarded_at, awarded_by
		FROM user_badges
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []*models.UserBadge
	for rows.Next() {
		ub := &models.UserBadge{}
		err := rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.AwardedAt, &ub.AwardedBy)
		if err != nil {
			return nil, err
		}
		userBadges = append(userBadges, ub)
	}
	return userBadges, nil
}

func (r *postgresUserBadgeRepository) HasBadge(userID, badgeID int) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_badges
			WHERE user_id = $1 AND badge_id = $2
		)`,
		userID, badgeID,
	).Scan(&exists)
	return exists, err
}

func (r *postgresUserBadgeRepository) FindByBadgeID(badgeID int) ([]*models.UserBadge, error) {
	rows, err := r.DB.Query(`
		SELECT ub.id, ub.user_id, ub.badge_id, ub.awarded_at, ub.awarded_by
		FROM user_badges ub
		WHERE ub.badge_id = $1
		ORDER BY ub.awarded_at DESC`,
		badgeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges []*models.UserBadge
	for rows.Next() {
		ub := &models.UserBadge{}
		err := rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.AwardedAt, &ub.AwardedBy)
		if err != nil {
			return nil, err
		}
		userBadges = append(userBadges, ub)
	}
	return userBadges, nil
}
