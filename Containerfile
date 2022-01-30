ARG VARIANT=stable-slim
FROM docker.io/library/debian:$VARIANT

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get clean && rm -fr /var/lib/apt/lists/*
