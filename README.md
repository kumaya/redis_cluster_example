# redis cluster and cluster client

redis_cluster_example repository demonstrates a way to 
 - create redis cluster and
 - use of redis cluster, using go redis cluster client.

#### Prerequisites

```
docker
go
```

#### Create a redis cluster

To create a redis cluster run
```bash
./prepareCluster/install.sh
```

This will create a cluster of 3 master nodes and 3 slave nodes.

#### Run examle go redis cluster client

```go
go build -o runner *.go
./runner
```

#### To update packages

```
export GO111MODULE=on; go mod tidy; go mod vendor; unset GO111MODULE
```
