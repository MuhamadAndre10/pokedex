package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPokacache(t *testing.T) {

	interval := 10 * time.Second

	testCache := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range testCache {
		t.Run(fmt.Sprintf("Test Case %v", i), func(t *testing.T) {
			cacheEntry := NewCache(interval)
			cacheEntry.Set(c.key, c.val)

			data, found := cacheEntry.Get(c.key)

			if !found {
				t.Errorf("cache with key %s not found", string(data))
			}

			assert.Equal(t, c.val, data, "harus sama")
		})

	}

}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	cacheEntry := NewCache(baseTime)

	cacheEntry.Set("https://example.com", []byte("testdata"))

	_, ok := cacheEntry.Get("https://example.com")
	assert.Equal(t, ok, true, "data ok")

	time.Sleep(waitTime)

	_, ok = cacheEntry.Get("https://example.com")
	assert.NotEqual(t, ok, true, "data sudah tidak ada")

}
