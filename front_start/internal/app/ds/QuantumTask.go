package ds

import "time"

type QuantumTask struct {
	ID_task        uint      `gorm:"column:id_task;primaryKey"`
	TaskStatus     string    `gorm:"column:task_status;size:255;not null;default:черновик"`
	CreationDate   time.Time `gorm:"column:creation_date;not null"`
	ID_user        uint      `gorm:"column:id_user;not null"`
	ConclusionDate time.Time `gorm:"column:conclusion_date"`
	ID_moderator   uint      `gorm:"column:id_moderator"`

	TaskDescription string  `gorm:"column:task_description"`
	Res_koeff_0     float32 `gorm:"column:res_koeff_0;default:1"`
	Res_koeff_1     float32 `gorm:"column:res_koeff_1;default:0"`

	// Несущая связь
	GatesDegrees []DegreesToGates `gorm:"foreignKey:ID_task;references:ID_task"`
	User         Users            `gorm:"foreignKey:ID_user;references:ID_user"`
	Moderator    Users            `gorm:"foreignKey:ID_moderator;references:ID_user"`
}
