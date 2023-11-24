package cache

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"job-portal-api/internal/models"
)

//go:generate mockgen -source cache.go -destination mock_cache.go -package cache

type Cache struct {
	Rd *redis.Client
}

type UserCache interface {
	SetRedisKey(key string, value models.Job)
	CheckRedisKey(key string) (models.Job, error)
	SetRedisKeyOtp(key string, value string) error
	GetRedisKeyOtp(key string) (string, error)
	DelRedisKey(email string) error
}

func NewCache(rd *redis.Client) (UserCache, error) {
	if rd == nil {
		return nil, errors.New("db cannot be null")
	}
	return &Cache{
		Rd: rd,
	}, nil
}
