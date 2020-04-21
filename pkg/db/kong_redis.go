package db

import (
	"errors"
	"gopkg.in/redis.v5"
	"sync"
)

var Redis sync.Map

func GetRedis(DBAddr, Auth string, DBIndex int) (rds *redis.Client, err error) {

	if DBAddr == "" {
		return nil, errors.New("DBConnect Is Nil")
	}

	v, ok := Redis.Load(DBAddr)
	if ok {
		rds = v.(*redis.Client)
	} else {
		options := redis.Options{Addr: DBAddr, Password: Auth, DB: DBIndex}
		client := redis.NewClient(&options)
		_, err := client.Ping().Result()
		if err != nil {
			return nil, err
		}

		Redis.LoadOrStore(DBAddr, client)
	}

	return rds, err
}
