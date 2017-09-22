FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch

ENV GODEBUG=netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/croudsupport/drone-rancher.git"
LABEL org.label-schema.name="Drone Rancher"
LABEL org.label-schema.vendor="Croud"

ADD release/linux/amd64/drone-rancher /bin/

ENTRYPOINT ["/bin/drone-rancher"]
