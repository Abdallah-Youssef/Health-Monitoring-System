FROM golang:alpine3.15 AS build

LABEL maintainer="abdallah.yossf@gmail.com"
WORKDIR /app

COPY ./main.go ./


RUN go mod init mockservice
RUN go mod tidy

# Build the executable
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./mockservice  


FROM alpine
COPY --from=build /app/mockservice /
CMD ["/mockservice"]





