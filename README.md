# Getting started

## Prerequisite
### Setup Docker

```
# start redis and redis-sentinel
docker-compose up -d

# login to redis or redis-sentinel container to check
docker-exec -it <sentinel-01/redis-01> /bin/bash
```

### Development

Ensure docker-compose is up and running, then run `make dev-run`

#### Manual cURL Testing

1. /sentinel/register - POST

```
curl -X POST http://localhost:8282/sentinel/register  \
   -H "Content-Type: application/json" \
   -d '{"name": "test-sentinel", "hosts": "10.218.123.41:26379,10.218.123.42:26379,10.218.123.43:26379"}' 
```

2. /sentinel - GET (all)

```
curl -X GET http://localhost:8282/sentinel \
   -H "Content-Type: application/json"
```

3. /sentinel/<id>

```
curl -X GET http://localhost:8282/sentinel/<id> \
   -H "Content-Type: application/json"

```