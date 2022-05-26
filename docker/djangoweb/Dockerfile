FROM abdallahyossf/hadoop:base
#install 
RUN apt install python3-pip -y
## Django web server
WORKDIR /django_app

# set environment variables  
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1  

# Copy the app
COPY ./django_app/ ./
#  install all dependencies  
RUN pip3 install -r requirements.txt
#install numpy
RUN pip3 install numpy 
#install duckDB
RUN pip3 install duckdb==0.3.4 
# Django app port  
EXPOSE 8001
# start server  
CMD ["python3", "./manage.py", "runserver"]
