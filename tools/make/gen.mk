.PHONY: gen-labels
gen-labels:
	@$(LOG_TARGET)
	go run ./tools/annotations_prep/main.go --input ./api/label/label.yaml --output ./api/label/labels.gen.go --collection_type label

.PHONY: gen-annotations
gen-annotations:
	@$(LOG_TARGET)
	go run ./tools/annotations_prep/main.go --input ./api/annotation/annotations.yaml --output ./api/annotation/annotations.gen.go --collection_type annotation
