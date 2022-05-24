FROM python:3.10.0a7-alpine3.13



## Django web server
WORKDIR /django_app

# set environment variables  
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1  

# Copy the app
COPY ./django_app/ ./
#  install all dependencies  
RUN pip3 install -r requirements.txt
#install duckDB
RUN pip3 install duckdb==0.3.4 
# Django app port  
EXPOSE 8001
# start server  
CMD ["python3", "./manage.py", "runserver"]
