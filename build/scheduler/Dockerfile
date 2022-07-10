# Build the scheduler binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/scheduler/main.go main.go
COPY apis/ apis/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a \
    -ldflags="-w -s \
    -X k8s.io/component-base/version.gitVersion=v1.24.0-moirai-scheduler-$(shell date +%Y%m%d) \
    -X k8s.io/component-base/version.gitCommit=$(shell git rev-parse HEAD) \
    -X k8s.io/component-base/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
    "\
    -o scheduler main.go

# Use distroless as minimal base image to package the scheduler binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/scheduler .
USER 65532:65532

ENTRYPOINT ["/scheduler"]
