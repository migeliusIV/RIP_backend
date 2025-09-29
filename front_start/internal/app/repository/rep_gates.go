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
