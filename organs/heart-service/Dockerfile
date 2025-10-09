# File: organs/heart-service/Dockerfile

# --- Build Stage ---
# This stage compiles our Go application into a static binary.
FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/heart-service .

# --- Final Stage ---
# This stage creates the minimal final image for our organ.
FROM alpine:latest

# It's good practice to run containers as a non-root user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=builder /app/heart-service /heart-service
EXPOSE 8080
CMD ["/heart-service"]
