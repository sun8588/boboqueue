package Model

import (
	"Config"
	"Lib/redis/redis"
	"fmt"
)
var redisChan chan *redis.Client
func init(){
poolSize:=Config.GetRedisConfig().PoolSize
redisChan=make(chan *redis.Client ,poolSize)
	for i:=0;i<poolSize;i++{
		redisChan<-redis.New("tcp:"+Config.GetRedisConfig().Addr+":"+Config.GetRedisConfig().Port, 0, "")
	}
}

func GetRedisConn()*redis.Client{
	conn:= <-redisChan
	fmt.Printf("now chan=%d\n",len(redisChan))
	return conn
	//	redisConn.Set("tes",123)
	//	retn,err:=redisconn.Get("tes")

}
func RestoreRedisConn(redisConn *redis.Client){
	redisChan<-redisConn
	fmt.Printf("count chan=%d\n",len(redisChan))
}