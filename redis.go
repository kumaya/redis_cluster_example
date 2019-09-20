package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type redisClient struct {
	clusterCL *redis.ClusterClient
}

func getRedisCluster(addr []string) *redis.ClusterClient {
	var cluster *redis.ClusterClient
	var retries = 1
	for {
		clusterOpts := redis.ClusterOptions{
			Addrs:           addr,
			MaxRetryBackoff: time.Second * 2,
		}
		cluster = redis.NewClusterClient(&clusterOpts)
		_, err := cluster.Ping().Result()
		if err != nil {
			time.Sleep(time.Duration(retries) * time.Second)
		} else {
			return cluster
		}
	}
}

// NewRedisClient returns a new redis ring client
func NewRedisClient(addrs []string) StorageDriver {
	clusterCl := getRedisCluster(addrs)
	return &redisClient{clusterCl}
}

// Set caches the item.
func (r *redisClient) Set(key string, object interface{}, ttl time.Duration) error {
	fmt.Printf("Set: (%s: %v)\n", key, object)
	sOp := r.clusterCL.Set(key, object, ttl)
	if sOp.Err() != nil {
		return sOp.Err()
	}
	return nil
}

// Get gets the object for the given key.
func (r *redisClient) Get(key string) (interface{}, error) {
	fmt.Printf("Get: (%s)\n", key)
	gOp := r.clusterCL.Get(key)
	if gOp.Err() != nil {
		return nil, gOp.Err()
	}
	res, err := gOp.Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete deletes the object and key for the given key from redis ring.
func (r *redisClient) Delete(key string) error {
	fmt.Printf("Delete: (%s)\n", key)
	dOp := r.clusterCL.Del(key)
	if dOp.Err() != nil {
		return dOp.Err()
	}
	_, err := dOp.Result()
	if err != nil {
		return err
	}
	return nil
}

// Exists reports whether object for the given key exists.
func (r *redisClient) Exists(key string) bool {
	fmt.Printf("Exists: (%s)\n", key)
	eOp := r.clusterCL.Exists(key)
	if eOp.Err() != nil {
		return false
	}
	_, err := eOp.Result()
	if err != nil {
		return false
	}
	return true
}

// Close closes the ring client, releasing any open resources.
func (r *redisClient) Close() error {
	fmt.Printf("Closing redis cluster client...\n")
	return r.clusterCL.Close()
}

func redisClusterOps() {
	addrs := []string{":7001", ":7002", ":7003", ":7004", ":7005", ":7006"}
	fmt.Println("redis cluster client...")

	cluster := NewRedisClient(addrs)
	// close...
	defer cluster.Close()

	key := "exampleKey"
	value := "exampleValue"
	ttl := 3 * time.Second

	// set...
	err := cluster.Set(key, value, ttl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// get...
	res, err := cluster.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}

	// exists...
	res = cluster.Exists(key)
	fmt.Printf("%v\n", res)

	// unknown key...
	res = cluster.Exists("randomNonExistentKey")
	fmt.Printf("%v\n", res)

	// expire
	fmt.Printf("sleeping for expire...\n")
	time.Sleep(5 * time.Second)

	// get...
	res, err = cluster.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}

	// exists...
	res = cluster.Exists(key)
	fmt.Printf("%v\n", res)

	// delete...
	// set...
	err = cluster.Set(key, value, ttl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	res = cluster.Delete(key)
	fmt.Printf("%v\n", res)

	// get...
	res, err = cluster.Get(key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("%v\n", res)
	}
}
