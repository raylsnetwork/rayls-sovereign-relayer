# Build context must be the repo root:
#   docker build -f public-relayer/Dockerfile .
#
# The build stage always runs on the native architecture and cross-compiles to
# the target. Go with CGO_ENABLED=0 cross-compiles cleanly, so arm64 images do
# not pay for QEMU emulation.
FROM --platform=$BUILDPLATFORM golang:1.26.4-bookworm AS build

WORKDIR /app

# Dependencies first — this layer is reused until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build-pubrelayer

FROM scratch

WORKDIR /app

COPY --from=build /app/build/pubrelayer /app/public-relayer

COPY --from=build /app/contracts /app/contracts

COPY --from=build /app/private-relayer/repository/migrations /app/repository/migrations

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/public-relayer", "run"]
