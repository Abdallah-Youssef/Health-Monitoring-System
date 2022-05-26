#!/bin/bash

start-dfs.sh
start-yarn.sh
sleep 10
hadoop jar /root/mapredjob.jar MapReduceJob /health/health_0.csv /output/p.parquet
rm -rf /root/p.parquet
hdfs dfs -get -f /output/p.parquet /root/p.parquet
parquet-tools show -n 10 /root/p.parquet
