package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"

	"sweepstake/conf"
)

// GetConn get the connection for redis
func GetConn() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.RedisConf.Host, conf.RedisConf.Port))
	if err != nil {
		log.Println("connect redis error", err)
		return nil, err
	}

	return conn, nil
}

//InitRedis used to init award_time and award_remain_num
func InitRedis() error {
	conn, err := GetConn()
	if err != nil {
		log.Println("redis conn is nil")
		return err
	}
	defer conn.Close()

	startTime, _ := time.Parse(conf.InitTimeConf.Layout, conf.InitTimeConf.StartTime)
	conn.Send("ZADD", "award_remain_num", conf.AwardConf.A, "A")
	conn.Send("ZADD", "award_remain_num", conf.AwardConf.B, "B")
	conn.Send("ZADD", "award_remain_num", conf.AwardConf.C, "C")
	conn.Send("HSET", "award_time", "A", startTime.Format(time.RFC3339))
	conn.Send("HSET", "award_time", "B", startTime.Format(time.RFC3339))
	conn.Send("HSET", "award_time", "C", startTime.Format(time.RFC3339))
	conn.Flush()

	for i := 0; i < 6; i++ {
		_, err := conn.Receive()
		if err != nil {
			log.Printf("conn send error, %s", err)
			return err
		}
	}

	return nil
}
