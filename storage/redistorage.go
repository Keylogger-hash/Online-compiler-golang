package storage

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)



func NewRedisClient() *redis.Client{
	opt := &redis.Options{
		Network: "tcp",
		Addr: 	 "localhost:6379",
		DB: 0,
		Password: "",
	}
	client := redis.NewClient(opt)
	return client
}



func AddRedisValue(client *redis.Client,key string, value []byte) error{
	ctx := context.Background()
	err := client.Set(ctx,key,value,0).Err()
	return err
}
func GetRedisValue(client *redis.Client,key string) (string,error){
	ctx := context.Background()
	value, err := client.Get(ctx,key).Result()
	return value,err
}