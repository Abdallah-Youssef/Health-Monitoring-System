FROM golang:alpine3.15 AS go-build

LABEL maintainer="abdallah.yossf@gmail.com"

RUN apk update
RUN apk add nano iputils

WORKDIR /healthserver
COPY ./healthserver/ ./


# Download the required modules
RUN go mod init healthserver
RUN go mod tidy

# Build executable
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server


FROM maven:alpine as maven-build
WORKDIR /mapreduce

# Download packages (separate steps to allow caching the packages)
COPY ./mapreduce/pom.xml ./
RUN mvn install

# Copy the source code and create the jar
COPY ./mapreduce/src/ ./src
RUN mvn package

FROM abdallahyossf/hadoop
COPY --from=go-build /healthserver/server /root
COPY --from=maven-build /mapreduce/target/json-to-parquet-0.0.1-SNAPSHOT.jar /root/mapreducejob.jar
COPY ./scheduler.sh ./



# this runs the server that listens for healthmessages from the mockservices
# CMD ["/root/server"] 

# this runs the mapreduce job, which generates a ".parquet" file in the hdfs /output directory
# CMD ["hadoop", "jar", "/root/mapreducejob.jar", "MapReduceJob", "/messages.txt",  "/output"]

# TODO: create some script that:
# 1: starts the /root/server in the background
# 2: runs the map reduce job periodically
# 3: runs the spark job periodically
# 4: update the databases using the previous outputs
CMD [ "sh", "-c", "service ssh start; sleep 10d"]



