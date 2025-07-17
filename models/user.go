package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"id"`
	Nama      string    `json:"nama" bson:"nama"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"` // ✅ Tambahkan ini
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
