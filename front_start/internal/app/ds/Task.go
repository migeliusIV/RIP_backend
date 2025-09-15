package ds

import "time"

type Task struct {
	ID_task        uint      `gorm:"column:id_task;primaryKey"`
	TeskStatus     string    `gorm:"column:tesk_status;not null;default:черновик"`
	CreationDate   time.Time `gorm:"column:creation_date;not null"`
	ID_user        uint      `gorm:"column:id_user;not null"`
	ConclusionDate time.Time `gorm:"column:conclusion_date"`
	ID_moderator   uint      `gorm:"column:id_moderator"`

	// subject area
	TaskDescription string `gorm:"column:task_description"`
	Result          string

	// Связи
	Task DegreesToGates `gorm:"foreignKey:ID_task;references:ID_task"`
	// User      Users `gorm:"foreignKey:ID_user;references:ID_user"`
	// Moderator Users `gorm:"foreignKey:ID_moderator;references:ID_user"`
}
