ALL_JSONNET_FILES = find . -type d \( -name 'vendor' \) -prune -false -o \( -name '*.jsonnet' -o -name '*.libsonnet' \)

test:
	@ jt
fmt:
	@ $(ALL_JSONNET_FILES) | xargs jsonnetfmt -i

lint:
	@ $(ALL_JSONNET_FILES) | xargs -I {} jsonnet-lint {}

.PHONY: test fmt lint
