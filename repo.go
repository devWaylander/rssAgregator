package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewRepo(db *sqlx.DB) *repo {
	return &repo{db: db}
}

func (r *repo) NewUser(ctx context.Context, name string) (error, bool) {
	uuid := uuid.New()
	userModel := User{}
	query := `SELECT user FROM users WHERE name = $1`
	existStatus := r.db.GetContext(ctx, &userModel, query, name)
	if userModel.ID != "" {
		return existStatus, true
	}

	query = `INSERT INTO users (id, name) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, uuid, name)
	if err != nil {
		logrus.Error("Cannot insert new user: ", err)
	}

	return err, false
}
