package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisKvs struct {
	client *redis.Client
}

func NewRedisKvs() (*RedisKvs) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	// if _, err := client.Ping().Result(); err != nil {
	// 	return nil, err
	// }

	return &RedisKvs{
		client: client,
	}
}

func (rk *RedisKvs) Get(key string) (string, bool) {
	val, err := rk.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func (rk *RedisKvs) Set(key, value string) {
	rk.client.Set(ctx, key, value, 0)
}

func (rk *RedisKvs) Delete(key string) {
	rk.client.Del(ctx, key)
}
