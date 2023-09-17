# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN GOOS=linux go build -o shipping cmd/main.go

EXPOSE 8080

# Run
CMD ["./shipping"]
