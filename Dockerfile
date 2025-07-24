FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-w -s" .

FROM scratch

WORKDIR /app

COPY --from=builder /app/static ./static

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["/app/server"]