FROM ubuntu:20.04

LABEL maintainer="abdallah.yossf@gmail.com"

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /root

# Download hadoop 3.3.0 (this is the first RUN because it is the heaviest)
RUN apt-get update && apt-get upgrade -y
RUN apt-get install wget -y
RUN wget "http://archive.apache.org/dist/hadoop/common/hadoop-3.0.0/hadoop-3.0.0.tar.gz"

# Install java & utils
RUN apt-get install openjdk-8-jdk -y
RUN apt-get install vim nano net-tools -y
RUN apt-get install openssh-client openssh-server -y
RUN apt-get install iputils-ping -y
RUN apt-get install python3
RUN apt-get --assume-yes install python3-pip
RUN python3 -m pip install parquet-tools

RUN tar -xvzf hadoop-3.0.0.tar.gz  && \
  mv hadoop-3.0.0 /usr/local/hadoop && \
  rm hadoop-3.0.0.tar.gz




# set environment variable
ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64 
ENV HADOOP_HOME=/usr/local/hadoop 
# ENV HADOOP_ROOT_LOGGER=DEBUG,console
ENV PATH=$PATH:/usr/local/hadoop/bin:/usr/local/hadoop/sbin 


# ssh without key
# RUN ssh-keygen -t rsa -f ~/.ssh/id_rsa -P '' && \
#    cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys

#RUN mkdir -p ~/hdfs/namenode && \ 
#    mkdir -p ~/hdfs/datanode && \
#    mkdir $HADOOP_HOME/logs 

COPY config/* /tmp/
COPY ./messages.txt /root/
COPY ./mapredjob.jar /root/
COPY ./test-mapred.sh /root/

RUN mkdir ~/.ssh && \
    mv /tmp/ssh_config ~/.ssh/config && \
    mv /tmp/id_rsa ~/.ssh/id_rsa && \
    mv /tmp/id_rsa.pub ~/.ssh/id_rsa.pub && \
    cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
    
 
RUN mv /tmp/hadoop-env.sh /usr/local/hadoop/etc/hadoop/hadoop-env.sh && \
    mv /tmp/hdfs-site.xml $HADOOP_HOME/etc/hadoop/hdfs-site.xml && \ 
    mv /tmp/core-site.xml $HADOOP_HOME/etc/hadoop/core-site.xml && \
    mv /tmp/mapred-site.xml $HADOOP_HOME/etc/hadoop/mapred-site.xml && \
    mv /tmp/yarn-site.xml $HADOOP_HOME/etc/hadoop/yarn-site.xml && \
    mv /tmp/workers $HADOOP_HOME/etc/hadoop/workers && \
    mv /tmp/start-hadoop.sh ~/start-hadoop.sh && \
    mv /tmp/run-wordcount.sh ~/run-wordcount.sh

RUN chmod +x ~/start-hadoop.sh && \
    chmod +x ~/run-wordcount.sh && \
    chmod +x $HADOOP_HOME/sbin/start-dfs.sh && \
    chmod +x $HADOOP_HOME/sbin/start-yarn.sh && \
    chmod 600 ~/.ssh/id_rsa
    

    
# format namenode
#RUN /usr/local/hadoop/bin/hdfs namenode -format

CMD [ "sh", "-c", "service ssh start; bash"]




