package models

type User struct {
	ID         string `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	CreatedAt  int64  `json:"created_at" db:"created_at"`
	UpdatedAt  int64  `json:"updated_at" db:"updated_at"`
}
