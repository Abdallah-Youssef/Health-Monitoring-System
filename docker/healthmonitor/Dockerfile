FROM golang:alpine3.15 AS build

LABEL maintainer="abdallah.yossf@gmail.com"

RUN apk update
RUN apk add nano iputils

WORKDIR /app

COPY ./main.go ./
COPY ./healthmonitor.go ./
COPY ./httpserver.go ./
COPY ./structs.go ./

RUN mkdir resources 
COPY ./resources/* ./resources/

# Download the required modules
RUN go mod init healthmonitor
RUN go mod tidy

# # Build the executable
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o healthmonitor


# FROM alpine
# COPY --from=build /app/healthmonitor /
# CMD ["/healthmonitor"]





