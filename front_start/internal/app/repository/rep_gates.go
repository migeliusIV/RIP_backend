package repository

import (
	"fmt"
	"front_start/internal/app/ds"
)

func (r *Repository) GetGates() ([]ds.Gate, error) {
	var gates []ds.Gate

	err := r.db.Find(&gates).Error
	if err != nil {
		return nil, err
	}

	if len(gates) == 0 {
		return nil, fmt.Errorf("gates not found")
	}
	return gates, nil
}

func (r *Repository) GetGatesByName(title string) ([]ds.Gate, error) {
	var gates []ds.Gate
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&gates).Error
	if err != nil {
		return nil, err
	}
	return gates, nil
}

func (r *Repository) GetGateByID(id int) (*ds.Gate, error) {
	var gate ds.Gate
	err := r.db.First(&gate, id).Error
	if err != nil {
		return nil, err
	}
	return &gate, nil
}

// ---- JSON API helpers for services ----

func (r *Repository) ListServices(title string) ([]ds.Gate, error) {
    var gates []ds.Gate
    q := r.db
    if title != "" {
        q = q.Where("title ILIKE ?", "%"+title+"%")
    }
    if err := q.Find(&gates).Error; err != nil {
        return nil, err
    }
    return gates, nil
}

func (r *Repository) CreateService(g *ds.Gate) error {
    return r.db.Create(g).Error
}

func (r *Repository) UpdateService(id uint, title, description, fullInfo, theAxis string, status *bool) (*ds.Gate, error) {
    var gate ds.Gate
    if err := r.db.First(&gate, id).Error; err != nil {
        return nil, err
    }
    if title != "" {
        gate.Title = title
    }
    if description != "" {
        gate.Description = description
    }
    if fullInfo != "" {
        gate.FullInfo = fullInfo
    }
    if theAxis != "" {
        gate.TheAxis = theAxis
    }
    if status != nil {
        gate.Status = *status
    }
    if err := r.db.Save(&gate).Error; err != nil {
        return nil, err
    }
    return &gate, nil
}

func (r *Repository) DeleteService(id uint) error {
    return r.db.Delete(&ds.Gate{}, id).Error
}

func (r *Repository) SetServiceImage(id uint, image string) error {
    return r.db.Model(&ds.Gate{}).Where("id_gate = ?", id).Update("image", image).Error
}