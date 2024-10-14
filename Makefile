.PHONY: build run lint format
build:
	CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o ./bin/gatepoint-server ./cmd/gatepoint-server/

run:
	bin/gatepoint-server server

lint:
	make proto_lint
	golint ./...

format: proto_format

.PHONY: proto_lint proto_format
proto_lint:
	buf lint

proto_format:
	buf format api -d --exit-code



.PHONY: server-image
server-image:
	docker buildx build -f tools/docker/gatepoint-server/Dockerfile -t release.daocloud.io/skoala/gatepoint:v0.12 . --platform linux/amd64,linux/arm64 --push

.PHONY: _run
_run:
	@$(MAKE) --warn-undefined-variables -f tools/make/common.mk $(MAKECMDGOALS)

$(if $(MAKECMDGOALS),$(MAKECMDGOALS): %: _run)
