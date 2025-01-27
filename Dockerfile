FROM scratch
COPY dsrepo /dsrepo
LABEL org.opencontainers.image.source https://github.com/davidjspooner/dsrepo
ENTRYPOINT ["/dsrepo"]
