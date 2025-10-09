# File: organs/brain-service/Dockerfile

# --- Build Stage ---
FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/brain-service .

# --- Final Stage ---
FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /app/brain-service /brain-service
EXPOSE 8080
CMD ["/brain-service"]
