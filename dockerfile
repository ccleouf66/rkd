FROM ubuntu:20.04 AS builder

ENV ENV="/root/.bashrc" \
    TZ=Europe \
    EDITOR=vi \
    LANG=en_US.UTF-8

WORKDIR /go/src/rkd

COPY . /go/src/rkd/

ADD https://golang.org/dl/go1.15.2.linux-amd64.tar.gz /tmp

RUN    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
       apt-get update && \
       apt-get install -y sudo git build-essential make libdevmapper-dev libgpgme-dev libostree-dev curl libassuan-dev libbtrfs-dev && \
       tar -C /usr/local -xzf /tmp/go1.15.2.linux-amd64.tar.gz && \
       rm /tmp/go1.15.2.linux-amd64.tar.gz && \
       ln -s /usr/local/go/bin/* /usr/local/bin/ && \
       go get && \
       GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/src/rkd/rkd-linux-amd64

FROM debian:stretch-slim
WORKDIR /go/bin/
COPY --from='builder' /go/src/rkd/rkd-linux-amd64 /go/bin/rkd-linux-amd64
COPY --from='builder' /go/src/rkd/policy.json /go/bin/
COPY --from='builder' /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
run apt update && \
    apt install -y libgpgme-dev libdevmapper-dev && \
    rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/go/bin/rkd-linux-amd64"]