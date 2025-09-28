package ds

type DegreesToGates struct {
	ID_gate uint `gorm:"primaryKey;column:id_gate;not null"`
	ID_task uint `gorm:"primaryKey;column:id_task;not null"`
	Degrees *float32
	// Связи, необходимые для работы сервиса
	Gate Gate        `gorm:"foreignKey:ID_gate;references:ID_gate"`
	Task QuantumTask `gorm:"foreignKey:ID_task;references:ID_task"`
}
