FROM golang:1.22.0 as builder
WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o custom-controller cmd/custom-controller/main.go

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /build/custom-controller /app/custom-controller
COPY --from=builder /build/.kube/ /app/.kube/

CMD ["/app/custom-controller"]