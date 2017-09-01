# drone-rancher

[![Build Status](https://drone.seattleslow.com/api/badges/josmo/drone-rancher/status.svg)](https://drone.seattleslow.com/josmo/drone-rancher)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)
[![Go Doc](https://godoc.org/github.com/josmo/drone-rancher?status.svg)](http://godoc.org/github.com/josmo/drone-rancher)
[![Go Report](https://goreportcard.com/badge/github.com/josmo/drone-rancher)](https://goreportcard.com/report/github.com/josmo/drone-rancher)
[![](https://images.microbadger.com/badges/image/peloton/drone-rancher.svg)](https://microbadger.com/images/peloton/drone-rancher "Get your own image badge on microbadger.com")

Drone plugin to deploy or update a project on Rancher. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

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
  peloton/drone-rancher 
```
