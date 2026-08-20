# Build context must be the repo root:
#   docker build -f cts/Dockerfile .
FROM golang:1.26.4-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 make build-cts

FROM scratch

WORKDIR /app

COPY --from=build /app/build/cts /app/rayls-cts

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/rayls-cts", "run"]
