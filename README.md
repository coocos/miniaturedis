# miniaturedis

A miniature implementation of [Redis](https://redis.io).

## Description

Miniaturedis implements a tiny Redis server based on the the [RESP protocol specification](https://redis.io/docs/reference/protocol-spec/), i.e. you can use `redis-cli` or any other Redis client implementation to connect to it and issue commands:

```shell
# Start miniaturedis in the background
$ ./miniaturedis &

# Start a Redis client
$ redis-cli
127.0.0.1:6379> GET some-key
(nil)
127.0.0.1:6379> SET some-key some-value
OK
127.0.0.1:6379> GET some-key
"some-value"
127.0.0.1:6379> GET some-key
```

## Features

As the name suggests, Miniaturedis currently only supports a tiny subset of Redis commands:

* [GET](https://redis.io/commands/get/) (with no options)
* [SET](https://redis.io/commands/set/) (with no options)

There is also no snapshot behaviour, i.e. nothing is persisted.

## How to build

Assuming you have Go installed:

```shell
go build -o . ./...
./miniaturedis
```