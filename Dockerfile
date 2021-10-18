ARG GO_VERSION="1.16"

FROM golang:$GO_VERSION

ARG TARGETPLATFORM
ARG PKG_NAME=ksops

# Match Argo CD's build
ENV GO111MODULE=on

# Define kustomize config location
ENV XDG_CONFIG_HOME=$HOME/.config

# Export templated Go env variables
RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) && \
    export GOARM=$(echo ${TARGETPLATFORM} | cut -d / -f3 | cut -c2-)

WORKDIR /go/src/github.com/argyle-engineering/ksops

ADD . .

# Perform the build
RUN rm ksops|| true
RUN rm -rf ${XDG_CONFIG_HOME}/kustomize/plugin/argyle.com/v1/ || true
RUN rm -rf ${HOME}/sigs.k8s.io/kustomize/plugin/argyle.com/v1/ || true
RUN go install
RUN go build -o ksops

CMD ["kustomize", "version"]
