package redis

import (
	"context"
	"strconv"

	"github.com/adamkoro/adventcalendar-backend/lib/env"
	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     env.GetRedisHost() + ":" + strconv.Itoa(env.GetRedisPort()),
		Password: env.GetRedisPassword(),
		DB:       env.GetRedisDb(),
	})
}

func Close(client *redis.Client) {
	client.Close()
}

func Ping(client *redis.Client) error {
	_, err := client.Ping(context.Background()).Result()
	return err
}

func Set(client *redis.Client, key string, value interface{}) error {
	return client.Set(context.Background(), key, value, 0).Err()
}

func Get(client *redis.Client, key string) (string, error) {
	return client.Get(context.Background(), key).Result()
}

func Del(client *redis.Client, key string) error {
	return client.Del(context.Background(), key).Err()
}
