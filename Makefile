.EXPORT_ALL_VARIABLES:

help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

GIT_TAG?=$(shell git describe --tags --abbrev=0 2>/dev/null || printf 'dev')
GIT_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD | tr /_ - | cut -c1-16 )
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
VERSION?=${GIT_TAG}-${GIT_BRANCH}-${GIT_COMMIT}

VERSION:
	@printf ${VERSION}

DOCKER_REGISTRY?=10.10.0.116:5000
DOCKER_REPO?=galactica/faucet
IMG_TAG?=${DOCKER_REGISTRY}/${DOCKER_REPO}:${VERSION}
DOCKER_BUILDKIT=1

IMG_TAG:
	@printf ${IMG_TAG}

docker-build:
	docker build -t ${IMG_TAG} . 

docker-push:
	docker push ${IMG_TAG}