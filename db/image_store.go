package db

import (
	"context"
	"database/sql"
	"image_processing/errors"
	"image_processing/types"
)

type ImageStore interface {
	InsertImage(context.Context, *types.Image) (*types.Image, *errors.Error)
	GetImageByID(context.Context, int) (*types.Image, *errors.Error)
}

type PostgresImageStore struct {
	conn *sql.DB
}

func NewPostgresImageStore(conn *sql.DB) *PostgresImageStore {
	return &PostgresImageStore{conn: conn}
}
func (s *PostgresImageStore) InsertImage(ctx context.Context, img *types.Image) (*types.Image, *errors.Error) {
	res, err := s.conn.QueryContext(ctx, "INSERT INTO images (name,user_id,format) VALUES ($1, $2, $3) RETURNING id", img.Name, img.UserID, img.Format)

	if err != nil {
		return nil, errors.ErrDB(err.Error())
	}
	defer res.Close()
	if !res.Next() {
		return nil,
			errors.ErrDB("no rows returned from insert " + err.Error())
	}
	err = res.Scan(&img.ID)
	if err != nil {
		return nil, errors.ErrServer(err.Error())
	}

	return img, nil
}
func (s *PostgresImageStore) GetImageByID(ctx context.Context, id int) (*types.Image, *errors.Error) {
	image := new(types.Image)
	res, err := s.conn.QueryContext(ctx, "SELECT * FROM images WHERE id=$1", id)
	if !res.Next() {
		return nil, errors.ErrResourceNotFound(string(id))
	}
	if err != nil {
		return nil, errors.ErrDB(err.Error())
	}
	defer res.Close()
	err = res.Scan(&image.ID, &image.UserID, &image.Name, &image.Format)
	if err != nil {
		return nil, errors.ErrServer(err.Error())
	}
	return image, nil
}
