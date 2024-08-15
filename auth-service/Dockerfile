FROM golang:alpine3.19

USER root
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN echo environment=production > configs/.env

RUN go build cmd/api/main.go

ENTRYPOINT [ "./main" ]