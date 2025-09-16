package repository

import (
	"errors"
	"fmt"
	"front_start/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetDraftTask(userID uint) (*ds.Task, error) {
	var task ds.Task

	err := r.db.Where("id_user = ? AND tesk_status = ?", userID, ds.StatusDraft).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) CreateTask(task *ds.Task) error {
	return r.db.Create(task).Error
}

// AddGateToTask добавляет гейт в задачу (создает запись в таблице m-m).
func (r *Repository) AddGateToTask(taskID, gateID uint) error {
	// Проверяем, нет ли уже такого фактора в заявке, чтобы избежать дублей
	var count int64
	r.db.Model(&ds.DegreesToGates{}).Where("id_task = ? AND id_gate = ?", taskID, gateID).Count(&count)
	if count > 0 {
		return errors.New("gate already in task")
	}

	link := ds.DegreesToGates{
		ID_task: taskID,
		ID_gate: gateID,
		Degrees: 0,
	}
	return r.db.Create(&link).Error
}

// GetTaskWithGates получает задачу со всеми связанными гейтами и их углами поворота.
func (r *Repository) GetTaskWithGates(taskID uint) (*ds.Task, error) {
	var task ds.Task

	// Используем Preload для загрузки связанных данных через связующую таблицу
	err := r.db.
		Preload("Task.Gate").
		First(&task, taskID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("error fetching task: %v", err)
	}

	// Проверяем, что задача не удалена (адаптируйте под ваши статусы)
	if task.TeskStatus == ds.StatusDeleted {
		return nil, errors.New("task not found or has been deleted")
	}

	return &task, nil
}

// LogicallyDeleteFrax выполняет логическое удаление заявки через чистый SQL UPDATE.
func (r *Repository) LogicallyDeleteFrax(taskID uint) error {
	// Используем Exec для выполнения "сырого" SQL-запроса
	result := r.db.Exec("UPDATE task SET tesk_status = ? WHERE id_task = ?", ds.StatusDeleted, taskID)
	return result.Error
}
