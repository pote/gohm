package gohm

import(
	"github.com/garyburd/redigo/redis"
	"github.com/pote/redisurl"
)

type Connection struct {
	RedisPool *redis.Pool
	luaSave   *redis.Script
}

func NewConnection(r... *redis.Pool) (*Connection, error) {
	if len(r) < 1 {
		pool, err := redisurl.NewPool(3, 200, "240s")
		if err != nil {
			return &Connection{}, err
		}	

  return NewConnectionWithPool(pool), nil
	} else {
    return NewConnectionWithPool(r[0]), nil
	}

}

func NewConnectionWithPool(pool *redis.Pool) *Connection {
	c := &Connection{
		RedisPool: pool,
	}

	c.luaSave = redis.NewScript(0, LUA_SAVE)

	return c
}
