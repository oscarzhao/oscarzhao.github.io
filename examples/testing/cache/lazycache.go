package cache

import (
	"sync"
	"time"

	"github.com/oscarzhao/oscarzhao.github.io/examples/testing/thirdpartyapi"
)

//go:generate mockery -name=LazyCache

// LazyCache defines the methods for the cache
type LazyCache interface {
	Get(key string) (data interface{}, err error)
}

// NewLazyCache instantiates a default lazy cache implementation
func NewLazyCache(client thirdpartyapi.Client, timeout time.Duration) LazyCache {
	return &lazyCacheImpl{
		cacheStore:       make(map[string]cacheValueType),
		thirdPartyClient: client,
		timeout:          timeout,
	}
}

type cacheValueType struct {
	data        interface{}
	lastUpdated time.Time
}

type lazyCacheImpl struct {
	sync.RWMutex
	cacheStore       map[string]cacheValueType
	thirdPartyClient thirdpartyapi.Client
	timeout          time.Duration // cache would expire after timeout
}

// Get implements LazyCache interface
func (c *lazyCacheImpl) Get(key string) (data interface{}, err error) {
	c.RLock()
	val := c.cacheStore[key]
	c.RUnlock()

	timeNow := time.Now()
	if timeNow.After(val.lastUpdated.Add(c.timeout)) {
		// fetch data from third party service and update cache
		latest, err := c.thirdPartyClient.Get(key)
		if err != nil {
			return nil, err
		}

		val = cacheValueType{latest, timeNow}

		c.Lock()
		c.cacheStore[key] = val
		c.Unlock()
	}

	return val.data, nil
}
