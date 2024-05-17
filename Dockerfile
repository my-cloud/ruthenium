FROM golang:1.19 AS builder
WORKDIR /app
COPY ./validatornode ./validatornode
COPY ./accessnode ./accessnode
ADD go.mod .
ADD go.sum .

RUN CGO_ENABLED=0 go build -o validatornode validatornode/main.go
RUN CGO_ENABLED=0 go build -o accessnode accessnode/main.go

FROM debian:11.9
USER nonroot
WORKDIR /app
COPY --from=builder /app/validatornode /app
COPY --from=builder /app/accessnode /app
