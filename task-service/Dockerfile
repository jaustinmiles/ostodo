FROM golang:1.19-alpine

LABEL base.name="ostodo.task-service"

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o task-service .


EXPOSE 9090

ENTRYPOINT ["./task-service"]
