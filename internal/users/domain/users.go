package domain

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	DeletedAt time.Time `bson:"deleted_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	SoftDelete(id string) error
	Delete(id string) error
	List(page, limit int, nameFilter, emailFilter string) ([]*User, int, error)
}

type UserUsecase interface {
	GetProfile(userID string) (*User, error)
	UpdateProfile(userID, name string) (*User, error)
	DeleteProfile(userID string) error
	ListUsers(page, limit int, nameFilter, emailFilter string) ([]*User, int, error)
}
