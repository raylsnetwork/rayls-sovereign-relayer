# Build context must be the repo root:
#   docker build -f cts/Dockerfile .
#
# The build stage always runs on the native architecture and cross-compiles to
# the target. Go with CGO_ENABLED=0 cross-compiles cleanly, so arm64 images do
# not pay for QEMU emulation.
FROM --platform=$BUILDPLATFORM golang:1.26.4-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build-cts

FROM scratch

WORKDIR /app

COPY --from=build /app/build/cts /app/rayls-cts

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/rayls-cts", "run"]
