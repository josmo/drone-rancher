FROM microsoft/nanoserver:latest

ENV GODEBUG=netdns=go

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/josmo/drone-rancher.git"
LABEL org.label-schema.name="Drone Rancher"
LABEL org.label-schema.vendor="Josmo"
LABEL org.label-schema.schema-version="1.0"

ADD release/windows/amd64/drone-rancher /bin/
ENTRYPOINT [ "/bin/drone-rancher" ]