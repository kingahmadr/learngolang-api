# Step 1: Use the official Golang image as the base image
FROM golang:1.20-alpine as builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Step 4: Copy the source code into the container
COPY . .

# Step 5: Build the Go application
RUN go build -o main .

# Step 6: Use a smaller image to reduce the size of the final image
FROM alpine:latest  

# Step 7: Set the Current Working Directory inside the container
WORKDIR /root/

# Step 8: Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Step 9: Expose the port the app runs on
EXPOSE 8080

# Step 10: Command to run the executable
CMD ["./main"]