FROM golang:1.20

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY go.mod .
COPY go.sum .

RUN apt update -y
RUN apt install -y postgresql-client

# Build the Go application
RUN go build -o userapi ./cmd/userapi/main.go

ADD docker/userapi/userapi-startup.sh .

ENTRYPOINT ["./userapi-startup.sh"]