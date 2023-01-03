package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yapi-teklif/internal/pkg/database/config"
)

type ICacheDB interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) (int64, error)
	Client() *redis.Client
	TTl(key string) *redis.DurationCmd
}

var ctx = context.Background()

type CacheDB struct {
	Redis *redis.Client
}

func NewCacheClient() *CacheDB {
	databaseConfig := config.NewDatabase()
	redisDB := redis.NewClient(&redis.Options{
		Addr:     databaseConfig.Redis.DBAddress,
		Password: databaseConfig.Redis.DBPassword,
		DB:       databaseConfig.Redis.DB,
	})
	pong, err := redisDB.Ping(ctx).Result()

	if err != nil {

		log.Fatal(err)

	}

	// return pong if server is online

	fmt.Println(pong)
	return &CacheDB{
		Redis: redisDB,
	}
}

const prefix = "GO"

func (s *CacheDB) Set(key string, value interface{}, expiration time.Duration) error {
	key = fmt.Sprintf("%s_%s", prefix, key)
	return s.Redis.Set(ctx, key, value, expiration).Err()
}

func (s *CacheDB) Get(key string) (string, error) {
	key = fmt.Sprintf("%s_%s", prefix, key)
	return s.Redis.Get(ctx, key).Result()
}

func (s *CacheDB) Del(key string) (int64, error) {
	key = fmt.Sprintf("%s_%s", prefix, key)
	return s.Redis.Del(ctx, key).Result()
}

func (s *CacheDB) TTl(key string) *redis.DurationCmd {
	key = fmt.Sprintf("%s_%s", prefix, key)
	return s.Redis.TTL(ctx, key)
}

func (s *CacheDB) Client() *redis.Client {
	return s.Redis
}
