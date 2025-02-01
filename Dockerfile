ARG DEBIAN_IMAGE=golang:1.23.5-bookworm
ARG BASE=gcr.io/distroless/static-debian11:nonroot
ARG COREDNS_VERSION=v1.12.0
ARG COREDNS_LIBVIRT_VERSION=v0.1.2
FROM ${DEBIAN_IMAGE} AS build
SHELL [ "/bin/sh", "-ec" ]

RUN export DEBCONF_NONINTERACTIVE_SEEN=true \
           DEBIAN_FRONTEND=noninteractive \
           DEBIAN_PRIORITY=critical \
           TERM=linux ; \
    apt-get -qq update ; \
    apt-get -yyqq upgrade ; \
    apt-get -yyqq install ca-certificates libcap2-bin libvirt-dev; \
    apt-get clean
COPY ./docker /docker
WORKDIR /docker
RUN export COREDNS_URI="github.com/coredns/coredns@${COREDNS_VERSION}"
RUN export COREDNS_LIBVIRT_URI="github.com/ironpinguin/coredns-libvirt@${COREDNS_LIBVIRT_VERSION}"
RUN go get ${COREDNS_URI} ${COREDNS_LIBVIRT_URI} ; \
    go build -o /coredns \
        -ldflags "-s -w -X github.com/coredns/coredns/coremain.Version=${COREDNS_VERSION}" /docker/main.go
RUN setcap cap_net_bind_service=+ep /coredns
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
