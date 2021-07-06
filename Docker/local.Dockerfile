FROM autodesk-docker-build-images.***REMOVED***/hardened-build/golang-1.16:latest as BASE
LABEL maintainer="shore@autodesk.com"

ARG JT_VERSION="0.0.6"
ARG JT_FILE_NAME="jsonnet-test_${JT_VERSION}_Linux_x86_64.tar.gz"

ARG JB_VERSION="v0.4.0"
ARG JB_FILE_NAME="jb-linux-amd64"

ARG JSONNET_VERSION='0.17.0'
ARG JSONNET_FILE_NAME="go-jsonnet_${JSONNET_VERSION}_Linux_x86_64.tar.gz"

ARG SPIN_CLI_VERSION="1.22.0"
ARG SPIN_CLI_FILE_NAME="spin"


WORKDIR /tmp/build

RUN echo "Installing Jsonnet Bundler (${JB_VERSION}), jsonnet-test (v${JT_VERSION}), (go) Jsonnet (v${JSONNET_VERSION}), spin-cli (v${SPIN_CLI_VERSION})" && \
    # Jsonnet-Bundler
    wget -q https://github.com/jsonnet-bundler/jsonnet-bundler/releases/download/${JB_VERSION}/${JB_FILE_NAME} && \
    chmod +x ${JB_FILE_NAME} && \
    mv ${JB_FILE_NAME} jb && \ 
    # Jsonnet-test
    wget -q https://***REMOVED***/***REMOVED***/team-shore-generic/jsonnet-test/${JT_VERSION}/linux/amd64/${JT_FILE_NAME} && \
    tar -xzvf ${JT_FILE_NAME} && \
    chmod +x jsonnet-test && \
    mv jsonnet-test jt && \
    # Jsonnet
    wget -q https://github.com/google/go-jsonnet/releases/download/v${JSONNET_VERSION}/${JSONNET_FILE_NAME} && \
    tar -xzvf ${JSONNET_FILE_NAME} && \
    chmod +x jsonnet && \
    chmod +x jsonnetfmt && \
    # spin-cli
    wget -q https://storage.googleapis.com/spinnaker-artifacts/spin/${SPIN_CLI_VERSION}/linux/amd64/${SPIN_CLI_FILE_NAME} && \
    chmod +x ${SPIN_CLI_FILE_NAME}


# Final Container
FROM autodesk-docker-build-images.***REMOVED***/hardened-build/golang-1.16:latest

WORKDIR /shore

RUN apk add git make --no-cache

COPY --from=BASE /tmp/build/jb /usr/local/bin/jb
COPY --from=BASE /tmp/build/jt /usr/local/bin/jt
COPY --from=BASE /tmp/build/jsonnet /usr/local/bin/jsonnet
COPY --from=BASE /tmp/build/jsonnetfmt /usr/local/bin/jsonnetfmt
COPY --from=BASE /tmp/build/spin /usr/local/bin/spin

# Copy over the source, and install from it.
COPY / .
RUN make setup && \
    go build -o shore cmd/shore/shore.go
