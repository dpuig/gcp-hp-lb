# Use a more recent Golang version for up-to-date features and security
FROM golang:1.20-alpine as builder

# Add Maintainer Info (Always a good practice)
LABEL maintainer="Daniel Puig Gerarde <dpuigerarde@gmail.com>"

# Install minimal dependencies
RUN apk --no-cache add ca-certificates tzdata

# Workdir optimization
WORKDIR /app

# Copy go.mod and go.sum for efficient dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the optimized Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o main .

### Start a new stage ######
FROM scratch

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the necessary certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/

WORKDIR /app

EXPOSE 8080

CMD ["./main"] 