FROM alpine:3.13

COPY rabdis /bin/rabdis

RUN addgroup -g 1000 -S rabdis && \
    adduser -u 1000 -S rabdis -G rabdis && \
    chown rabdis:rabdis /bin/rabdis

USER rabdis:rabdis

EXPOSE 9090

ENTRYPOINT ["/bin/rabdis"]
