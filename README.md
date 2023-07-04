### Setup Docker

```
# start redis and redis-sentinel
docker-compose up -d

# login to redis or redis-sentinel container to check
docker-exec -it <sentinel-01/redis-01> /bin/bash
```
