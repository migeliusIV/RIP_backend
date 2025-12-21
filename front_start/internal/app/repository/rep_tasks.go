package repository

import (
	"errors"
	"fmt"
	"front_start/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetDraftTask(userID uint) (*ds.QuantumTask, error) {
	var task ds.QuantumTask

	err := r.db.Where("id_user = ? AND task_status = ?", userID, ds.StatusDraft).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) CreateTask(task *ds.QuantumTask) error {
	// Set creation date to current time if not already set
	if task.CreationDate.IsZero() {
		task.CreationDate = time.Now()
	}
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
	}
	return r.db.Create(&link).Error
}

// GetTaskWithGates получает задачу со всеми связанными гейтами и их углами поворота.
func (r *Repository) GetTaskWithGates(taskID uint) (*ds.QuantumTask, error) {
	var task ds.QuantumTask

	// Используем Preload для загрузки связанных данных через связующую таблицу
	err := r.db.
		Preload("GatesDegrees.Gate").
		First(&task, taskID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("error fetching task: %v", err)
	}

	// Проверяем, что задача не удалена (адаптируйте под ваши статусы)
	if task.TaskStatus == ds.StatusDeleted {
		return nil, errors.New("task not found or has been deleted")
	}

	return &task, nil
}

// LogicallyDeleteTask выполняет логическое удаление заявки через чистый SQL UPDATE.
func (r *Repository) LogicallyDeleteTask(taskID uint) error {
	// Используем Exec для выполнения "сырого" SQL-запроса
	result := r.db.Exec("UPDATE quantum_tasks SET task_status = ? WHERE id_task = ?", ds.StatusDeleted, taskID)
	return result.Error
}

// ---- JSON task repo ----

// В репозитории сделайте оба метода возвращать один тип
func (r *Repository) ListTasks(status, from, to string) ([]*ds.QuantumTask, error) {
	var tasks []*ds.QuantumTask
	q := r.db.Preload("GatesDegrees.Gate").Model(&ds.QuantumTask{})
	if status != "" {
		q = q.Where("task_status = ?", status)
	} else {
		q = q.Where("task_status NOT IN ?", []string{ds.StatusDeleted, ds.StatusDraft})
	}
	if from != "" {
		q = q.Where("creation_date >= ?", from)
	}
	if to != "" {
		q = q.Where("creation_date <= ?", to)
	}
	if err := q.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) ListTasksByUser(userID uint, status, from, to string) ([]*ds.QuantumTask, error) {
	var tasks []*ds.QuantumTask
	// Включим логирование SQL
	//r.db = r.db.Debug()
	query := r.db.Preload("GatesDegrees.Gate").Where("id_user = ?", userID)

	if status != "" {
		query = query.Where("task_status = ?", status)
	} else {
		query = query.Where("task_status NOT IN ?", []string{ds.StatusDeleted, ds.StatusDraft})
	}
	if from != "" {
		query = query.Where("creation_date >= ?", from)
	}
	if to != "" {
		query = query.Where("creation_date <= ?", to)
	}

	err := query.Order("creation_date DESC").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(id uint, description string) (*ds.QuantumTask, error) {
	var task ds.QuantumTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	task.TaskDescription = description
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) FormTask(id uint, formDate time.Time) (*ds.QuantumTask, error) {
	var task ds.QuantumTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	if task.TaskStatus != ds.StatusDraft {
		return nil, errors.New("only draft can be formed")
	}
	task.TaskStatus = ds.StatusFormed
	task.FormedDate = formDate
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) ResolveTask(id uint, moderatorID uint, action string, resolvedAt time.Time) (*ds.QuantumTask, error) {
	var task ds.QuantumTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	if task.TaskStatus != ds.StatusFormed {
		return nil, errors.New("only formed can be resolved")
	}
	if action == "complete" {
		task.TaskStatus = ds.StatusCompleted
	} else if action == "reject" {
		task.TaskStatus = ds.StatusRejected
	} else {
		return nil, errors.New("invalid action")
	}
	task.ID_moderator = &moderatorID
	task.ConclusionDate = resolvedAt
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) DeleteTask(id uint) error {
	return r.db.Model(&ds.QuantumTask{}).Where("id_task = ?", id).Update("task_status", ds.StatusDeleted).Error
}

// Helper repos for m-m
func (r *Repository) RemoveGateFromTask(taskID, gateID uint) error {
	return r.db.Where("id_task = ? AND id_gate = ?", taskID, gateID).Delete(&ds.DegreesToGates{}).Error
}

func (r *Repository) UpdateDegrees(taskID, gateID uint, degrees float32) error {
	return r.db.Model(&ds.DegreesToGates{}).Where("id_task = ? AND id_gate = ?", taskID, gateID).Update("degrees", degrees).Error
}

// UpdateQuantumTaskResult обновляет амплитуды в квантовой задаче после расчёта во внешнем сервисе
func (r *Repository) UpdateQuantumTaskResult(taskID uint, k0, k1 float32) error {
	return r.db.Model(&ds.QuantumTask{}).
		Where("id_task = ?", taskID).
		Updates(map[string]interface{}{
			"res_koeff_0": k0,
			"res_koeff_1": k1,
		}).Error
}
