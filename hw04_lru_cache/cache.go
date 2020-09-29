package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	cap   int
	queue List
	cache map[Key]*listItem
	sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, ok := c.cache[key]; !ok {
		return nil, false
	} else {
		c.queue.MoveToFront(c.cache[key])
		c.cache[key] = c.queue.Front()
		return c.cache[key].Value.(*cacheItem).value, true
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if _, ok := c.cache[key]; !ok {
		if c.queue.Len() >= c.cap {
			c.Clear()
		}
		newCacheItem := &cacheItem{value: value, key: key}
		c.queue.PushFront(newCacheItem)
		c.cache[key] = c.queue.Front()
		return false
	} else {
		c.cache[key].Value.(*cacheItem).value = value
		c.queue.MoveToFront(c.cache[key])
		c.cache[key] = c.queue.Front()
		return true
	}
}

func (c *lruCache) Clear() {
	lastItem := c.queue.Back().Value.(*cacheItem)
	c.queue.Remove(c.queue.Back())
	delete(c.cache, lastItem.key)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		cap:   capacity,
		queue: NewList(),
		cache: make(map[Key]*listItem),
	}
}
