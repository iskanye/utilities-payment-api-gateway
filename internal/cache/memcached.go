package cache

import (
	"fmt"
	"net"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
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

func (c *Cache) Set(key string) error {
	const op = "memcached.Set"

	item := &memcache.Item{
		Key:        key,
		Value:      nil,
		Expiration: c.ttl,
	}

	err := c.cl.Add(item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) Get(key string) error {
	const op = "memcached.Get"

	_, err := c.cl.Get(key)
	if err == memcache.ErrCacheMiss {
		return ErrCacheMiss
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
