FROM golang:1.22.0 as builder
WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bff cmd/bff/main.go

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /build/bff /app/bff
COPY --from=builder /build/.env /app/.env
COPY --from=builder /build/.crt /app/.crt

CMD ["./bff"]