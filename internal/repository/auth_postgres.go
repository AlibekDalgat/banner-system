package repository

import (
	"banner-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func (p *AuthPostgres) GetUser(login, password string) (models.User, error) {
	return models.User{}, nil
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}
