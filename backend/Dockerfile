FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o oqu ./cmd/api/main.go

FROM alpine:3.23.3
COPY --from=builder app/migrations/ migrations
COPY --from=builder app/oqu /oqu
ENTRYPOINT ["/oqu"]