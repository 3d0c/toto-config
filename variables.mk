### project name
BINARY				:= toto-config

### version
MAJOR				:= 0
MINOR				:= 0
PATCH				:= 1
PRODUCT_VERSION 	:= ${MAJOR}.${MINOR}.${PATCH}
COMMIT				:= $(shell git rev-parse --short HEAD)
FULL_VERSION		:= ${PRODUCT_VERSION}-${COMMIT}
DATE				:= $(shell date +%FT%T%z)

### go vars
GO_ENV_VARS	= GOSUMDB=off GO111MODULE=on
BUILD_FLAGS	= -ldflags "-X github.com/3d0c/toto-config/cmd/toto-config.version=${FULL_VERSION} -X github.com/3d0c/toto-config/cmd/toto-config.date=${DATE}"

### variables for testing
COVERAGE_OUT		:= coverage.out
COVERAGE_HTML		:= coverage.html
TMP					:= /tmp
