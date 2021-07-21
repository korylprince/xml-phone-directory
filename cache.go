package main

import (
	"sync"
	"time"
)

type Cache struct {
	entries  []*Entry
	mu       *sync.Mutex
	interval time.Duration
	expires  time.Time
}

func NewCache(cacheTime time.Duration) *Cache {
	return &Cache{mu: new(sync.Mutex), interval: cacheTime}
}

func (c *Cache) GetEntries(dsn string) ([]*Entry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.entries == nil || time.Now().After(c.expires) {
		entries, err := GetEntries(dsn)
		if err != nil {
			return nil, err
		}

		c.entries = entries
		c.expires = time.Now().Add(c.interval)
		return entries, nil
	}

	return c.entries, nil
}
