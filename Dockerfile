# Build stage
FROM golang:1.16-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build aviatus/vividora

# Runtime stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the build stage
COPY --from=build /app/vividora .

# Expose a port for the application (if applicable)
EXPOSE 8080

# Specify the command to run the application
CMD ["./vividora"]
