FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git curl make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make assets && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ohara .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/ohara .

EXPOSE 8080
ENTRYPOINT ["./ohara"]
CMD ["--data-dir", "/data"]
