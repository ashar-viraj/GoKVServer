package cache

import (
	"container/list"
	"os"
	"strconv"
	"sync"
)

const MAX_CACHE_SIZE = 5

type entry struct {
	key   int
	value string
}

type LRUCache struct {
	mu       sync.Mutex
	cache    map[int]*list.Element
	eviction *list.List
}

func NewCache() *LRUCache {
	return &LRUCache{
		cache:    make(map[int]*list.Element),
		eviction: list.New(),
	}
}

func (c *LRUCache) Get(key int) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.eviction.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return "", false
}

func (c *LRUCache) Put(key int, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key exists, update it and move to front
	if elem, found := c.cache[key]; found {
		elem.Value.(*entry).value = value
		c.eviction.MoveToFront(elem)
		return
	}

	// If full, evict LRU
	MAX_CACHE_SIZE, _ := strconv.Atoi(os.Getenv("MAX_CACHE_SIZE"))
	if c.eviction.Len() >= MAX_CACHE_SIZE {
		oldest := c.eviction.Back()
		if oldest != nil {
			c.eviction.Remove(oldest)
			delete(c.cache, oldest.Value.(*entry).key)
		}
	}

	elem := c.eviction.PushFront(&entry{key, value})
	c.cache[key] = elem
}

func (c *LRUCache) Delete(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cache[key]; found {
		c.eviction.Remove(elem)
		delete(c.cache, key)
	}
}
