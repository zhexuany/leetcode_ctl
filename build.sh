go build -v -v -ldflags="-X main.version=$(git describe --always --tags --abbrev=0) -X main.commit=$(git rev-parse HEAD) -X main.branch=$(git rev-parse --abbrev-ref HEAD)"
# git rev-parse HEAD
#git log --pretty=format:'%h' -n 1

