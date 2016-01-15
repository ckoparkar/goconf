package main

import (
	"log"

	redis "gopkg.in/redis.v3"
)

type Redis struct {
	*redis.Client
}

func NewRedisClient(opt *redis.Options) (*Redis, error) {
	rconn := redis.NewClient(opt)
	_, err := rconn.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Redis{rconn}, nil
}

func (r *Redis) GetAllKV() <-chan KV {
	keys := r.GetAllKeys()
	out := make(chan KV)
	go func() {
		for k := range keys {
			if v, err := r.Get(k).Result(); err == nil {
				// kvs = append(kvs, KV{Key: k, Value: v})
				out <- KV{Key: k, Value: v}
			}
		}
		close(out)
	}()
	return out
}

func (r *Redis) GetAllKeys() <-chan string {
	out := make(chan string)
	go func() {
		var cursor int64
		for {
			var err error
			cursor, keys, err := r.Scan(cursor, "", 100000).Result()
			if err != nil {
				log.Println(err)
			}
			for _, k := range keys {
				out <- k
			}
			if cursor == 0 {
				close(out)
				break
			}
		}
	}()
	return out
}
