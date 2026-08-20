# Third-Party Licenses

`rayls-privacy-relayer-api` is licensed under the Apache License, Version 2.0 (see
[LICENSE](./LICENSE) and [NOTICE](./NOTICE)).

It is built with the third-party Go modules listed below, which are compiled into the
`private-relayer`, `public-relayer`, and `cts` binaries. Each remains under its own license;
the authoritative license text and copyright for each is at the linked source for the exact
version used. This inventory was generated with `go-licenses` over the build graph and
reflects the modules resolved by `go.mod`/`go.sum`.

Copyleft note: **github.com/ethereum/go-ethereum** is used under the **LGPL-3.0** (its library
packages only — none of its GPL-3.0 `cmd/` programs are imported). Because the binaries are
statically linked, recipients may relink against a modified go-ethereum; its full license text
is bundled at [licenses/LGPL-3.0.txt](./licenses/LGPL-3.0.txt). All other dependencies are
permissive (Apache-2.0, MIT, BSD, ISC, CC0/Public Domain, Unlicense).

## LGPL-3.0 (library packages only) (1)

- **github.com/ethereum/go-ethereum** — https://github.com/ethereum/go-ethereum/blob/v1.16.9/COPYING.LESSER

## Apache-2.0 (74)

- **cloud.google.com/go/auth** — https://github.com/googleapis/google-cloud-go/blob/auth/v0.18.1/auth/LICENSE
- **cloud.google.com/go/auth/oauth2adapt** — https://github.com/googleapis/google-cloud-go/blob/auth/oauth2adapt/v0.2.8/auth/oauth2adapt/LICENSE
- **cloud.google.com/go/compute/metadata** — https://github.com/googleapis/google-cloud-go/blob/compute/metadata/v0.9.0/compute/metadata/LICENSE
- **cloud.google.com/go/iam** — https://github.com/googleapis/google-cloud-go/blob/iam/v1.5.3/iam/LICENSE
- **cloud.google.com/go/kms** — https://github.com/googleapis/google-cloud-go/blob/kms/v1.25.0/kms/LICENSE
- **cloud.google.com/go/longrunning** — https://github.com/googleapis/google-cloud-go/blob/longrunning/v0.8.0/longrunning/LICENSE
- **github.com/agoda-com/opentelemetry-go/otelslog** — https://github.com/agoda-com/opentelemetry-go/blob/otelslog/v0.3.0/otelslog/LICENSE
- **github.com/agoda-com/opentelemetry-logs-go** — https://github.com/agoda-com/opentelemetry-logs-go/blob/v0.6.0/LICENSE
- **github.com/aws/aws-sdk-go-v2** — https://github.com/aws/aws-sdk-go-v2/blob/v1.41.6/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/config** — https://github.com/aws/aws-sdk-go-v2/blob/config/v1.18.45/config/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/credentials** — https://github.com/aws/aws-sdk-go-v2/blob/credentials/v1.13.43/credentials/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/feature/ec2/imds** — https://github.com/aws/aws-sdk-go-v2/blob/feature/ec2/imds/v1.13.13/feature/ec2/imds/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/internal/configsources** — https://github.com/aws/aws-sdk-go-v2/blob/internal/configsources/v1.4.22/internal/configsources/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/internal/endpoints/v2** — https://github.com/aws/aws-sdk-go-v2/blob/internal/endpoints/v2.7.22/internal/endpoints/v2/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/internal/ini** — https://github.com/aws/aws-sdk-go-v2/blob/internal/ini/v1.3.45/internal/ini/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/service/internal/presigned-url** — https://github.com/aws/aws-sdk-go-v2/blob/service/internal/presigned-url/v1.9.37/service/internal/presigned-url/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/service/kms** — https://github.com/aws/aws-sdk-go-v2/blob/service/kms/v1.51.0/service/kms/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/service/sso** — https://github.com/aws/aws-sdk-go-v2/blob/service/sso/v1.15.2/service/sso/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/service/ssooidc** — https://github.com/aws/aws-sdk-go-v2/blob/service/ssooidc/v1.17.3/service/ssooidc/LICENSE.txt
- **github.com/aws/aws-sdk-go-v2/service/sts** — https://github.com/aws/aws-sdk-go-v2/blob/service/sts/v1.23.2/service/sts/LICENSE.txt
- **github.com/aws/smithy-go** — https://github.com/aws/smithy-go/blob/v1.25.0/LICENSE
- **github.com/cockroachdb/errors** — https://github.com/cockroachdb/errors/blob/v1.12.0/LICENSE
- **github.com/cockroachdb/logtags** — https://github.com/cockroachdb/logtags/blob/bb51bb14a506/LICENSE
- **github.com/cockroachdb/redact** — https://github.com/cockroachdb/redact/blob/v1.1.6/LICENSE
- **github.com/consensys/gnark-crypto** — https://github.com/consensys/gnark-crypto/blob/v0.19.2/LICENSE
- **github.com/containerd/errdefs** — https://github.com/containerd/errdefs/blob/v1.0.0/LICENSE
- **github.com/containerd/errdefs/pkg** — https://github.com/containerd/errdefs/blob/pkg/v0.3.0/pkg/LICENSE
- **github.com/containerd/log** — https://github.com/containerd/log/blob/v0.1.0/LICENSE
- **github.com/containerd/platforms** — https://github.com/containerd/platforms/blob/v0.2.1/LICENSE
- **github.com/crate-crypto/go-eth-kzg** — https://github.com/crate-crypto/go-eth-kzg/blob/v1.4.0/LICENSE
- **github.com/crate-crypto/go-ipa** — https://github.com/crate-crypto/go-ipa/blob/53bbb0ceb27a/LICENSE-APACHE
- **github.com/distribution/reference** — https://github.com/distribution/reference/blob/v0.6.0/LICENSE
- **github.com/docker/docker** — https://github.com/docker/docker/blob/v28.5.2/LICENSE
- **github.com/docker/go-connections** — https://github.com/docker/go-connections/blob/v0.6.0/LICENSE
- **github.com/docker/go-units** — https://github.com/docker/go-units/blob/v0.5.0/LICENSE
- **github.com/ebitengine/purego** — https://github.com/ebitengine/purego/blob/v0.9.1/LICENSE
- **github.com/go-logr/logr** — https://github.com/go-logr/logr/blob/v1.4.3/LICENSE
- **github.com/go-logr/stdr** — https://github.com/go-logr/stdr/blob/v1.2.2/LICENSE
- **github.com/google/s2a-go** — https://github.com/google/s2a-go/blob/v0.1.9/LICENSE.md
- **github.com/googleapis/enterprise-certificate-proxy/client** — https://github.com/googleapis/enterprise-certificate-proxy/blob/v0.3.11/LICENSE
- **github.com/iden3/go-iden3-crypto** — https://github.com/iden3/go-iden3-crypto/blob/v0.0.17/LICENSE-APACHE
- **github.com/klauspost/compress** — https://github.com/klauspost/compress/blob/v1.18.3/LICENSE
- **github.com/minio/sha256-simd** — https://github.com/minio/sha256-simd/blob/v1.0.1/LICENSE
- **github.com/moby/docker-image-spec/specs-go/v1** — https://github.com/moby/docker-image-spec/blob/v1.3.1/LICENSE
- **github.com/moby/go-archive** — https://github.com/moby/go-archive/blob/v0.2.0/LICENSE
- **github.com/moby/patternmatcher** — https://github.com/moby/patternmatcher/blob/v0.6.0/LICENSE
- **github.com/moby/sys/sequential** — https://github.com/moby/sys/blob/sequential/v0.6.0/sequential/LICENSE
- **github.com/moby/sys/user** — https://github.com/moby/sys/blob/user/v0.4.0/user/LICENSE
- **github.com/moby/term** — https://github.com/moby/term/blob/v0.5.2/LICENSE
- **github.com/nats-io/nats.go** — https://github.com/nats-io/nats.go/blob/v1.48.0/LICENSE
- **github.com/nats-io/nkeys** — https://github.com/nats-io/nkeys/blob/v0.4.14/LICENSE
- **github.com/nats-io/nuid** — https://github.com/nats-io/nuid/blob/v1.0.1/LICENSE
- **github.com/opencontainers/go-digest** — https://github.com/opencontainers/go-digest/blob/v1.0.0/LICENSE
- **github.com/opencontainers/image-spec/specs-go** — https://github.com/opencontainers/image-spec/blob/v1.1.1/LICENSE
- **github.com/spf13/afero** — https://github.com/spf13/afero/blob/v1.15.0/LICENSE.txt
- **github.com/spf13/cobra** — https://github.com/spf13/cobra/blob/v1.10.2/LICENSE.txt
- **go.opentelemetry.io/auto/sdk** — https://github.com/open-telemetry/opentelemetry-go-instrumentation/blob/sdk/v1.2.1/sdk/LICENSE
- **go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc** — https://github.com/open-telemetry/opentelemetry-go-contrib/blob/instrumentation/google.golang.org/grpc/otelgrpc/v0.61.0/instrumentation/google.golang.org/grpc/otelgrpc/LICENSE
- **go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp** — https://github.com/open-telemetry/opentelemetry-go-contrib/blob/instrumentation/net/http/otelhttp/v0.64.0/instrumentation/net/http/otelhttp/LICENSE
- **go.opentelemetry.io/contrib/instrumentation/runtime** — https://github.com/open-telemetry/opentelemetry-go-contrib/blob/instrumentation/runtime/v0.64.0/instrumentation/runtime/LICENSE
- **go.opentelemetry.io/otel** — https://github.com/open-telemetry/opentelemetry-go/blob/v1.43.0/LICENSE
- **go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp** — https://github.com/open-telemetry/opentelemetry-go/blob/exporters/otlp/otlpmetric/otlpmetrichttp/v1.43.0/exporters/otlp/otlpmetric/otlpmetrichttp/LICENSE
- **go.opentelemetry.io/otel/exporters/otlp/otlptrace** — https://github.com/open-telemetry/opentelemetry-go/blob/exporters/otlp/otlptrace/v1.43.0/exporters/otlp/otlptrace/LICENSE
- **go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp** — https://github.com/open-telemetry/opentelemetry-go/blob/exporters/otlp/otlptrace/otlptracehttp/v1.43.0/exporters/otlp/otlptrace/otlptracehttp/LICENSE
- **go.opentelemetry.io/otel/metric** — https://github.com/open-telemetry/opentelemetry-go/blob/metric/v1.43.0/metric/LICENSE
- **go.opentelemetry.io/otel/sdk** — https://github.com/open-telemetry/opentelemetry-go/blob/sdk/v1.43.0/sdk/LICENSE
- **go.opentelemetry.io/otel/sdk/metric** — https://github.com/open-telemetry/opentelemetry-go/blob/sdk/metric/v1.43.0/sdk/metric/LICENSE
- **go.opentelemetry.io/otel/trace** — https://github.com/open-telemetry/opentelemetry-go/blob/trace/v1.43.0/trace/LICENSE
- **go.opentelemetry.io/proto/otlp** — https://github.com/open-telemetry/opentelemetry-proto-go/blob/otlp/v1.10.0/otlp/LICENSE
- **google.golang.org/genproto/googleapis** — https://github.com/googleapis/go-genproto/blob/8636f8732409/LICENSE
- **google.golang.org/genproto/googleapis/api** — https://github.com/googleapis/go-genproto/blob/9d38bb4040a9/googleapis/api/LICENSE
- **google.golang.org/genproto/googleapis/rpc** — https://github.com/googleapis/go-genproto/blob/9d38bb4040a9/googleapis/rpc/LICENSE
- **google.golang.org/grpc** — https://github.com/grpc/grpc-go/blob/v1.80.0/LICENSE
- **gopkg.in/yaml.v2** — https://github.com/go-yaml/yaml/blob/v2.4.0/LICENSE

## MIT (49)

- **github.com/cenkalti/backoff/v4** — https://github.com/cenkalti/backoff/blob/v4.3.0/LICENSE
- **github.com/cenkalti/backoff/v5** — https://github.com/cenkalti/backoff/blob/v5.0.3/LICENSE
- **github.com/cespare/xxhash/v2** — https://github.com/cespare/xxhash/blob/v2.3.0/LICENSE.txt
- **github.com/clipperhouse/stringish** — https://github.com/clipperhouse/stringish/blob/v0.1.1/LICENSE
- **github.com/clipperhouse/uax29/v2** — https://github.com/clipperhouse/uax29/blob/v2.4.0/LICENSE
- **github.com/cpuguy83/dockercfg** — https://github.com/cpuguy83/dockercfg/blob/v0.3.2/LICENSE
- **github.com/deckarep/golang-set/v2** — https://github.com/deckarep/golang-set/blob/v2.8.0/LICENSE
- **github.com/emicklei/dot** — https://github.com/emicklei/dot/blob/v1.10.0/LICENSE
- **github.com/fatih/color** — https://github.com/fatih/color/blob/v1.18.0/LICENSE.md
- **github.com/felixge/httpsnoop** — https://github.com/felixge/httpsnoop/blob/v1.0.4/LICENSE.txt
- **github.com/ferranbt/fastssz** — https://github.com/ferranbt/fastssz/blob/v1.0.0/LICENSE
- **github.com/gabriel-vasile/mimetype** — https://github.com/gabriel-vasile/mimetype/blob/v1.4.12/LICENSE
- **github.com/getsentry/sentry-go** — https://github.com/getsentry/sentry-go/blob/v0.42.0/LICENSE
- **github.com/go-playground/locales** — https://github.com/go-playground/locales/blob/v0.14.1/LICENSE
- **github.com/go-playground/universal-translator** — https://github.com/go-playground/universal-translator/blob/v0.18.1/LICENSE
- **github.com/go-playground/validator/v10** — https://github.com/go-playground/validator/blob/v10.30.1/LICENSE
- **github.com/go-viper/mapstructure/v2** — https://github.com/go-viper/mapstructure/blob/v2.5.0/LICENSE
- **github.com/golang-migrate/migrate/v4** — https://github.com/golang-migrate/migrate/blob/v4.19.1/LICENSE
- **github.com/holiman/bloomfilter/v2** — https://github.com/holiman/bloomfilter
- **github.com/jackc/pgerrcode** — https://github.com/jackc/pgerrcode/blob/469b46aa5efa/LICENSE
- **github.com/jackc/pgpassfile** — https://github.com/jackc/pgpassfile/blob/v1.0.0/LICENSE
- **github.com/jackc/pgservicefile** — https://github.com/jackc/pgservicefile/blob/5a60cdf6a761/LICENSE
- **github.com/jackc/pgx/v5** — https://github.com/jackc/pgx/blob/v5.9.0/LICENSE
- **github.com/jackc/puddle/v2** — https://github.com/jackc/puddle/blob/v2.2.2/LICENSE
- **github.com/klauspost/compress/zstd/internal/xxhash** — https://github.com/klauspost/compress/blob/v1.18.3/zstd/internal/xxhash/LICENSE.txt
- **github.com/klauspost/cpuid/v2** — https://github.com/klauspost/cpuid/blob/v2.3.0/LICENSE
- **github.com/kr/pretty** — https://github.com/kr/pretty/blob/v0.3.1/License
- **github.com/kr/text** — https://github.com/kr/text/blob/v0.2.0/License
- **github.com/leodido/go-urn** — https://github.com/leodido/go-urn/blob/v1.4.0/LICENSE
- **github.com/mattn/go-colorable** — https://github.com/mattn/go-colorable/blob/v0.1.14/LICENSE
- **github.com/mattn/go-isatty** — https://github.com/mattn/go-isatty/blob/v0.0.20/LICENSE
- **github.com/mattn/go-runewidth** — https://github.com/mattn/go-runewidth/blob/v0.0.19/LICENSE
- **github.com/mcuadros/go-defaults** — https://github.com/mcuadros/go-defaults/blob/v1.2.0/LICENSE
- **github.com/mitchellh/mapstructure** — https://github.com/mitchellh/mapstructure/blob/v1.5.0/LICENSE
- **github.com/morikuni/aec** — https://github.com/morikuni/aec/blob/v1.1.0/LICENSE
- **github.com/olekukonko/tablewriter** — https://github.com/olekukonko/tablewriter/blob/v0.0.5/LICENSE.md
- **github.com/pelletier/go-toml/v2** — https://github.com/pelletier/go-toml/blob/v2.2.4/LICENSE
- **github.com/sagikazarmark/locafero** — https://github.com/sagikazarmark/locafero/blob/v0.12.0/LICENSE
- **github.com/sirupsen/logrus** — https://github.com/sirupsen/logrus/blob/v1.9.4/LICENSE
- **github.com/spf13/cast** — https://github.com/spf13/cast/blob/v1.10.0/LICENSE
- **github.com/spf13/viper** — https://github.com/spf13/viper/blob/v1.21.0/LICENSE
- **github.com/stretchr/testify** — https://github.com/stretchr/testify/blob/v1.11.1/LICENSE
- **github.com/subosito/gotenv** — https://github.com/subosito/gotenv/blob/v1.6.0/LICENSE
- **github.com/testcontainers/testcontainers-go** — https://github.com/testcontainers/testcontainers-go/blob/v0.40.0/LICENSE
- **github.com/testcontainers/testcontainers-go/modules/nats** — https://github.com/testcontainers/testcontainers-go/blob/modules/nats/v0.40.0/modules/nats/LICENSE
- **github.com/testcontainers/testcontainers-go/modules/postgres** — https://github.com/testcontainers/testcontainers-go/blob/modules/postgres/v0.40.0/modules/postgres/LICENSE
- **github.com/VictoriaMetrics/fastcache** — https://github.com/VictoriaMetrics/fastcache/blob/v1.13.0/LICENSE
- **go.yaml.in/yaml/v3** — https://github.com/yaml/go-yaml/blob/v3.0.4/LICENSE
- **gopkg.in/yaml.v3** — https://github.com/go-yaml/yaml/blob/v3.0.1/LICENSE

## BSD-3-Clause (33)

- **dario.cat/mergo** — https://github.com/imdario/mergo/blob/v1.0.2/LICENSE
- **github.com/aws/aws-sdk-go-v2/internal/sync/singleflight** — https://github.com/aws/aws-sdk-go-v2/blob/v1.41.6/internal/sync/singleflight/LICENSE
- **github.com/aws/smithy-go/internal/sync/singleflight** — https://github.com/aws/smithy-go/blob/v1.25.0/internal/sync/singleflight/LICENSE
- **github.com/bits-and-blooms/bitset** — https://github.com/bits-and-blooms/bitset/blob/v1.24.4/LICENSE
- **github.com/ethereum/go-bigmodexpfix/src** — https://github.com/ethereum/go-bigmodexpfix/blob/f9e208c548ab/LICENSE
- **github.com/ethereum/go-ethereum/crypto/bn256** — https://github.com/ethereum/go-ethereum/blob/v1.16.9/crypto/bn256/LICENSE
- **github.com/ethereum/go-ethereum/crypto/secp256k1** — https://github.com/ethereum/go-ethereum/blob/v1.16.9/crypto/secp256k1/LICENSE
- **github.com/fsnotify/fsnotify** — https://github.com/fsnotify/fsnotify/blob/v1.9.0/LICENSE
- **github.com/gofrs/flock** — https://github.com/gofrs/flock/blob/v0.13.0/LICENSE
- **github.com/gogo/protobuf** — https://github.com/gogo/protobuf/blob/v1.3.2/LICENSE
- **github.com/golang/snappy** — https://github.com/golang/snappy/blob/v1.0.0/LICENSE
- **github.com/google/uuid** — https://github.com/google/uuid/blob/v1.6.0/LICENSE
- **github.com/googleapis/gax-go/v2** — https://github.com/googleapis/gax-go/blob/v2.16.0/v2/LICENSE
- **github.com/grpc-ecosystem/grpc-gateway/v2** — https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.28.0/LICENSE
- **github.com/holiman/uint256** — https://github.com/holiman/uint256/blob/v1.3.2/COPYING
- **github.com/klauspost/compress/internal/snapref** — https://github.com/klauspost/compress/blob/v1.18.3/internal/snapref/LICENSE
- **github.com/pmezard/go-difflib/difflib** — https://github.com/pmezard/go-difflib/blob/5d4384ee4fb2/LICENSE
- **github.com/rogpeppe/go-internal/fmtsort** — https://github.com/rogpeppe/go-internal/blob/v1.14.1/LICENSE
- **github.com/shirou/gopsutil** — https://github.com/shirou/gopsutil/blob/v3.21.11/LICENSE
- **github.com/shirou/gopsutil/v4** — https://github.com/shirou/gopsutil/blob/v4.25.12/LICENSE
- **github.com/spf13/pflag** — https://github.com/spf13/pflag/blob/v1.0.10/LICENSE
- **github.com/tklauser/go-sysconf** — https://github.com/tklauser/go-sysconf/blob/v0.3.16/LICENSE
- **golang.org/x/crypto** — https://cs.opensource.google/go/x/crypto/+/v0.49.0:LICENSE
- **golang.org/x/exp/maps** — https://cs.opensource.google/go/x/exp/+/716be562:LICENSE
- **golang.org/x/net** — https://cs.opensource.google/go/x/net/+/v0.52.0:LICENSE
- **golang.org/x/oauth2** — https://cs.opensource.google/go/x/oauth2/+/v0.35.0:LICENSE
- **golang.org/x/sync** — https://cs.opensource.google/go/x/sync/+/v0.20.0:LICENSE
- **golang.org/x/sys** — https://cs.opensource.google/go/x/sys/+/v0.42.0:LICENSE
- **golang.org/x/text** — https://cs.opensource.google/go/x/text/+/v0.35.0:LICENSE
- **golang.org/x/time/rate** — https://cs.opensource.google/go/x/time/+/v0.14.0:LICENSE
- **google.golang.org/api** — https://github.com/googleapis/google-api-go-client/blob/v0.265.0/LICENSE
- **google.golang.org/api/internal/third_party/uritemplates** — https://github.com/googleapis/google-api-go-client/blob/v0.265.0/internal/third_party/uritemplates/LICENSE
- **google.golang.org/protobuf** — https://github.com/protocolbuffers/protobuf-go/blob/v1.36.11/LICENSE

## BSD-2-Clause (3)

- **github.com/gorilla/websocket** — https://github.com/gorilla/websocket/blob/v1.5.3/LICENSE
- **github.com/magiconair/properties** — https://github.com/magiconair/properties/blob/v1.8.10/LICENSE.md
- **github.com/pkg/errors** — https://github.com/pkg/errors/blob/v0.9.1/LICENSE

## BSD-2-Clause-FreeBSD (1)

- **github.com/ethereum/go-ethereum/metrics** — https://github.com/ethereum/go-ethereum/blob/v1.16.9/metrics/LICENSE

## ISC (1)

- **github.com/davecgh/go-spew/spew** — https://github.com/davecgh/go-spew/blob/d8f796af33cc/LICENSE

## Unlicense (1)

- **github.com/ethereum/go-verkle** — https://github.com/ethereum/go-verkle/blob/v0.2.2/LICENSE

## Public Domain (CC0) (1)

- **github.com/dchest/blake512** — https://github.com/dchest/blake512

