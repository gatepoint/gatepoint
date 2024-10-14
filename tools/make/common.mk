SHELL:=/bin/bash

ROOT_PACKAGE=github.com/gatepoint/gatepoint

RELEASE_VERSION=$(shell cat VERSION)

GIT_COMMIT:=$(shell git rev-parse HEAD)

# Supported Platforms for building multiarch binaries.
PLATFORMS ?= darwin_amd64 darwin_arm64 linux_amd64 linux_arm64

# Set Root Directory Path
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell pwd -P))
endif

# Set Output Directory Path
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/bin
endif

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# List commands in cmd directory for building targets
COMMANDS ?= $(wildcard ${ROOT_DIR}/cmd/*)
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

# REV is the short git sha of latest commit.
REV=$(shell git rev-parse --short HEAD)

# Log the running target
LOG_TARGET = echo -e "\033[0;32m===========> Running $@ ... \033[0m"

define log
echo -e "\033[36m===========>$1\033[0m"
endef

define errorlog
echo -e "\033[0;31m===========>$1\033[0m"
endef

include tools/make/lint.mk
include tools/make/api.mk
include tools/make/golang.mk
include tools/make/gen.mk
