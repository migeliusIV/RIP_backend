package ds

type Users struct {
	ID_user  uint   `gorm:"column:id_user;primaryKey"`
	Login    string `gorm:"column:login;not null;size:255;unique"`
	Password string `gorm:"column:password;size:255;not null"`
	IsAdmin  bool   `gorm:"column:is_admin;default:false"`
	// Связь для приложения
	//Tasks []QuantumTask `gorm:"foreignKey:ID_user;references:ID_user" json:"-"`
}
