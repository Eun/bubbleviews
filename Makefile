ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
.PHONY: linter
linter:
	-mkdir ${ROOT_DIR}/.make
	envsubst < .golangci.yml > ${ROOT_DIR}/.make/.golangci.yml
	docker run --rm --workdir /work --volume ${ROOT_DIR}:/work golangci/golangci-lint:v1.50.1 golangci-lint run --config /work/.make/.golangci.yml
