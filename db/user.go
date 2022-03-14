package db

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/jasanchez1/Dpricing/models"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}

	rows, err := db.Conn.Query("SELECT * FROM users ORDER BY ID DESC")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Name, &user.Description, &user.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (db Database) AddUser(user *models.User) error {
	var id uuid.UUID
	var createdAt string
	query := `INSERT INTO users (name, description) VALUES ($1, $2) RETURNING user_id, created_at`
	err := db.Conn.QueryRow(query, user.Name, user.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}

	user.UserID = id
	user.CreatedAt = createdAt
	return nil
}

func (db Database) GetUserById(userId uuid.UUID) (models.User, error) {
	user := models.User{}

	query := `SELECT * FROM users WHERE id = $1;`
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.UserID, &user.Name, &user.Description, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) DeleteUser(userId uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateUser(userId uuid.UUID, userData models.User) (models.User, error) {
	user := models.User{}
	query := `UPDATE users SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Conn.QueryRow(query, userData.Name, userData.Description, userId).Scan(&user.UserID, &user.Name, &user.Description, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}

	return user, nil
}
