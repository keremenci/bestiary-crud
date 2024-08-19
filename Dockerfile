# build stage
FROM golang:alpine AS builder

# Install git and curl for dependencies
RUN apk add --no-cache git curl

WORKDIR /app

# Install migrate tool
RUN GOBIN=/app go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

RUN go build -o main .

# run stage
FROM alpine:latest AS run

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Copy the migrate tool and migrations
COPY --from=builder /app/migrate .
COPY --from=builder /app/db/migrations ./db/migrations

# Copy config
COPY --from=builder /app/config/config.yml ./config/config.yml

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run migrations followed by the application
CMD ["sh", "-c", "/root/migrate -database $BESTIARY_DATABASE_URL -path db/migrations up && ./main"]