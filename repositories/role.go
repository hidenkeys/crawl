package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository struct {
	BaseRepository[models.Role]
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{
		BaseRepository: BaseRepository[models.Role]{DB: db},
	}
}

func (r *RoleRepository) AssignRoleToUser(userID, roleID uuid.UUID) error {
	return r.DB.
		Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING", userID, roleID).
		Error
}

func (r *RoleRepository) RemoveRoleFromUser(userID, roleID uuid.UUID) error {
	return r.DB.
		Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).
		Error
}

func (r *RoleRepository) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).
		Error
	return roles, err
}
