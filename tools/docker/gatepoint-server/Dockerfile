FROM --platform=$BUILDPLATFORM golang:1.22 as builder
ARG TARGETARCH

RUN apt-get update && apt-get install --no-install-recommends -y \
    protobuf-compiler && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \

RUN buf --version

WORKDIR /workspace
COPY ./ .
RUN make api_install_dep
RUN make api_gen
RUN make build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/gatepoint-server ./
COPY --from=builder /workspace/config/ ./config/

USER 65532:65532

CMD ["./gatepoint-server", "server"]

EXPOSE 8081 9091
