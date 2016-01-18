# drone-rancher

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-rancher/status.svg)](http://beta.drone.io/drone-plugins/drone-rancher)
[![](https://badge.imagelayers.io/plugins/drone-rancher:latest.svg)](https://imagelayers.io/?images=plugins/drone-rancher:latest 'Get your own badge on imagelayers.io')

Drone plugin for deploying to Rancher

## Usage

```
./drone-rancher <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
    }
}
EOF
```

## Docker

Build the Docker container using `make`:

```
make deps build docker
```

### Example

```sh
docker run -i plugins/drone-rancher <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
    }
}
EOF
```
