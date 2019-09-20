package main

import (
	"fmt"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type redisUniversalClient struct {
	universalCl redis.UniversalClient
	cacheClient *cache.Codec
}

func getRedisClient(addr []string) redis.UniversalClient {
	var universalCl redis.UniversalClient
	var retries = 1
	for {
		universalOpts := redis.UniversalOptions{
			Addrs:           addr,
			MaxRetryBackoff: time.Second * 2,
		}
		universalCl = redis.NewUniversalClient(&universalOpts)
		_, err := universalCl.Ping().Result()
		if err != nil {
			time.Sleep(time.Duration(retries) * time.Second)
		} else {
			return universalCl
		}
	}
}

// NewRedisClient returns a new redis ring client
func NewRedisUniversalClient(addrs []string) StorageDriver {
	universalCl := getRedisClient(addrs)
	codec := &cache.Codec{
		Redis: universalCl,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return &redisUniversalClient{universalCl, codec}
}

// Set caches the item.
func (r *redisUniversalClient) Set(key string, object interface{}, ttl time.Duration) error {
	fmt.Printf("Set: (%s: %v)\n", key, object)
	item := &cache.Item{
		Key:        key,
		Object:     object,
		Expiration: ttl,
	}
	return r.cacheClient.Set(item)
}

// Get gets the object for the given key.
func (r *redisUniversalClient) Get(key string) (interface{}, error) {
	fmt.Printf("Get: (%s)\n", key)
	var object interface{}
	err := r.cacheClient.Get(key, &object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Delete deletes the object and key for the given key from redis ring.
func (r *redisUniversalClient) Delete(key string) error {
	fmt.Printf("Delete: (%s)\n", key)
	return r.cacheClient.Delete(key)
}

// Exists reports whether object for the given key exists.
func (r *redisUniversalClient) Exists(key string) bool {
	fmt.Printf("Exists: (%s)\n", key)
	return r.cacheClient.Exists(key)
}

// Close closes the ring client, releasing any open resources.
func (r *redisUniversalClient) Close() error {
	fmt.Printf("Closing redis universal client...\n")
	return r.universalCl.Close()
}

func redisUniversalClientOps() {
	addrs := []string{":7001", ":7002", ":7003", ":7004", ":7005", ":7006"}
	fmt.Println("redis universal client...")

	universalCl := NewRedisUniversalClient(addrs)
	// close...
	defer universalCl.Close()

	key := "exampleKey"
	value := "exampleValue"
	ttl := 3 * time.Second

	// set...
	err := universalCl.Set(key, value, ttl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// get...
	res, err := universalCl.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}

	// exists...
	res = universalCl.Exists(key)
	fmt.Printf("%v\n", res)

	// unknown key...
	res = universalCl.Exists("randomNonExistentKey")
	fmt.Printf("%v\n", res)

	// expire
	fmt.Printf("sleeping for expire...\n")
	time.Sleep(5 * time.Second)

	// get...
	res, err = universalCl.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}

	// exists...
	res = universalCl.Exists(key)
	fmt.Printf("%v\n", res)

	// delete...
	// set...
	err = universalCl.Set(key, value, ttl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	res = universalCl.Delete(key)
	fmt.Printf("%v\n", res)

	// get...
	res, err = universalCl.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}
}
