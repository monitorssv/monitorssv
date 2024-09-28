VERSION ?= latest
DOCKER_REPO := monitorssv
PLATFORMS := amd64 arm64
COMPONENTS := api web

.PHONY: build
build:
	go mod tidy
	go build -o monitorssv cmd/monitorssv/main.go

.PHONY: build-web
build-web:
	cd web && yarn install && npm run build

define docker_build
.PHONY: docker-$(1)-$(2)
docker-$(1)-$(2): $(if $(filter web,$(1)),build-web)
	docker buildx build --platform linux/$(2) -f docker/Dockerfile-$(1) -t $(DOCKER_REPO)/monitorssv-$(1):$(VERSION)-$(2) .
endef

$(foreach component,$(COMPONENTS),\
	$(foreach platform,$(PLATFORMS),\
		$(eval $(call docker_build,$(component),$(platform)))))

.PHONY: docker-all
docker-all: $(foreach component,$(COMPONENTS),\
				$(foreach platform,$(PLATFORMS),\
					docker-$(component)-$(platform)))

define docker_push
.PHONY: push-$(1)-$(2)
push-$(1)-$(2):
	docker push $(DOCKER_REPO)/monitorssv-$(1):$(VERSION)-$(2)
endef

$(foreach component,$(COMPONENTS),\
	$(foreach platform,$(PLATFORMS),\
		$(eval $(call docker_push,$(component),$(platform)))))

.PHONY: push-all
push-all: $(foreach component,$(COMPONENTS),\
			$(foreach platform,$(PLATFORMS),\
				push-$(component)-$(platform)))

$(foreach component,$(COMPONENTS),\
			$(foreach platform,$(PLATFORMS),\
				push-$(component)-$(platform))):

define docker_manifest
.PHONY: manifest-$(1)
manifest-$(1):
	docker manifest create $(DOCKER_REPO)/monitorssv-$(1):$(VERSION) \
		$(foreach platform,$(PLATFORMS),$(DOCKER_REPO)/monitorssv-$(1):$(VERSION)-$(platform))
	docker manifest push $(DOCKER_REPO)/monitorssv-$(1):$(VERSION)
endef

$(foreach component,$(COMPONENTS),\
	$(eval $(call docker_manifest,$(component))))

.PHONY: manifest-all
manifest-all: $(foreach component,$(COMPONENTS),manifest-$(component))

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build               - Build the go application"
	@echo "  build-web           - Build the web application"
	@echo "  docker-all          - Build all docker images"
	@echo "  docker-web-arm64    - Build web arm64 docker image"
	@echo "  docker-api-arm64    - Build api arm64 docker image"
	@echo "  docker-web-amd64    - Build web amd64 docker image"
	@echo "  docker-api-amd64    - Build api amd64 docker image"
	@echo "  push-all            - Push all docker images"
	@echo "  manifest-all        - Create and push multi-arch manifests"
	@echo "  help                - Show this help message"