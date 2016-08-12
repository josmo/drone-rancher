# drone-rancher

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-rancher/status.svg)](http://beta.drone.io/drone-plugins/drone-rancher)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-rancher/coverage.svg)](https://aircover.co/drone-plugins/drone-rancher)
[![](https://badge.imagelayers.io/plugins/drone-rancher:latest.svg)](https://imagelayers.io/?images=plugins/drone-rancher:latest 'Get your own badge on imagelayers.io')

Drone plugin to deploy or update a project on Rancher. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `make`:

```
make deps build
```

## Image

Build the docker images using `make`:

```
make docker
```

## Usage:

Build and publish from your current working directory:

```
docker run --rm \
  -e PLUGIN_URL=rancher.youdomain.com \
  -e PLUGIN_ACCESS_KEY=4ccesskey \
  -e PLUGIN_SECRET_KEY=secretkey \
  -e PLUGIN_SERVICE=stack/service \
  -e PLUGIN_IMAGE=yourregistry.com/image:tag \
  -e PLUGIN_START_FIRST=true \
  -e PLUGIN_CONFIRM=true \
  -e PLUGIN_TIMEOUT=20 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  --privileged \
  plugins/drone-rancher
```
