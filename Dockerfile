# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-rancher .

FROM gliderlabs/alpine:3.2
RUN apk add --update \
  ca-certificates
ADD drone-rancher /bin/
ENTRYPOINT ["/bin/drone-rancher"]

