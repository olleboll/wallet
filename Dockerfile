FROM golang:1.21

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build cmd/main.go

EXPOSE 8000