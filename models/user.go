package models

import "gorm.io/gorm"

// User => Struct yang mewakili tabel "users" di database
type User struct {
	gorm.Model          // Menambahkan kolom ID, CreatedAt, UpdatedAt, DeletedAt
	Name  string `json:"name"`   // Kolom "name"
	Email string `json:"email"`  // Kolom "email"
}
