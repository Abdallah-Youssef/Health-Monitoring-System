FROM ubuntu:20.04
ENV DEBIAN_FRONTEND noninteractive

WORKDIR /root

# Download hadoop 3.3.0 (this is the first RUN because it is the heaviest)
RUN apt-get update && apt-get upgrade -y
RUN apt-get install wget -y
RUN apt --assume-yes install python3
RUN apt --assume-yes install python3-pip

RUN python3 -m pip install -U flask-restful

RUN wget https://dlcdn.apache.org/spark/spark-3.2.1/spark-3.2.1-bin-hadoop3.2.tgz
RUN tar xf spark-3.2.1-bin-hadoop3.2.tgz
RUN mv spark-3.2.1-bin-hadoop3.2/ /opt/spark


# Install java & utils
RUN apt-get install openjdk-8-jdk -y
RUN apt-get install vim nano net-tools -y
RUN apt-get install openssh-client openssh-server -y
RUN apt-get install iputils-ping -y


# set environment variable
ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64 
ENV SPARK_HOME=/opt/spark
ENV PATH=$PATH:$SPARK_HOME/bin:$SPARK_HOME/sbin


# RUN mkdir ~/.ssh && \
#     mv /tmp/ssh_config ~/.ssh/config && \
#     mv /tmp/id_rsa ~/.ssh/id_rsa && \
#     mv /tmp/id_rsa.pub ~/.ssh/id_rsa.pub && \
#     cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys

EXPOSE 5000
EXPOSE 9000

COPY ./sparkjob.jar /root/    
COPY ./server.py /root/    
 
CMD [ "python3", "server.py"]




