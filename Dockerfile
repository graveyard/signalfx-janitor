FROM debian

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update && apt-get install -y ca-certificates curl

RUN curl -L https://github.com/Clever/batchcli/releases/download/0.0.12/batchcli-v0.0.12-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components 1

COPY signalfx-janitor /bin/signalfx-janitor

ENTRYPOINT ["/usr/local/bin/batchcli"]
CMD ["--cmd", "/bin/signalfx-janitor"]
