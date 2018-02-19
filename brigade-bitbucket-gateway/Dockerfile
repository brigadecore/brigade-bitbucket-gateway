FROM alpine:3.7

RUN apk update && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates

COPY rootfs/brigade-bitbucket-gateway /usr/bin/brigade-bitbucket-gateway
COPY rootfs/gitssh.sh /gitssh.sh

ENV GIT_SSH=/gitssh.sh

CMD /usr/bin/brigade-bitbucket-gateway
