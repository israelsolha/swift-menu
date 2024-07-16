package mysql

import (
	"database/sql"
	"swift-menu-session/internal/domain/entities"
	"time"
)

type UserGateway struct {
	db *sql.DB
}

func NewUserGateway(db *sql.DB) *UserGateway {
	return &UserGateway{
		db: db,
	}
}

func (u *UserGateway) CreateUser(user entities.User) (entities.User, error) {
	result, err := u.db.Exec("INSERT INTO users (email, profile_picture, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Email, user.ProfilePicture, user.Name, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return entities.User{}, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return entities.User{}, err
	}
	user.ID = &userID
	return user, nil
}

func (u *UserGateway) GetUserByID(id int) (entities.User, error) {
	var userDB User
	err := u.db.QueryRow("SELECT id, email, profile_picture, name FROM users WHERE id = ?", id).
		Scan(&userDB.ID, &userDB.Email, &userDB.ProfilePicture, &userDB.Name)
	if err != nil {
		return entities.User{}, err
	}
	return toDomainUser(userDB), nil
}

func (u *UserGateway) GetUserByEmail(email string) (entities.User, error) {
	var userDB User
	err := u.db.QueryRow("SELECT id, email, profile_picture, name FROM users WHERE email = ?", email).
		Scan(&userDB.ID, &userDB.Email, &userDB.ProfilePicture, &userDB.Name)
	if err != nil {
		return entities.User{}, err
	}
	return toDomainUser(userDB), nil
}

func toDomainUser(user User) entities.User {
	return entities.User{
		ID:             &user.ID,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Name:           user.Name,
	}
}
