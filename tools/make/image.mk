include tools/make/env.mk

DOCKER := docker

IMAGES_DIR := $(wildcard ${ROOT_DIR}tools/docker/*)

IMAGES ?= gatepoint-server
IMAGE_PLATFORMS ?= linux_amd64 linux_arm64


# Convert to linux/arm64,linux/amd64
$(eval BUILDX_PLATFORMS = $(shell echo "${IMAGE_PLATFORMS}" | sed "s# #,#;s#_#/#g"))

.PHONY: image.build
image.build: $(addprefix image.build., $(IMAGES))

.PHONY: image.build.%
image.build.%: go.build.linux_$(GOARCH).%
	@$(LOG_TARGET)
	$(eval COMMAND := $(word 1,$(subst ., ,$*)))
	$(eval IMAGES := $(COMMAND))
	@$(call log, "Building image $(IMAGES):$(TAG) in linux/$(GOARCH)")
	$(eval BUILD_SUFFIX := --pull --load -t $(IMAGE):$(TAG) -f $(ROOT_DIR)/tools/docker/$(IMAGES)/Dockerfile bin)
	@$(call log, "Creating image tag $(REGISTRY)/$(IMAGES):$(TAG) in linux/$(GOARCH)")
	$(DOCKER) buildx build --platform linux/$(GOARCH) $(BUILD_SUFFIX)

.PHONY: image.multiarch.setup
image.multiarch.setup: image.verify image.multiarch.verify image.multiarch.emulate
	@$(LOG_TARGET)
	docker buildx rm $(BUILDX_CONTEXT) || :
	docker buildx create --use --name $(BUILDX_CONTEXT) --platform "${BUILDX_PLATFORMS}"

.PHONY: image.multiarch.emulate $(EMULATE_TARGETS)
image.multiarch.emulate: $(EMULATE_TARGETS)
$(EMULATE_TARGETS): image.multiarch.emulate.%:
	@$(LOG_TARGET)
# Install QEMU emulator, the same emulator as the host will report an error but can safe ignore
	docker run --rm --privileged tonistiigi/binfmt --install linux/$*

.PHONY: image.build.multiarch
image.build.multiarch:
	@$(LOG_TARGET)
	docker buildx build bin -f "$(ROOT_DIR)/tools/docker/$(IMAGES)/Dockerfile" -t "${IMAGE}:${TAG}" --platform "${BUILDX_PLATFORMS}"

.PHONY: image.push.multiarch
image.push.multiarch:
	@$(LOG_TARGET)
	docker buildx build bin -f "$(ROOT_DIR)/tools/docker/$(IMAGES)/Dockerfile" -t "${IMAGE}:${TAG}" --platform "${BUILDX_PLATFORMS}" --sbom=false --provenance=false --push


.PHONY: image
image: ## Build docker images for host platform. See Option PLATFORM and BINS.
image: image.build

.PHONY: image-multiarch
image-multiarch: ## Build docker images for multiple platforms. See Option PLATFORMS and IMAGES.
image-multiarch: image.multiarch.setup go.build.multiarch image.build.multiarch

.PHONY: push
push: ## Push docker images to registry.
push: image.push

.PHONY: push-multiarch
push-multiarch: ## Push docker images for multiple platforms to registry.
push-multiarch: image.multiarch.setup go.build.multiarch image.push.multiarch
