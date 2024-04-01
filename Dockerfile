FROM golang:1.19 as builder
WORKDIR /app
COPY ./validatornode ./validatornode
COPY ./observernode ./observernode
COPY ./config ./config
ADD go.mod .
ADD go.sum .

RUN CGO_ENABLED=0 go build -o validatornode validatornode/main.go
RUN CGO_ENABLED=0 go build -o observernode observernode/main.go

FROM gcr.io/distroless/static-debian11
USER nonroot
WORKDIR /app
COPY --from=builder /app/config /app/config
COPY --from=builder /app/validatornode /app
COPY --from=builder /app/observernode /app
