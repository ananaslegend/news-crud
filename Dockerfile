FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN  cd cmd/news && go build -o /bin/news

FROM alpine:latest

COPY --from=builder /bin/news /app/news

RUN chmod +x /app/news

CMD ["/app/news"]