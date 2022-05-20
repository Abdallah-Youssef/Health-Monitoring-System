FROM golang AS build

LABEL maintainer="abdallah.yossf@gmail.com"
WORKDIR /app

COPY ./main.go ./

# Download the required modules
RUN go mod init healthmonitor
RUN go mod tidy

# Build the executable
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./healthmonitor  


FROM alpine
COPY --from=build /app/healthmonitor /
CMD ["/healthmonitor"]





