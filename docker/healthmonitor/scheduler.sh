#!/bin/sh

# Start mockserver listener
/root/server &

while true
do
  hadoop jar /root/mapreducejob.jar MapReduceJob /messages.csv /output

  rm  /root/data/part-r-00000.parquet
  hdfs dfs -get /output/part-r-00000.parquet /root/data/
  
  scp /root/data/part-r-00000.parquet root@djangoweb:/django_app/django_app/output/batch.parquet
  sleep 30m
done

