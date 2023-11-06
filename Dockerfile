FROM golang:1.21.3-alpine

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the files for building
COPY ./ ./

RUN go build -o /app/ ./...

# Use more minimal base image for final binary
FROM scratch

EXPOSE 8080

# Copy created binary from first stage
COPY --from=0 /app/battlesnake-server /bin/battlesnake-server

CMD ["/bin/battlesnake-server"]