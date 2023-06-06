FROM golang:1.20

WORKDIR /app

COPY cmd cmd
COPY go.mod .
COPY go.sum .

# Build the Go application
RUN go build -o userapi ./cmd/userapi.go

RUN apt update -y
RUN apt install -y postgresql-client

ADD bin/userapi-startup.sh .

ENTRYPOINT ["./userapi-startup.sh"]