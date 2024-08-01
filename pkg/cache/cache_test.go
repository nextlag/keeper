package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	tests := []struct {
		name               string
		expiration         time.Duration
		sleepDuration      time.Duration
		expectedFoundAfter bool
		expectedValueAfter string
	}{
		{
			name:               "Value found immediately after setting",
			expiration:         1 * time.Millisecond,
			sleepDuration:      0 * time.Millisecond,
			expectedFoundAfter: true,
			expectedValueAfter: "testValue",
		},
		{
			name:               "Value expired after sleep",
			expiration:         1 * time.Millisecond,
			sleepDuration:      2 * time.Millisecond,
			expectedFoundAfter: false,
			expectedValueAfter: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := New(tt.expiration, tt.expiration)
			key := "testKey"
			value := "testValue"

			_, found := cache.Get(key)
			assert.False(t, found, "Key should not be found in cache initially")

			cache.Set(key, value)
			time.Sleep(tt.sleepDuration)

			cachedValue, found := cache.Get(key)
			assert.Equal(t, tt.expectedFoundAfter, found, "Unexpected found status after sleep duration")
			if found {
				assert.Equal(t, tt.expectedValueAfter, cachedValue, "Cached value should match the expected value")
			}
		})
	}
}
