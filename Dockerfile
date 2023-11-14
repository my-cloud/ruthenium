FROM golang:1.19 as builder
WORKDIR /app
COPY ./cmd ./cmd
COPY ./config ./config
COPY ./domain ./domain
COPY ./infrastructure ./infrastructure
ADD go.mod .
ADD go.sum .

RUN CGO_ENABLED=0 go build -o validatornode cmd/validatornode/main.go
RUN CGO_ENABLED=0 go build -o observernode cmd/observernode/main.go

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /app/config /app/config
COPY --from=builder /app/validatornode /app
COPY --from=builder /app/observernode /app
