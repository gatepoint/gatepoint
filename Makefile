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

.PHONY: api_gen api_install_dep api_clean
api_install_dep:
	# go env -w GOPROXY=https://goproxy.cn,direct
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/jstemmer/go-junit-report@latest
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest
	go install github.com/rakyll/statik@latest
	go install istio.io/tools/cmd/protoc-gen-golang-jsonshim@latest
	go install istio.io/tools/cmd/protoc-gen-golang-deepcopy@latest

api_gen:
	buf generate
	#protoc -I . -I third_party \
#		--go_out=paths=source_relative:. \
#		--go-grpc_out=paths=source_relative:. \
#		--grpc-gateway_out=paths=source_relative:. \
#		--openapiv2_out=logtostderr=true:. \
#		--openapiv2_opt allow_merge=true \
#		--openapiv2_opt output_format=json \
#		--openapiv2_opt merge_file_name="gatepoint." \
#		api/gatepoint/v1/gatepoint.proto api/general/v1/demo.proto
	cp api/gen/gatepoint.swagger.json swagger-ui/gatepoint.swagger.json

api_clean:
	#rm -rf api/gen
	rm -f api/*/*/*.pb.go api/*/*/*.pb.gw.go api/*/*/*.swagger.json api/*/*/*.pb.validate.go
	rm -rf dist/sdk/*
	rm -rf third_party/swagger-ui/*.swagger.json
	rm -rf dist/swagger/*.swagger.json
	rm -rf *.swagger.json

.PHONY: server-image
server-image:
	docker buildx build -f tools/docker/gatepoint-server/Dockerfile -t gatepoint/gatepoint:latest . --platform linux/amd64,linux/arm64 --push
