# telnet-chat server

telnet chat server implemented in golang.

the server works as independent or with redis. the server can horizontally scale if working with redis backend.

## requiremens

* go 1.6+

## getting started

type the following command to compile

```
go build .
./telnet-chat
```

## modes

this chat server works as standalone or with redis backend

### standalone

server works as standalone by default

```
./telnet-chat
```

### redis backend

server works with redis if redis host is specified

```
REDIS_HOST=localhost:6379 ./telnet-chat
```

