FROM golang:1.18 as builder
WORKDIR /app
COPY ./src ./src
COPY ./test ./test
ADD go.mod .
ADD go.sum .


RUN CGO_ENABLED=0 go build -o node src/node/main.go
RUN CGO_ENABLED=0 go build -o ui src/ui/main.go


FROM gcr.io/distroless/static-debian11
WORKDIR /app/templates
COPY --from=builder /app/src/ui/templates .
WORKDIR /app
COPY --from=builder /app/node /app
COPY --from=builder /app/ui /app


