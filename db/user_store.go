package db

import (
	"context"
	"database/sql"
	"image_processing/errors"

	"image_processing/types"
)

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, *errors.Error)
}

type PostgresUserStore struct {
	conn *sql.DB
}

func NewPostgresUserStore(conn *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{conn: conn}
}
func (s *PostgresUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, *errors.Error) {
	res, err := s.conn.QueryContext(ctx, "SELECT * FROM users WHERE username=$1", user.Username)
	if res.Next() {
		return nil, errors.ErrBadRequest("ss")
	}
	if err != nil {
		return nil, errors.ErrDB(err.Error())
	}
	res, err = s.conn.QueryContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.EncryptedPassword)

	if err != nil {
		return nil, errors.ErrDB(err.Error())
	}
	defer res.Close()
	if !res.Next() {
		return nil, errors.ErrDB("no rows returned from insert")
	}
	err = res.Scan(&user.ID)
	if err != nil {
		return nil, errors.ErrServer(err.Error())
	}
	return user, nil
}
