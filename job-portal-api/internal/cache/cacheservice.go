package cache

import (
	"context"
	"encoding/json"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
	"time"
)

func (c *Cache) CheckRedisKey(key string) (models.Job, error) {
	var ctx = context.Background()
	val, err := c.Rd.Get(ctx, key).Result()
	if err == redis2.Nil {
		return models.Job{}, err

	}
	var job models.Job
	err = json.Unmarshal([]byte(val), &job)
	if err != nil {
		log.Err(err)
	}
	return job, nil
}
func (c *Cache) SetRedisKey(key string, value models.Job) {
	var ctx = context.Background()
	jobdata, err := json.Marshal(value)
	if err != nil {
		log.Err(err)
		return
	}
	data := string(jobdata)
	err = c.Rd.Set(ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		log.Err(err)
		return
	}

}