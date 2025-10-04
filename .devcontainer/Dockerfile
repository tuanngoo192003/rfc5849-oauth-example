# Golang on Alpine
FROM golang:1.24-alpine

RUN apk update && apk add --no-cache \
    git \
    curl \
    wget \
    unzip \
    make \
    vim \
    net-tools \
    iputils \
    jq \
    bash \
    tar \
    && rm -rf /var/cache/apk/*

WORKDIR /work
