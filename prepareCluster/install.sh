#!/bin/bash

# cleanup
for i in `seq 1 6`;
do
  docker rm -f redis$i
done
docker network rm redisCluster

#docker pull redis
docker network create redisCluster

for i in `seq 1 6`;
do
 docker run -d -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf -p 700$i:6379 --name "redis$i" --net redisCluster redis redis-server /usr/local/etc/redis/redis.conf;
done

# create cluster
nodes=$(for ind in `seq 1 6`; do echo -n "$(docker inspect -f '{{(index .NetworkSettings.Networks "redisCluster").IPAddress}}' "redis$ind")"':6379 '; done)
echo "yes" | docker exec -i redis1 redis-cli --cluster create $nodes --cluster-replicas 1

echo ""
docker exec redis1 redis-cli cluster nodes
