FROM golang:1.21.1-alpine

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /api cmd/api/main.go

EXPOSE 8083

ENTRYPOINT ["/api"]
