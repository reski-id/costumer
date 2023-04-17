package utils

import "github.com/go-redis/redis"

// RedisConfig represents the Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ConnectRedis initializes and returns a Redis client
func ConnectRedis(config RedisConfig) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
