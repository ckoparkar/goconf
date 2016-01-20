package main

import (
	"io"
	"log"
	"strings"

	redis "gopkg.in/redis.v3"
)

type Redis struct {
	*redis.Client
}

func NewRedisClient(hostport string) (*Redis, error) {
	opt := &redis.Options{Addr: hostport}
	rconn := redis.NewClient(opt)
	_, err := rconn.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Redis{rconn}, nil
}

func (r *Redis) GetAllKV() <-chan KV {
	keys := r.GetAllKeys()
	out := make(chan KV, 10)
	go func() {
		for k := range keys {
			if v, err := r.Get(k).Result(); err == nil {
				out <- KV{Key: k, Value: v}
			}
		}
		close(out)
	}()
	return out
}
func (r *Redis) GetAllACL() <-chan KV {
	out := make(chan KV)
	go func() {
		close(out)
	}()
	return out
}

func (r *Redis) GetAllKeys() <-chan string {
	out := make(chan string, 10)
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

func (r *Redis) GetACL(token string) []string {
	tokenExists, _ := r.HExists("acl", token).Result()
	if token == "" || !tokenExists {
		token = "anonymous"
	}
	aclString, _ := r.HGet("acl", token).Result()
	acls := strings.Split(aclString, ",")
	return acls
}

func (r *Redis) SetKV(kv KV) error {
	return nil
}

func (r *Redis) SetACL(kv KV) error {
	return nil
}

func (r *Redis) Backup(w io.Writer) (int, error) {
	return 0, nil
}
