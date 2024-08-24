package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"image_processing/errors"
	"strconv"
	"time"
)

type LinkStore interface {
	InsertLink(context.Context, int64, string) *errors.Error
	GetLinkID(context.Context, string) (int, *errors.Error)
}

type RedisLinkStore struct {
	conn *redis.Client
}

func NewRedisLinkStore(client *redis.Client) *RedisLinkStore {
	return &RedisLinkStore{conn: client}
}

func (s *RedisLinkStore) InsertLink(ctx context.Context, id int64, link string) *errors.Error {
	err := s.conn.Set(ctx, link, id, time.Hour*24).Err()
	if err != nil {
		return errors.ErrDB(err.Error())
	}
	return nil
}
func (s *RedisLinkStore) GetLinkID(ctx context.Context, link string) (int, *errors.Error) {
	val, err := s.conn.Get(ctx, link).Result()
	if err != nil {
		errors.ErrDB(err.Error())
	}
	id, _ := strconv.Atoi(val)
	return id, nil
}
