# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/api

# Run stage (migrations + seed data are embedded in the binary)
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

COPY --from=build /out/server /app/server

WORKDIR /app
EXPOSE 8080
CMD ["./server"]
