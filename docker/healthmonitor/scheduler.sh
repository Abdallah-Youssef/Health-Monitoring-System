#!/bin/sh

# Start mockserver listener
/root/server &

while true
do
  hadoop jar /root/mapreducejob.jar MapReduceJob /messages.csv /output

  rm  /root/data/part-r-00000.parquet
  hdfs dfs -get /output/part-r-00000.parquet /root/data/batch/
  hdfs dfs -truncate 0 /new.csv 

  scp /root/data/batch/part-r-00000.parquet root@djangoweb:/django_app/django_app/output/batch.parquet
  sleep 30m
done &

while true
do
  # run spark job
  # get output from hdfs
  # put in djangoweb's container
  sleep 5m
done &

