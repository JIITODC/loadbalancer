# loadbalancer

## How to run locally 
This command will spwan multiple toy load servers on docker 
```
docker compose up
```

Now to run loadbalancer fo 
```
go run ./docker-lb
```

Now you can ping 
```
curl localhost:8000
```
