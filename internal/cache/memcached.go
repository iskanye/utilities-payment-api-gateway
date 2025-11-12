package cache

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
)

type Cache struct {
	cl  *memcache.Client
	ttl int32
}

func New(host string, port int, ttl int32) *Cache {
	cl := memcache.New(net.JoinHostPort(host, strconv.Itoa(port)))

	return &Cache{
		cl:  cl,
		ttl: ttl,
	}
}

func (c *Cache) Set(key string, val jwt.TokenPayload) error {
	const op = "memcached.Set"

	value, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: c.ttl,
	}

	err = c.cl.Add(item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) Get(key string) (jwt.TokenPayload, error) {
	const op = "memcached.Get"

	item, err := c.cl.Get(key)
	if err == memcache.ErrCacheMiss {
		return jwt.TokenPayload{}, ErrCacheMiss
	}
	if err != nil {
		return jwt.TokenPayload{}, fmt.Errorf("%s: %w", op, err)
	}

	var payload jwt.TokenPayload
	err = json.Unmarshal(item.Value, &payload)
	if err != nil {
		return jwt.TokenPayload{}, fmt.Errorf("%s: %w", op, err)
	}

	return payload, nil
}
