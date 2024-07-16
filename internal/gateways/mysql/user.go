package mysql

import "time"

type User struct {
	ID             int64     `db:"id"`
	Email          string    `db:"email"`
	ProfilePicture string    `db:"profile_picture"`
	Name           string    `db:"name"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
