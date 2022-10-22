package limitip

import (
	"github.com/go-redis/redis"
)

// new redis client
func NewRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return client
}

// test connection
func TestRedisConnection(client *redis.Client) error {
	_, err := client.Ping().Result()
	return err
}

// return ping if server is online
func Ping(client *redis.Client) (string, error) {
	return client.Ping().Result()
}
