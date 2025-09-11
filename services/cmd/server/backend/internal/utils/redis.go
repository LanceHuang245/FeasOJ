package utils

import (
	"FeasOJ/internal/config"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// ConnectRedis 连接到Redis并返回redis.Client对象
func ConnectRedis() *redis.Client {
	cfg := config.GlobalConfig.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	})
	return rdb
}

// SetCache 数据缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ConnectRedis().Set(key, jsonData, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string, dest interface{}) error {
	val, err := ConnectRedis().Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return nil // 缓存未命中
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
