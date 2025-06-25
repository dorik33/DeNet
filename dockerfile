FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine

WORKDIR /app  

COPY --from=builder /build/main /app/main
COPY --from=builder /build/.env  /app/.env
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /go/bin/goose /app/goose

CMD ["./main"]