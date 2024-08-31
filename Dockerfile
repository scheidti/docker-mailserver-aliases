FROM golang:1.23.0-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=builder /app/docker-mailserver-aliases /app/docker-mailserver-aliases
ENV GIN_MODE=release
ENTRYPOINT ["/app/docker-mailserver-aliases"]