FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /task-service

RUN go install

EXPOSE 9090

CMD ["/task-service"]
