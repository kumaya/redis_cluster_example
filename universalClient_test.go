package main

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

// newTestRedisUniversalClient returns a NewRedisUniversalClient
func newTestRedisUniversalClient() StorageDriver {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := NewRedisUniversalClient([]string{mr.Addr()})

	return client
}

func TestRedisUniversalClient_Set(t *testing.T) {
	key := "testKey"
	value := "testValue"
	ttl := 1 * time.Second
	uc := newTestRedisUniversalClient()

	err := uc.Set(key, value, ttl)
	if err != nil {
		t.Errorf("%v", err)
	}
	assert.True(t, uc.Exists(key))
}

func TestRedisUniversalClient_Get(t *testing.T) {
	key := "testKey"
	value := "testValue"
	ttl := 1 * time.Second
	uc := newTestRedisUniversalClient()

	err := uc.Set(key, value, ttl)
	if err != nil {
		t.Errorf("%v", err)
	}

	res, err := uc.Get(key)
	if err != nil {
		t.Errorf("%v", err)
	}
	assert.Equal(t, res, value, "Expected: %s; Got: %s", value, res)
}

func TestRedisUniversalClient_Exists(t *testing.T) {
	key := "testKey"
	value := "testValue"
	ttl := 1 * time.Second
	uc := newTestRedisUniversalClient()

	assert.False(t, uc.Exists(key))

	uc.Set(key, value, ttl)

	assert.True(t, uc.Exists(key))
}

func TestRedisUniversalClient_Delete(t *testing.T) {
	key := "testKey"
	value := "testValue"
	ttl := 1 * time.Second
	uc := newTestRedisUniversalClient()

	assert.False(t, uc.Exists(key))

	uc.Set(key, value, ttl)

	assert.True(t, uc.Exists(key))

	err := uc.Delete(key)
	if err != nil {
		t.Errorf("%v", err)
	}

	assert.False(t, uc.Exists(key))
}
