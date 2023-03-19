package models

type LoginData struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	IsAdmin  bool   `gorm:"not null;default:false" json:"is_admin"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
