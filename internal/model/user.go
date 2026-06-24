package model

type UserID uint64

type User struct {
	ID       UserID
	Name     string
	Email    string
	Password string
}
