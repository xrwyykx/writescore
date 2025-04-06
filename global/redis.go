package global

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisConn *redis.Client

func GetRedisConn() *redis.Client {
	return redisConn
}
func InitRedis() {
	addr := fmt.Sprintf("%s:%d", viper.Get("redis.host"), viper.GetInt("redis.port"))
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	c := context.Background()
	pong, err := redisConn.Ping(c).Result()
	if err != nil {
		fmt.Printf("连接redis失败：%v\n", err)
		return
	}
	fmt.Printf("连接redis成功:%s\n", pong)
}
