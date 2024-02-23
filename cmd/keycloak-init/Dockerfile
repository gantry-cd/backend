FROM golang:1.22-alpine as builder
WORKDIR /build

COPY . .

RUN go build -o /build/keycloak-init cmd/keycloak-init/main.go

FROM scratch
COPY --from=builder /build/keycloak-init /usr/local/bin/keycloak-init

ENV KEYCLOAK_HOST=
ENV KEYCLOAK_ADMIN_ID=
ENV KEYCLOAK_ADMIN_PASSWORD=


ENTRYPOINT ["/usr/local/bin/keycloak-init"]