FROM golang:1.22-alpine AS builder

# Install git and SSH client
RUN apk add --no-cache git openssh-client

# Install Terraform
RUN apk add --no-cache curl \
    && curl -LO https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip \
    && unzip terraform_1.5.7_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && rm terraform_1.5.7_linux_amd64.zip

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/deployer

# Start a new stage from scratch
FROM alpine:latest

# Install necessary tools
RUN apk add --no-cache ca-certificates git openssh-client

# Copy Terraform from the builder stage
COPY --from=builder /usr/local/bin/terraform /usr/local/bin/terraform

# Copy the built application from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory
WORKDIR /app

# Run the application
CMD ["./main"]
