# Start from the official golang image
FROM golang:latest

# Set environment variables
ENV PROJECT_DIR=/go/src/app \
    GO111MODULE=on \
    CGO_ENABLED=0

# Set the working directory
WORKDIR $PROJECT_DIR

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Install CompileDaemon
RUN go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon

# Command to run the application using CompileDaemon
CMD CompileDaemon -build="go build -o /go/src/app/build/app" -command="/go/src/app/build/app"
