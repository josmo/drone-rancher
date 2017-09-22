
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0


go build -v -ldflags -a -o release/linux/amd64/drone-rancher
docker build . --tag croudtech/drone-rancher:latest
