FROM bitnami/kubectl as kubectl

FROM debian:trixie-slim

RUN apt update
RUN apt install gettext-base -y

RUN apt-get autoremove --purge && \
    apt-get update && apt-get upgrade -y && \
    apt-get clean && rm -rf /var/lib/apt/lists /var/cache/apt/archives

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/

USER 1001
ENTRYPOINT [ "kubectl" ]
CMD [ "--help" ]