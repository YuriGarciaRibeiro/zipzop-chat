package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string   `json:"name,omitempty" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Salt         string    `json:"-" db:"salt"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Para registro de novos usuários
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name,omitempty" validate:"omitempty,e164"`
	Password string `json:"password" validate:"required,min=8"`
}

// Para resposta pública (sem campos sensíveis)
type PublicUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string   `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ToPublic converte User para PublicUser
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}
