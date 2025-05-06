FROM golang:1.23.5 AS builder

# where to application will be built
WORKDIR /app 
# copy go mod and sum files to install dependencies
COPY go.mod ./ 
# install new dependencies
RUN go mod download
# copy the rest of the application code
COPY . .

# CGO_ENABLED=0 <- ensures Go program can run on any linux system
# GOOS=linux <- sets the target OS to linux
# GOARCH=amd64 <- sets the target architecture to 64-bit
# go build -o main . <- builds the application and outputs it to a file named main
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# set the working directory in the container
FROM alpine:latest

# set the working directory in the container
WORKDIR /root/
# copy the binary from the builder stage to the current stage
COPY --from=builder /app/main .
# copy the .env file to the current stage
COPY .env .

# set the permissions for the main binary
RUN chmod +x /root/main
# expose port 8080 for the application
EXPOSE 8080
# set the command to run the application
CMD ["/root/main"]