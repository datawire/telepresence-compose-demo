FROM golang:1.20

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY go.mod .
COPY go.sum .

RUN apt update -y
RUN apt install -y postgresql-client

# Build the Go application
RUN go build -o contactapi ./cmd/contactapi/main.go

ADD docker/contactapi/contactapi-startup.sh .

ENTRYPOINT ["./contactapi-startup.sh"]