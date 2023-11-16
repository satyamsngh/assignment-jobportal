package cache

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"job-portal-api/internal/models"
)

type Cache struct {
	Rd *redis.Client
}

type UserCache interface {
	SetRedisKey(key string, value models.Job)
	CheckRedisKey(key string) (models.Job, error)
}

func NewCache(rd *redis.Client) (UserCache, error) {
	if rd == nil {
		return nil, errors.New("db cannot be null")
	}
	return &Cache{
		Rd: rd,
	}, nil
}
