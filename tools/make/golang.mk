VERSION_PACKAGE := github.com/gatepoint/gatepoint/internal/version

GO_LDFLAGS += -X $(VERSION_PACKAGE).binVersion=$(shell cat VERSION) \
	-X $(VERSION_PACKAGE).gitCommitID=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).buildDate=$(shell TZ=Asia/Shanghai date +%FT%T%z) \
	-X $(VERSION_PACKAGE).gitBranch=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: build-multiarch
build-multiarch: go.build.multiarch

.PHONY: go.build.multiarch
go.build.multiarch: $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))

.PHONY: go.build
go.build: $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

.PHONY: go.build.%
go.build.%:
	@$(LOG_TARGET)
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@$(call log, "Building binary $(COMMAND) with commit $(REV) for $(OS) $(ARCH)")
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(COMMAND) -ldflags "$(GO_LDFLAGS)" $(ROOT_PACKAGE)/cmd/$(COMMAND)
