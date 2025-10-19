# Use a specific minor version for stability
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies before copying the rest of the code
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o bin .

# Expose the required port
EXPOSE 5000

# Ensure logs can be written to the mounted volume (No changes needed)
ENTRYPOINT [ "/app/bin" ]