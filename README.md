# drone-rancher

[![Build Status](https://drone.seattleslow.com/api/badges/josmo/drone-rancher/status.svg)](https://drone.seattleslow.com/josmo/drone-rancher)

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
  peloton/drone-rancher 
```
