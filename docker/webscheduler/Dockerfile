# I use a specific tag so that when hadoop:latest changes, rebuilding this file can still use the cached hadoop:base
FROM abdallahyossf/hadoop:base

RUN apt install python3-pip -y
# setup environment variable 
ENV DockerHOME=/home/app/webapp  

## Compiling HealthMessagesJob.jar
WORKDIR /mapreduce
COPY ./mapreduce/ ./

RUN javac -classpath $(hadoop classpath):/mapreduce/lib/* -d ./ ./MapReduceJob.java ./HealthMessage.java
RUN jar -cvf HealthMessagesJob.jar -C /mapreduce .
# RUN hadoop jar HealthMessagesJob.jar MapReduceJob /messages.txt /temp /temp/part-r-00000 /output


## Django web server
WORKDIR /django_app

# set environment variables  
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1  

# install dependencies  
RUN pip install --upgrade pip
# Copy the app
COPY ./django_app/ ./

#  install all dependencies  
RUN pip install -r requirements.txt 
# Django app port  
EXPOSE 8000
# start server  
CMD python manage.py runserver  