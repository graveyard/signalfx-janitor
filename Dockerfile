FROM gliderlabs/alpine:3.3
ENTRYPOINT ["/bin/signalfx-janitor"]

RUN apk-install ca-certificates

COPY signalfx-janitor /bin/signalfx-janitor
