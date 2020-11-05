# drone-rancher

[![Build Status](https://cloud.drone.io/api/badges/josmo/drone-rancher/status.svg)](https://cloud.drone.io/josmo/drone-rancher)
[![Go Doc](https://godoc.org/github.com/josmo/drone-rancher?status.svg)](http://godoc.org/github.com/josmo/drone-rancher)
[![Go Report](https://goreportcard.com/badge/github.com/josmo/drone-rancher)](https://goreportcard.com/report/github.com/josmo/drone-rancher)
[![](https://images.microbadger.com/badges/image/pelotech/drone-rancher.svg)](https://microbadger.com/images/pelotech/drone-rancher "Get your own image badge on microbadger.com")

Drone plugin to deploy or update a project on Rancher 1.x only. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `drone cli`:

```
drone exec
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
  pelotech/drone-rancher 
```

### Contribution

This repo is setup in a way that if you enable a personal drone server to build your fork it will
 build and publish your image (makes it easier to test PRs and use the image till the contributions get merged)
 
* Build local ```DRONE_REPO_OWNER=josmo DRONE_REPO_NAME=drone-rancher drone exec```
* on your server just make sure you have DOCKER_USERNAME, DOCKER_PASSWORD, and PLUGIN_REPO set as secrets
