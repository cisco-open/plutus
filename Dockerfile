# Copyright 2022 Cisco Systems, Inc. and its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

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