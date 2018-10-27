package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	errNotFound = errors.New("not found")
)

//go:generate mockery -name=LazyCache

// LazyCache defines the methods for the cache
type LazyCache interface {
	Get(key string) (data *profileapi.Profile, err error)
}

type cacheValueType struct {
	data        *profileapi.Profile
	lastUpdated time.Time
}

type lazyCacheImpl struct {
	sync.RWMutex
	memStore   map[string]cacheValueType
	profileAPI profileapi.UserProfileAPI
	timeout    time.Duration // cache would expire after timeout
}

// NewLazyCache instantiates a default lazy cache implementation
func NewLazyCache(profileClient profileapi.UserProfileAPI, timeout time.Duration) LazyCache {
	return &lazyCacheImpl{
		memStore:   make(map[string]cacheValueType),
		profileAPI: profileAPI,
		timeout:    timeout,
	}
}

// Get implements LazyCache interface
func (c *lazyCacheImpl) Get(key string) (data interface{}, err error) {
	c.RLock()
	val := c.memStore[key]
	c.RUnlock()

	timeNow := time.Now()
	if timeNow.After(val.lastUpdated.Add(c.timeout)) {
		// fetch data from redis and update cache
		latest, err := c.profileClient.Get(key)
		if err != nil && err != profileapi.ErrNotFound {
			return nil, err
		}

		val = cacheValueType{latest, timeNow}

		c.Lock()
		c.memStore[key] = val
		c.Unlock()

		return val.data, nil
	}

	return val.data, err
}
