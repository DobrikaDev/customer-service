FROM golang:1.25.1 AS builder

WORKDIR /workspace

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/customer-service .

FROM alpine:3.20

WORKDIR /app

RUN adduser -D -u 10001 appuser

COPY --from=builder /out/customer-service /app/customer-service
COPY deployments /app/deployments

USER appuser

EXPOSE 8082

ENTRYPOINT ["./customer-service"]

