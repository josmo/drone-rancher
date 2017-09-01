FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch

ENV GODEBUG=netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/josmo/drone-rancher.git"
LABEL org.label-schema.name="Drone Rancher"
LABEL org.label-schema.vendor="Josmo"
LABEL org.label-schema.schema-version="1.0"

ADD release/linux/arm/drone-rancher /bin/
ENTRYPOINT ["/bin/drone-rancher"]