package ds

type DegreesToGates struct {
	ID_gate uint `gorm:"primaryKey;column:id_gate;not null"`
	ID_task uint `gorm:"primaryKey;column:id_task;not null"`
	Degrees float32
}
