package ds

type Gate struct {
	ID_gate     uint    `gorm:"column:id_gate;primaryKey"`
	Title       string  `gorm:"column:title;size:255;not null;default:gate-no-name;unique"`
	Description string  `gorm:"column:description;not null"`
	Status      bool    `gorm:"column:status; default:true"`
	Image       *string `gorm:"column:image"`

	// subject area
	FullInfo string  `gorm:"column:full_info;not null"`
	TheAxis  *string `gorm:"column:the_axis"`
	/*
		Удалил согласно правкам к ЛР 1
		IsEditable bool `gorm:"default:false"`
	*/

	// Несущая связь
	Degrees []DegreesToGates `gorm:"foreignKey:ID_gate;references:ID_gate"`
}
