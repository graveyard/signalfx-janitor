FROM debian

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update && apt-get install -y ca-certificates
COPY build/batchcli-v0.0.12-linux-amd64/batchcli /usr/local/bin/batchcli
RUN chmod +x /usr/local/bin/batchcli

COPY signalfx-janitor /bin/signalfx-janitor

ENTRYPOINT ["/usr/local/bin/batchcli"]
CMD ["--cmd", "/bin/signalfx-janitor"]
