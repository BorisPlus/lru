package lru

import (
	"fmt"
	"sync"
)

type Key string

// KeyValue - в хранилище будет учтена пара.
//
// Пара пригодится при извлесении элемента из списка и
// необходимостью поиска в карте, в частности, при
// очистке абсолютно заполненного кэша.
type KeyValue struct {
	key   Key
	value interface{}
}

func NewKeyValuePair(key Key, value interface{}) *KeyValue {
	return &KeyValue{key, value}
}

// Key - получить значение KeyValue.key.
func (kv *KeyValue) Key() Key {
	return kv.key
}

// Value - получить значение KeyValue.value.
func (kv *KeyValue) Value() interface{} {
	return kv.value
}

// String - наглядное представление значения KeyValue-структуры.
func (kv KeyValue) String() string {
	return fmt.Sprintf("%q->%q", kv.key, kv.value)
}

// LruCache - структура кэша.
type LruCache struct {
	capacity int
	list     Lister
	items    map[Key]*ListItem // Примерно как key-value database
	mutex    *sync.RWMutex
}

// Set - уставновка элемента в кэш.
func (cache *LruCache) Set(key Key, value interface{}) bool {
	if cache.capacity == 0 {
		return false
	}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item, exists := cache.items[key]
	if exists {
		item.Data = KeyValue{key, value}
		cache.list.MoveToFront(item)
		return true
	}
	if cache.capacity == cache.list.Len() {
		back := cache.list.Back()
		delete(cache.items, back.Data.(KeyValue).key)
		cache.list.Remove(back)
	}
	cache.items[key] = cache.list.PushFront(KeyValue{key, value})
	return false
}

// Get - получение элемента из кэша.
func (cache *LruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	item, exists := cache.items[key]
	if exists {
		cache.list.MoveToFront(item)
		return item.Data.(KeyValue).value, true
	}
	return nil, false
}

// Clear - "очистка" кэша.
func (cache *LruCache) Clear() {
	*cache = *(NewCache(cache.capacity).(*LruCache))
}

// Cacher - интерфейс хранения кэша.
type Cacher interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// NewCache - функция-конструктор кэша.
func NewCache(capacity int) Cacher {
	return &LruCache{
		capacity: capacity,
		list:     NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    &sync.RWMutex{},
	}
}
