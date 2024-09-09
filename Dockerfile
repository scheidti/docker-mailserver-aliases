FROM node:20 AS frontend-builder

WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend .
RUN npm run build

FROM golang:1.23.0-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./

COPY . .
COPY --from=frontend-builder /app/dist /app/frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=builder /app/docker-mailserver-aliases /app/docker-mailserver-aliases
ENV GIN_MODE=release
ENTRYPOINT ["/app/docker-mailserver-aliases"]