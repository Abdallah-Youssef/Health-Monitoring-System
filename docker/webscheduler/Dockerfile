FROM abdallahyossf/hadoop

WORKDIR /app

COPY ./HealthMessage.java ./
COPY ./MapReduceJob.java ./

RUN mkdir ./lib
COPY ./java-json.jar ./lib/

RUN javac -classpath $(hadoop classpath):/app/lib/* -d ./ ./MapReduceJob.java ./HealthMessage.java
RUN jar -cvf HealthMessagesJob.jar -C /app .

# RUN hadoop jar HealthMessagesJob.jar MapReduceJob /input.txt /output