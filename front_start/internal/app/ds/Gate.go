package ds

type Gate struct {
	ID_gate     uint   `gorm:"column:id_gate;primaryKey"`
	Title       string `gorm:"not null;unique"`
	Description string `gorm:"not null"`
	Status      bool   `gorm:"default:true"`
	Image       string

	// subject area
	FullInfo   string
	IsEditable bool `gorm:"default:false"`
	TheAxis    string

	// связи
	//Gate DegreesToGates `gorm:"foreignKey:ID_gate;references:ID_gate"`
}
