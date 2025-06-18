package model

type User struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username string
	Email    string `gorm:"uniqueIndex"`
	Password string
}
