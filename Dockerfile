FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test ./... -v
RUN CGO_ENABLED=0 GOOS=linux go build -o currency-converter ./cmd/app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/currency-converter .
CMD ["./currency-converter"]
