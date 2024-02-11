FROM golang:1.22.0 as builder
WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook-server cmd/webhook-server/main.go

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /build/webhook-server /app/webhook-server

CMD ["./webhook-server"]