DOCKER_REPOSITORY ?= spotify/bigtable-emulator

DOCKER ?= docker
export DOCKER

DOCKER_TAG = $(shell cat version.txt)
DOCKER_LATEST_TAG = latest$(subst SNAPSHOT,snapshot,$(findstring -SNAPSHOT,$(DOCKER_TAG)))

.PHONY: all image push

all: test

image:
	mkdir -p build
	cp Dockerfile.in build/Dockerfile
	cp bigtable-server.go build/bigtable-server.go
	$(DOCKER) build --no-cache=true -t $(DOCKER_REPOSITORY):$(DOCKER_TAG) build

test: image
	bash -ex runtest $(DOCKER_REPOSITORY):$(DOCKER_TAG)

push: test
	$(DOCKER) push $(DOCKER_REPOSITORY):$(DOCKER_TAG)
	$(DOCKER) tag $(DOCKER_REPOSITORY):$(DOCKER_TAG) $(DOCKER_REPOSITORY):$(DOCKER_LATEST_TAG)
	$(DOCKER) push $(DOCKER_REPOSITORY):$(DOCKER_LATEST_TAG)
