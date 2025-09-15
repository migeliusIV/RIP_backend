package ds

type DegreesToGates struct {
	ID_gate uint `gorm:"primaryKey;column:id_gate;not null"`
	ID_task uint `gorm:"primaryKey;column:id_task;not null"`
	Degrees float32

	//Связи
	Gates Gate `gorm:"foreignKey:ID_gate;references:ID_gate" json:"-"`
	Tasks Task `gorm:"foreignKey:ID_task;references:ID_task" json:"-"`
}
