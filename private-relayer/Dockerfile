# Build context must be the repo root:
#   docker build -f private-relayer/Dockerfile .
FROM golang:1.26.4-bookworm AS build

WORKDIR /app

COPY . .
RUN rm -rf go.sum
RUN go clean -modcache
RUN ls
RUN go mod tidy
RUN CGO_ENABLED=0 make

FROM scratch

# Create the /app directory in the container
WORKDIR /app

# Copy the binary from the build stage to the current stage
COPY --from=build /app/build/relayer /app/rayls-sovereign-relayer

# Coppy the migrations folder and files
COPY --from=build /app/private-relayer/repository/migrations /app/repository/migrations

# Copy the contract bindings (used for custom error decoding)
COPY --from=build /app/contracts /app/contracts

# Copy the SSL certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set the binary as the entry point of the container
ENTRYPOINT ["/app/rayls-sovereign-relayer", "run"]
