package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisConnection() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: "gsr_redis:6379",
	})

	return &Redis{
		Client: client,
	}
}

func (r *Redis) Get(course string) (string, error) {
	c := r.Client

	courseKey := fmt.Sprintf("course_%s", course)
	s, err := c.Get(courseKey).Result()
	if err != nil {
		return "", err
	}

	return s, nil
}

func (r *Redis) Set(course, shardUrl string) (string, error) {
	c := r.Client

	courseKey := fmt.Sprintf("course_%s", course)

	ttl, _ := time.ParseDuration(
		fmt.Sprintf("%sm", os.Getenv("REDIS_EXPIRATION_TTL")),
	)

	cmdStatus := c.Set(courseKey, shardUrl, ttl)
	s, err := cmdStatus.Result()
	if err != nil {
		return "", err
	}

	return s, nil
}
