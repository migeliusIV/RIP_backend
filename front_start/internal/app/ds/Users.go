package ds

type Users struct {
	ID_user  uint   `gorm:"primaryKey"`
	Login    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`

	// Связи
	Tasks []Task `gorm:"foreignKey:ID_user" json:"-"`
}

/*Работает корректно*/
