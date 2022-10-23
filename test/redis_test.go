package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedisSet(t *testing.T) {
	err := rdb.Set("key2", "zyj", time.Second*60).Err()
	if err != nil {
		panic(err)
	}

}
func TestRedisGet(t *testing.T) {
	val, err := rdb.Get("key2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
