# First time build
    cd docker
    docker build -t abdallahyossf/hadoop ./hadoop
    docker build -t abdallahyossf/healthmonitor ./healthmonitor
    docker build -t abdallahyossf/mockservice ./mockservice
    docker build -t abdallahyossf/djangoweb ./djangoweb
---
  Run this only once after your first hadoop build and don't run it again:

    docker tag abdallahyossf/hadoop abdallahyossf/hadoop:base

# Deploy
    docker swarm init
    docker network create -d overlay --attachable main

### hadoop cluster
    docker stack deploy -c docker-compose.hadoop.yml hadoop

### backend cluster (Health monitor + django)
    docker stack deploy -c docker-compose.backend.yml backend

### healthmonior
    docker service create --name healthmonitor --hostname healthmonitor --network main --mount type=bind,source="$(pwd)"/data,target=/root/data abdallahyossf/healthmonitor sleep 10d

  any service that needs its ip to be resolved by its service name should have a `--hostname` property, i think

### mockservice
    docker service create --name mockservice --network main abdallahyossf/mockservice sleep 10d
### djangoweb
    docker service create --name djangoweb --hostname djangoweb --network main -p 8001:8001 abdallahyossf/djangoweb
