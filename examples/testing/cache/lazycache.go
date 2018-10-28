package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/oscarzhao/oscarzhao.github.io/examples/testing/thirdpartyapi"
)

var (
	errNotFound = errors.New("not found")
)

//go:generate mockery -name=LazyCache

// LazyCache defines the methods for the cache
type LazyCache interface {
	Get(key string) (data interface{}, err error)
}

type cacheValueType struct {
	data        interface{}
	lastUpdated time.Time
}

type lazyCacheImpl struct {
	sync.RWMutex
	memStore         map[string]cacheValueType
	thirdPartyClient thirdpartyapi.Client
	timeout          time.Duration // cache would expire after timeout
}

// NewLazyCache instantiates a default lazy cache implementation
func NewLazyCache(client thirdpartyapi.Client, timeout time.Duration) LazyCache {
	return &lazyCacheImpl{
		memStore:         make(map[string]cacheValueType),
		thirdPartyClient: client,
		timeout:          timeout,
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
		latest, err := c.thirdPartyClient.Get(key)
		if err != nil && err != thirdpartyapi.ErrNotFound {
			return nil, err
		}

		val = cacheValueType{latest, timeNow}

		c.Lock()
		c.memStore[key] = val
		c.Unlock()
	}

	if val.data == nil {
		return nil, errNotFound
	}

	return val.data, nil
}
