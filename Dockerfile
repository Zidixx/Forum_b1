FROM golang:alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o forum-app ./cmd/web/

FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /app/forum-app .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/sql ./sql
COPY --from=builder /app/tls ./tls

EXPOSE 8080

CMD ["./forum-app"]