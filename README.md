# drone-rancher

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-rancher/status.svg)](http://beta.drone.io/drone-plugins/drone-rancher)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-rancher/coverage.svg)](https://aircover.co/drone-plugins/drone-rancher)
[![](https://badge.imagelayers.io/plugins/drone-rancher:latest.svg)](https://imagelayers.io/?images=plugins/drone-rancher:latest 'Get your own badge on imagelayers.io')

Drone plugin to deploy or update a project on Rancher

## Binary

Build the binary using `make`:

```
make deps build
```

### Example

```sh
./drone-rancher <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "url": "https://example.rancher.com",
        "access_key": "1234567abcdefg",
        "secret_key": "abcdefg1234567",
        "service": "drone/drone",
        "docker_image": "drone/drone:latest"
    }
}
EOF
```

## Docker

Build the container using `make`:

```
make deps docker
```

### Example

```sh
docker run -i plugins/drone-rancher <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "url": "https://example.rancher.com",
        "access_key": "1234567abcdefg",
        "secret_key": "abcdefg1234567",
        "service": "drone/drone",
        "docker_image": "drone/drone:latest"
    }
}
EOF
```
