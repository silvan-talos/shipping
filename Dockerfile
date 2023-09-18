# syntax=docker/dockerfile:1

FROM public.ecr.aws/docker/library/golang:1.21.1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN GOOS=linux go build -o shipping cmd/main.go

EXPOSE 8080

# Run
CMD ["./shipping"]
