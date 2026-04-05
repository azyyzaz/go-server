package user

import "time"

type User struct {
	ID        uint
	Username  string
	Password  string // bcrypt hash
	Name      string
	Email     string
	CreatedAt time.Time
}
