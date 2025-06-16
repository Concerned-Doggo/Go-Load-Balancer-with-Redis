# Dockerfile

# Use official Go image
FROM golang:1.24.3

# Set working directory inside the container
WORKDIR /app

# Copy code from current directory into container
COPY . . 

# Build the Go app into an executable named "main"
RUN go build -o main .

# Expose port 8080 (for HTTP traffic)
EXPOSE 8080

# Run the built executable
CMD ["./main"]


