package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/moceviciusda/pokeCLIpse-server/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		dbUser.ID,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
		dbUser.Username,
		dbUser.Password,
	}
}