FROM golang:alpine AS builder

WORKDIR /build

RUN apk update && apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    && rm -rf /var/cache/apk/*

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/scheduler/main.go

FROM alpine:latest AS final

WORKDIR /app

RUN apk update && apk add --no-cache \
    sqlite \
    && rm -rf /var/cache/apk/*

COPY --from=builder /build/main ./
COPY --from=builder /build/web ./web
COPY --from=builder /build/config.yaml ./

EXPOSE 7540

CMD ["./main"]