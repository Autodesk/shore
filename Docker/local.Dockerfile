FROM golang:1.19 as BASE
LABEL maintainer="shore@autodesk.com"

ARG JB_VERSION="v0.5.1"
ARG JB_FILE_NAME="jb-linux-amd64"

ARG JSONNET_VERSION='0.19.1'
ARG JSONNET_FILE_NAME="go-jsonnet_${JSONNET_VERSION}_Linux_x86_64.tar.gz"

ARG SPIN_CLI_VERSION="1.27.1"
ARG SPIN_CLI_FILE_NAME="spin"


WORKDIR /tmp/build

RUN echo "Installing Jsonnet Bundler (${JB_VERSION}), (go) Jsonnet (v${JSONNET_VERSION}), spin-cli (v${SPIN_CLI_VERSION})" && \
    # Jsonnet-Bundler
    wget -q https://github.com/jsonnet-bundler/jsonnet-bundler/releases/download/${JB_VERSION}/${JB_FILE_NAME} && \
    chmod +x ${JB_FILE_NAME} && \
    mv ${JB_FILE_NAME} jb && \
    # Jsonnet
    wget -q https://github.com/google/go-jsonnet/releases/download/v${JSONNET_VERSION}/${JSONNET_FILE_NAME} && \
    tar -xzvf ${JSONNET_FILE_NAME} && \
    chmod +x jsonnet && \
    chmod +x jsonnetfmt && \
    # spin-cli
    wget -q https://storage.googleapis.com/spinnaker-artifacts/spin/${SPIN_CLI_VERSION}/linux/amd64/${SPIN_CLI_FILE_NAME} && \
    chmod +x ${SPIN_CLI_FILE_NAME}

# Jsonnet-test
COPY jsonnet/tools/jsonnet-test jsonnet-test

RUN echo "Installing Jsonnet Test (local)" && \
    cd jsonnet-test && \
    make setup && \
    make build && \
    chmod +x jt && \
    mv jt ../.

# Final Container
FROM golang:1.19

WORKDIR /shore

RUN apt-get install git make

COPY --from=BASE /tmp/build/jb /usr/local/bin/jb
COPY --from=BASE /tmp/build/jt /usr/local/bin/jt
COPY --from=BASE /tmp/build/jsonnet /usr/local/bin/jsonnet
COPY --from=BASE /tmp/build/jsonnetfmt /usr/local/bin/jsonnetfmt
COPY --from=BASE /tmp/build/spin /usr/local/bin/spin

# Copy over the source, and install from it.
COPY / .
RUN make setup && \
    go build -o shore cmd/shore/shore.go
