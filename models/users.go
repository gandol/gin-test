package models

type Users struct {
	ID       uint   `json:"id" gorm:"primary_key;unique"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
