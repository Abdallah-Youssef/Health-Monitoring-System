# I use a specific tag so that when hadoop:latest changes, rebuilding this file can still use the cached hadoop:base
FROM abdallahyossf/hadoop:base

RUN apt install python3-pip -y
RUN apt install maven -y


## Compiling HealthMessagesJob.jar
WORKDIR /mapreduce
COPY ./mapreduce/ ./
ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64 
RUN mvn package
# RUN javac -classpath $(hadoop classpath):/mapreduce/lib/* -d ./ ./MapReduceJob.java ./HealthMessage.java ./ThriftHealthMessage.java
# RUN jar -cvf HealthMessagesJob.jar -C /mapreduce .
# RUN hadoop jar HealthMessagesJob.jar MapReduceJob /messages.txt  /output

CMD ["sleep", "1d"]

# ## Django web server
# WORKDIR /django_app

# # set environment variables  
# ENV PYTHONDONTWRITEBYTECODE 1
# ENV PYTHONUNBUFFERED 1  

# # install dependencies  
# RUN pip install --upgrade pip
# # Copy the app
# COPY ./django_app/ ./

# #  install all dependencies  
# RUN pip install -r requirements.txt 
# # Django app port  
# EXPOSE 8000
# # start server  
# CMD python manage.py runserver  