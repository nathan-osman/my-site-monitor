PKG = github.com/nathan-osman/my-site-monitor
CMD = msm

CWD = $(shell pwd)
UID = $(shell id -u)
GID = $(shell id -g)

# Find all Go source files (excluding the cache path)
SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')

# Find all source files that comprise the UI
UIFILES = $(shell find ui/public ui/src) ui/package-lock.json

all: dist/${CMD}

# Build the standalone executable
dist/${CMD}: ${SOURCES} server/ab0x.go | cache/lib cache/src/${PKG} dist
	@docker run \
	    --rm \
	    -e HOME=/tmp \
	    -e CGO_ENABLED=0 \
	    -e GIT_COMMITTER_NAME=a \
	    -e GIT_COMMITTER_EMAIL=b \
	    -u ${UID}:${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    golang:latest \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}

# Create a Go source file with the static files
server/ab0x.go: dist/fileb0x b0x.yaml .dep-static
	@dist/fileb0x b0x.yaml

# Create the fileb0x executable needed for embedding files
dist/fileb0x: | cache/lib cache/src/${PKG} dist
	@docker run \
	    --rm \
	    -e HOME=/tmp \
	    -e CGO_ENABLED=0 \
	    -e GIT_COMMITTER_NAME=a \
	    -e GIT_COMMITTER_EMAIL=b \
	    -u ${UID}:${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    golang:latest \
	    go get -pkgdir /go/lib github.com/UnnoTed/fileb0x

# Build the UI
.dep-static: ${UIFILES} .dep-node_modules
	@docker run \
	    --rm \
	    -e HOME=/tmp \
	    -u ${UID}:${GID} \
	    -v ${CWD}/ui:/usr/src/ui \
	    -w /usr/src/ui \
	    node:latest \
	    npm run build
	@touch .dep-static

# Fetch NPM packages for building the UI
.dep-node_modules: ui/package.json
	@docker run \
	    --rm \
	    -e HOME=/tmp \
	    -u ${UID}:${GID} \
	    -v ${CWD}/ui:/usr/src/ui \
	    -w /usr/src/ui \
	    node:latest \
	    npm install
	@touch .dep-node_modules

cache/lib:
	@mkdir -p cache/lib

cache/src/${PKG}:
	@mkdir -p cache/src/${PKG}

dist:
	@mkdir dist

clean:
	@rm -rf .dep-* cache dist server/ab0x.go ui/{dist,node_modules}

.PHONY: clean