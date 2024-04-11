package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

type Reader interface {
	GetPool() *redis.Pool
	Ping() error
	Set(key string, data interface{}, expire int) error
	Exists(key string) bool
	Get(key string) (string, error)
	GetFormatted(key string, resp interface{}) error
	Delete(key string) (bool, error)
	LikeDeletes(strKey string) error
}

type Base struct {
	Pool *redis.Pool
}

func New(host string, port string, pwd string, idleTimeout string, maxIdle int, maxActive int) Reader {
	idleTime, _ := time.ParseDuration(idleTimeout)

	return &Base{
		Pool: &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTime,
			Dial: func() (redis.Conn, error) {
				dial, err := redis.Dial("tcp", host+":"+port)
				if err != nil {
					return nil, err
				}
				if pwd != "" {
					if _, cErr := dial.Do("AUTH", pwd); cErr != nil {
						if ccErr := dial.Close(); ccErr != nil {
							return nil, ccErr
						}
						return nil, cErr
					}
				}
				return dial, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (r *Base) GetPool() *redis.Pool {
	return r.Pool
}
func (r *Base) Ping() error {
	pool := r.Pool.Get()
	_, err := pool.Do("PING")
	return err
}

func (r *Base) Set(key string, data interface{}, expire int) error {
	pool := r.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("error closing the redis connection", "err", err)
		}
	}(pool)

	str, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err = pool.Do("SET", key, str); err != nil {
		return err
	}

	if expire > 0 {
		if _, err = pool.Do("EXPIRE", key, expire); err != nil {
			return err
		}
	}
	return nil

}

func (r *Base) Exists(key string) bool {
	pool := r.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("error closing the redis connection", "err", err)
		}
	}(pool)

	exists, err := redis.Bool(pool.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func (r *Base) Get(key string) (string, error) {
	pool := r.Pool.Get()
	defer func(conn redis.Conn) {
		if err := conn.Close(); err != nil {
			fmt.Printf("error closing the redis connection", "err", err)
		}
	}(pool)
	data, err := redis.String(pool.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return data, errors.Wrap(err, fmt.Sprintf("error retriving key %s", key))
	}
	return data, err
}

func (r *Base) GetFormatted(key string, resp interface{}) error {
	dataStr, err := r.Get(key)
	if err != nil {
		return err
	}

	if err1 := json.Unmarshal([]byte(dataStr), &resp); err1 != nil {
		return err1
	}
	return nil
}

func (r *Base) Delete(key string) (bool, error) {
	pool := r.Pool.Get()
	defer func(pool redis.Conn) {
		if err := pool.Close(); err != nil {
			fmt.Printf("error closing the redis connection", "err", err)
		}
	}(pool)

	return redis.Bool(pool.Do("DEL", key))
}

func (r *Base) LikeDeletes(strKey string) error {
	pool := r.Pool.Get()
	defer func(pool redis.Conn) {
		if err := pool.Close(); err != nil {
			fmt.Printf("error closing the redis connection", "err", err)
		}
	}(pool)

	keys, err := redis.Strings(pool.Do("KEYS", "*"+strKey+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, err = r.Delete(key); err != nil {
			return err
		}
	}
	return nil
}
