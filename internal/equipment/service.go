package equipment

import (
	"errors"
	"math"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

// Equipment methods
func (s *Service) CreateEquipment(userID uint, name, equipmentType, brand, model string, weight float64, notes string) (*Equipment, error) {
	if name == "" || equipmentType == "" {
		return nil, errors.New("name and type are required")
	}

	equipment := &Equipment{
		UserID: userID,
		Name:   name,
		Type:   equipmentType,
		Brand:  brand,
		Model:  model,
		Weight: weight,
		Notes:  notes,
		Active: true,
	}

	if err := s.repo.CreateEquipment(equipment); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *Service) GetUserEquipment(userID uint) ([]Equipment, error) {
	return s.repo.GetUserEquipment(userID)
}

func (s *Service) GetEquipmentByID(id, userID uint) (*Equipment, error) {
	return s.repo.GetEquipmentByID(id, userID)
}

func (s *Service) UpdateEquipment(id, userID uint, updates map[string]interface{}) (*Equipment, error) {
	equipment, err := s.repo.GetEquipmentByID(id, userID)
	if err != nil {
		return nil, err
	}

	if name, ok := updates["name"].(string); ok && name != "" {
		equipment.Name = name
	}
	if brand, ok := updates["brand"].(string); ok {
		equipment.Brand = brand
	}
	if model, ok := updates["model"].(string); ok {
		equipment.Model = model
	}
	if weight, ok := updates["weight"].(float64); ok {
		equipment.Weight = weight
	}
	if notes, ok := updates["notes"].(string); ok {
		equipment.Notes = notes
	}
	if active, ok := updates["active"].(bool); ok {
		equipment.Active = active
	}

	if err := s.repo.UpdateEquipment(equipment); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *Service) DeleteEquipment(id, userID uint) error {
	return s.repo.DeleteEquipment(id, userID)
}

// Training Zones calculation
func (s *Service) CalculateHRZones(maxHR, restingHR int) *HRZones {
	if maxHR <= 0 {
		return nil
	}

	// Karvonen formula (Heart Rate Reserve)
	hrr := maxHR - restingHR

	return &HRZones{
		Zone1: TrainingZone{
			ID:    1,
			Name:  "Recovery",
			Min:   restingHR + int(math.Round(float64(hrr)*0.50)),
			Max:   restingHR + int(math.Round(float64(hrr)*0.60)),
			Color: "#4285F4", // Blue
		},
		Zone2: TrainingZone{
			ID:    2,
			Name:  "Endurance",
			Min:   restingHR + int(math.Round(float64(hrr)*0.60)),
			Max:   restingHR + int(math.Round(float64(hrr)*0.70)),
			Color: "#34A853", // Green
		},
		Zone3: TrainingZone{
			ID:    3,
			Name:  "Tempo",
			Min:   restingHR + int(math.Round(float64(hrr)*0.70)),
			Max:   restingHR + int(math.Round(float64(hrr)*0.80)),
			Color: "#FBBC04", // Yellow
		},
		Zone4: TrainingZone{
			ID:    4,
			Name:  "Threshold",
			Min:   restingHR + int(math.Round(float64(hrr)*0.80)),
			Max:   restingHR + int(math.Round(float64(hrr)*0.90)),
			Color: "#EA4335", // Red
		},
		Zone5: TrainingZone{
			ID:    5,
			Name:  "VO2 Max",
			Min:   restingHR + int(math.Round(float64(hrr)*0.90)),
			Max:   maxHR,
			Color: "#A61C00", // Dark Red
		},
	}
}

func (s *Service) CalculatePowerZones(ftp int) *PowerZones {
	if ftp <= 0 {
		return nil
	}

	return &PowerZones{
		Zone1: TrainingZone{
			ID:    1,
			Name:  "Active Recovery",
			Min:   0,
			Max:   int(math.Round(float64(ftp) * 0.55)),
			Color: "#4285F4", // Blue
		},
		Zone2: TrainingZone{
			ID:    2,
			Name:  "Endurance",
			Min:   int(math.Round(float64(ftp) * 0.55)),
			Max:   int(math.Round(float64(ftp) * 0.75)),
			Color: "#34A853", // Green
		},
		Zone3: TrainingZone{
			ID:    3,
			Name:  "Sweet Spot",
			Min:   int(math.Round(float64(ftp) * 0.75)),
			Max:   int(math.Round(float64(ftp) * 0.90)),
			Color: "#FBBC04", // Yellow
		},
		Zone4: TrainingZone{
			ID:    4,
			Name:  "Threshold",
			Min:   int(math.Round(float64(ftp) * 0.90)),
			Max:   int(math.Round(float64(ftp) * 1.05)),
			Color: "#EA4335", // Red
		},
		Zone5: TrainingZone{
			ID:    5,
			Name:  "VO2 Max",
			Min:   int(math.Round(float64(ftp) * 1.05)),
			Max:   int(math.Round(float64(ftp) * 1.20)),
			Color: "#A61C00", // Dark Red
		},
		Zone6: TrainingZone{
			ID:    6,
			Name:  "Anaerobic",
			Min:   int(math.Round(float64(ftp) * 1.20)),
			Max:   int(math.Round(float64(ftp) * 1.50)),
			Color: "#800080", // Purple
		},
		Zone7: TrainingZone{
			ID:    7,
			Name:  "Neuromuscular",
			Min:   int(math.Round(float64(ftp) * 1.50)),
			Max:   int(math.Round(float64(ftp) * 2.00)),
			Color: "#FF1744", // Bright Red
		},
	}
}