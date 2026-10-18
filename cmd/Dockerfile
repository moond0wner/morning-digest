FROM golang:1.27.2-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/exe /app/cmd

FROM alpine:3.24.2
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /app/cmd/exe /app
COPY --from=builder /app/config.yaml /app/config.yaml
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]
