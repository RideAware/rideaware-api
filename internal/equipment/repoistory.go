package equipment

import (
	"errors"
	"rideaware/pkg/database"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateEquipment(equipment *Equipment) error {
	return database.DB.Create(equipment).Error
}

func (r *Repository) GetEquipmentByID(id, userID uint) (*Equipment, error) {
	var equipment Equipment
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		First(&equipment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("equipment not found")
		}
		return nil, err
	}
	return &equipment, nil
}

func (r *Repository) GetUserEquipment(userID uint) ([]Equipment, error) {
	var equipment []Equipment
	if err := database.DB.Where("user_id = ?", userID).
		Find(&equipment).Error; err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *Repository) GetActiveEquipment(userID uint) ([]Equipment, error) {
	var equipment []Equipment
	if err := database.DB.Where("user_id = ? AND active = ?", userID, true).
		Find(&equipment).Error; err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *Repository) UpdateEquipment(equipment *Equipment) error {
	return database.DB.Save(equipment).Error
}

func (r *Repository) DeleteEquipment(id, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&Equipment{}).Error
}