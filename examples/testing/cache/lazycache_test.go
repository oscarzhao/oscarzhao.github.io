package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/oscarzhao/oscarzhao.github.io/examples/testing/thirdpartyapi/mocks"
)

const (
	testKeyInCache = "key-in-cache"
	testValInCache = "value-in-cache"

	testTimeout = time.Second * 10
)

func TestGet_CacheHit(t *testing.T) {
	mockThirdParty := &mocks.Client{}

	mockCache := &lazyCacheImpl{
		memStore: map[string]cacheValueType{
			testKeyInCache: cacheValueType{testValInCache, time.Now()},
		},
		thirdPartyClient: mockThirdParty,
		timeout:          testTimeout,
	}

	// test cache hit, with valid value
	got, gotErr := mockCache.Get(testKeyInCache)
	require.Equal(t, nil, gotErr)
	require.Equal(t, testValInCache, got)

	mock.AssertExpectationsForObjects(t, mockThirdParty)
}

func TestGet_CacheHit_Expired_Update_Success(t *testing.T) {
	mockThirdParty := &mocks.Client{}
	mockThirdParty.On("Get", testKeyInCache).Return(testValInCache, nil).Once()

	timeTooOld := time.Now().Add(-testTimeout - time.Second)

	mockCache := &lazyCacheImpl{
		memStore: map[string]cacheValueType{
			testKeyInCache: cacheValueType{testValInCache, timeTooOld},
		},
		thirdPartyClient: mockThirdParty,
		timeout:          testTimeout,
	}

	// test cache miss, with valid value
	got, gotErr := mockCache.Get(testKeyInCache)
	require.Equal(t, nil, gotErr)
	require.Equal(t, testValInCache, got)

	mock.AssertExpectationsForObjects(t, mockThirdParty)
}

func TestGet_CacheMiss_Update_Success(t *testing.T) {
	mockThirdParty := &mocks.Client{}
	mockThirdParty.On("Get", testKeyInCache).Return(testValInCache, nil).Once()

	mockCache := &lazyCacheImpl{
		memStore:         map[string]cacheValueType{},
		thirdPartyClient: mockThirdParty,
		timeout:          testTimeout,
	}

	// test cache miss, with valid value
	got, gotErr := mockCache.Get(testKeyInCache)
	require.Equal(t, nil, gotErr)
	require.Equal(t, testValInCache, got)

	mock.AssertExpectationsForObjects(t, mockThirdParty)
}

func TestGet_CacheMiss_Update_Failure(t *testing.T) {
	errTest := errors.New("test error")
	mockThirdParty := &mocks.Client{}
	mockThirdParty.On("Get", testKeyInCache).Return(nil, errTest).Once()

	mockCache := &lazyCacheImpl{
		memStore:         map[string]cacheValueType{},
		thirdPartyClient: mockThirdParty,
		timeout:          testTimeout,
	}

	// test cache miss, fails to fetch from data source
	_, gotErr := mockCache.Get(testKeyInCache)
	require.Equal(t, errTest, gotErr)

	mock.AssertExpectationsForObjects(t, mockThirdParty)
}
