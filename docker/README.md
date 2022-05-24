# First time build
    cd docker
    docker build -t abdallahyossf/hadoop ./hadoop
    docker build -t abdallahyossf/healthmonitor ./healthmonitor
    docker build -t abdallahyossf/mockservice ./mockservice
    docker build -t abdallahyossf/scheduler ./webscheduler
    docker build -t abdallahyossf/djangoweb ./djangoweb
---
  Run this only once after your first hadoop build and don't run it again:

    docker tag abdallahyossf/hadoop abdallahyossf/hadoop:base

# Deploy
    docker swarm init
    docker network create -d overlay --attachable main

### hadoop cluster
    docker stack deploy -c docker-compose.hadoop.yml hadoop

### healthmonior
    docker service create --name healthmonitor --hostname healthmonitor --network main abdallahyossf/healthmonitor sleep 10d

  any service that needs its ip to be resolved by its service name should have a `--hostname` property

### mockservice
    docker service create --name mockservice --network main abdallahyossf/mockservice sleep 10d


### webserver / scheduler
    docker service create --name scheduler --network main  -p 5000:5000 abdallahyossf/scheduler

