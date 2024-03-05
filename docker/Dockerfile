# Use a more recent Golang version for up-to-date features and security
FROM golang:1.20-alpine as builder

# Add Maintainer Info (Always a good practice)
LABEL maintainer="Daniel Puig Gerarde <dpuigerarde@gmail.com>"

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

# Install minimal dependencies for the smallest image
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"] 
