# I use a specific tag so that when hadoop:latest changes, rebuilding this file can still use the cached hadoop:base
FROM abdallahyossf/hadoop:base

RUN apt install python3-pip -y
RUN pip3 install flask


## Compiling HealthMessagesJob.jar
WORKDIR /mapreduce
COPY ./mapreduce/ ./

RUN javac -classpath $(hadoop classpath):/mapreduce/lib/* -d ./ ./MapReduceJob.java ./HealthMessage.java
RUN jar -cvf HealthMessagesJob.jar -C /mapreduce .
# RUN hadoop jar HealthMessagesJob.jar MapReduceJob /messages.txt /temp /temp/part-r-00000 /output


## Flask web server
WORKDIR /flask
COPY ./flask/ ./
CMD ["python3", "app.py"]