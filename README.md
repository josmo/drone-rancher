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


## Docker

Build the container using `make`:

```
make deps docker
```

### Example

## Usage

Build and deploy from your current working directory:

```
docker run --rm                          \
  -e PLUGIN_URL=<source>                 \
  -e PLUGIN_ACCESS_KEY=<key>     \
  -e PLUGIN_SECRET_KEY=<secret>  \
  -e PLUGIN_SERVICE=<service>            \  
  -e PLUGIN_DOCKER_IMAGE=<image>         \
  -v $(pwd):$(pwd)                       \
  -w $(pwd)                              \
  plugins/rancher 
```
