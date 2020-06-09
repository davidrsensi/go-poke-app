




# Set the latest golang base image:
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /go-poke-app

# Build server:
COPY  ./ ./
RUN go mod download
RUN go build -o go-poke-app .

# Expose server port:
EXPOSE 8080

# App version (uncomment if the image is for non-development ends):
# ENV APP_VERSION "YOUR_VERSION"

# Run the server:
CMD ["./go-poke-app"]