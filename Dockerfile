############################
# STEP 1 build the image for creating the executable
############################
FROM golang:1.15.6-alpine3.12 as builder

# Install git + SSL ca certificates + make
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# Make to build go application
RUN apk update && apk add --no-cache git ca-certificates make unzip g++ && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser


WORKDIR /app
COPY . .

# Build the binary
RUN mkdir -p /app/plutus/
RUN touch /app/plutus/.emptyfile

# Compile the binary
RUN go build -o ./build/plutus

RUN mv ./build/* /app/plutus

############################
# STEP 2 build a small image with only the executable
############################
FROM alpine:3.12

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /bin/sh /bin/sh

# Copy our static executables
COPY --from=builder --chown=appuser:appuser /app/plutus /app/plutus

# Create a /tmp/ diretctory (required for go plugin for Unix Domain Socket)
COPY --from=builder --chown=appuser:appuser /app/plutus/.emptyfile /tmp/.emptyfile


# Use an unprivileged user.
USER appuser

WORKDIR /app/plutus

# Run the hello binary.
ENTRYPOINT ["./plutus", "start"]