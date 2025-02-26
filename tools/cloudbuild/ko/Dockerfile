# syntax=docker/dockerfile:1

ARG KO_VERSION="0.17.1"
ARG KO_CHECKSUM="4f0b979b59880b3232f47d79c940f2279165aaad15a11d7614e8a2c9e5c78c29"

FROM golang:latest

ADD --checksum=sha256:${KO_CHECKSUM} \
    https://github.com/ko-build/ko/releases/download/v${KO_VERSION}/ko_${KO_VERSION}_Linux_x86_64.tar.gz \
    ko.tar.gz
RUN tar xf -C / ko.tar.gz ko && chmod +x /ko
