FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN  cd cmd/app && go build -o /bin/app

FROM alpine:latest

COPY --from=builder /bin/app /app/news

RUN chmod +x /app/news

CMD ["/app/news"]