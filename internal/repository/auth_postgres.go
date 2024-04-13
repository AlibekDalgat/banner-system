package repository

import (
	"banner-system/internal/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (p *AuthPostgres) GetUser(login, password string) (models.User, bool, error) {
	var user models.User
	query := fmt.Sprintf("SELECT login FROM %s WHERE login=$1 AND password=$2", adminsTable)
	err := p.db.Get(&user, query, login, password)

	isAdmin := false
	if err != nil {
		if err == sql.ErrNoRows {
			query = fmt.Sprintf("SELECT login FROM %s WHERE login=$1 AND password=$2", usersTable)
			err = p.db.Get(&user, query, login, password)
			if err != nil {
				return models.User{}, false, err
			}
		} else {
			return models.User{}, false, err
		}
	} else {
		isAdmin = true
	}
	return user, isAdmin, nil
}
