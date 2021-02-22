ARG ALPINE_VERSION=3.13
FROM alpine:${ALPINE_VERSION}

COPY rabdis /bin/rabdis

RUN addgroup -g 1000 -S rabdis && \
    adduser -u 1000 -S rabdis -G rabdis && \
    mkdir -p /rabdis && \
    chown rabdis:rabdis /bin/rabdis

USER rabdis:rabdis

EXPOSE 9090

ENTRYPOINT ["/bin/rabdis"]
