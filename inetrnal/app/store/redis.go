package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"project/inetrnal/app/model"
	"time"
)

type RedisRepository struct {
	Redis *redis.Client
}

const (
	SomeDBKey = "somedbkey-%s-%d"
)

func (r *RedisRepository) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Redis.Set(ctx, key, value, expiration).Err()
}

func (r *RedisRepository) GetValue(ctx context.Context, key string) (interface{}, error) {
	var value string
	if err := r.Redis.Get(ctx, key).Scan(&value); err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisRepository) SetSomeDBCache(ctx context.Context, parameter1 string, parameter2 int, result *model.DbModel) error {
	key := fmt.Sprintf(SomeDBKey, parameter1, parameter2)
	r.Redis.Del(ctx, key)
	expirations := time.Hour * time.Duration(24)
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return r.Redis.Set(ctx, key, jsonData, expirations).Err()
}

func (r *RedisRepository) GetSomeDBCache(ctx context.Context, parameter1 string, parameter2 int) (*model.DbModel, error) {
	key := fmt.Sprintf(SomeDBKey, parameter1, parameter2)
	jsonData, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	data := &model.DbModel{}
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
