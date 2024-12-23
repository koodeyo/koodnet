KOODNET_CMD_PATH = "./cmd/koodnet"
CGO_ENABLED = 0
export CGO_ENABLED

# Set up OS specific bits
ifeq ($(OS),Windows_NT)
	CMD_SUFFIX = .exe
	NULL_FILE = nul
	# RIO on windows does pointer stuff that makes go vet angry
	VET_FLAGS = -unsafeptr=false
else
	CMD_SUFFIX =
	NULL_FILE = /dev/null
endif

# Only defined the build number if we haven't already
ifndef BUILD_NUMBER
	ifeq ($(shell git describe --exact-match 2>$(NULL_FILE)),)
		BUILD_NUMBER = $(shell git describe --abbrev=0 --match "v*" | cut -dv -f2)-$(shell git branch --show-current)-$(shell git describe --long --dirty | cut -d- -f2-)
	else
		BUILD_NUMBER = $(shell git describe --exact-match --dirty | cut -dv -f2)
	endif
endif

KOODNET_DOCKER_IMAGE_REPO ?= koodeyo/koodnet
KOODNET_API_DOCKER_IMAGE_REPO ?= koodeyo/koodnet-api
DOCKER_IMAGE_TAG ?= latest

LDFLAGS = -X main.Build=$(BUILD_NUMBER)

ALL_LINUX = linux-amd64 \
	linux-386 \
	linux-ppc64le \
	linux-arm-5 \
	linux-arm-6 \
	linux-arm-7 \
	linux-arm64 \
	linux-mips \
	linux-mipsle \
	linux-mips64 \
	linux-mips64le \
	linux-mips-softfloat \
	linux-riscv64 \
	linux-loong64

ALL_FREEBSD = freebsd-amd64 \
	freebsd-arm64

ALL_OPENBSD = openbsd-amd64 \
	openbsd-arm64

ALL_NETBSD = netbsd-amd64 \
 	netbsd-arm64

ALL = $(ALL_LINUX) \
	$(ALL_FREEBSD) \
	$(ALL_OPENBSD) \
	$(ALL_NETBSD) \
	darwin-amd64 \
	darwin-arm64 \
	windows-amd64 \
	windows-arm64


DOCKER_BIN = build/linux-amd64/koodnet build/linux-amd64/koodnet-api

all: $(ALL:%=build/%/koodnet) $(ALL:%=build/%/koodnet-api)

docker: docker/linux-$(shell go env GOARCH)

release: $(ALL:%=build/koodnet-%.tar.gz)

release-linux: $(ALL_LINUX:%=build/koodnet-%.tar.gz)

release-freebsd: $(ALL_FREEBSD:%=build/koodnet-%.tar.gz)

release-openbsd: $(ALL_OPENBSD:%=build/koodnet-%.tar.gz)

release-netbsd: $(ALL_NETBSD:%=build/koodnet-%.tar.gz)

release-boringcrypto: build/koodnet-linux-$(shell go env GOARCH)-boringcrypto.tar.gz

BUILD_ARGS += -trimpath

bin-windows: build/windows-amd64/koodnet.exe build/windows-amd64/koodnet-api.exe
	mv $? .

bin-windows-arm64: build/windows-arm64/koodnet.exe build/windows-arm64/koodnet-api.exe
	mv $? .

bin-darwin: build/darwin-amd64/koodnet build/darwin-amd64/koodnet-api
	mv $? .

bin-freebsd: build/freebsd-amd64/koodnet build/freebsd-amd64/koodnet-api
	mv $? .

bin-freebsd-arm64: build/freebsd-arm64/koodnet build/freebsd-arm64/koodnet-api
	mv $? .

bin-boringcrypto: build/linux-$(shell go env GOARCH)-boringcrypto/koodnet build/linux-$(shell go env GOARCH)-boringcrypto/koodnet-api
	mv $? .

bin-pkcs11: BUILD_ARGS += -tags pkcs11
bin-pkcs11: CGO_ENABLED = 1
bin-pkcs11: bin

bin:
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./koodnet${CMD_SUFFIX} ${KOODNET_CMD_PATH}
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./koodnet-api${CMD_SUFFIX} ./cmd/koodnet-api

install:
	go install $(BUILD_ARGS) -ldflags "$(LDFLAGS)" ${KOODNET_CMD_PATH}
	go install $(BUILD_ARGS) -ldflags "$(LDFLAGS)" ./cmd/koodnet-api

build/linux-arm-%: GOENV += GOARM=$(word 3, $(subst -, ,$*))
build/linux-mips-%: GOENV += GOMIPS=$(word 3, $(subst -, ,$*))

# Build an extra small binary for mips-softfloat
build/linux-mips-softfloat/%: LDFLAGS += -s -w

# boringcrypto
build/linux-amd64-boringcrypto/%: GOENV += GOEXPERIMENT=boringcrypto CGO_ENABLED=1
build/linux-arm64-boringcrypto/%: GOENV += GOEXPERIMENT=boringcrypto CGO_ENABLED=1

build/%/koodnet: .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
		go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${KOODNET_CMD_PATH}

build/%/koodnet-api: .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
		go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ./cmd/koodnet-api

build/%/koodnet.exe: build/%/koodnet
	mv $< $@

build/%/koodnet-api.exe: build/%/koodnet-api
	mv $< $@

build/koodnet-%.tar.gz: build/%/koodnet build/%/koodnet-api
	tar -zcv -C build/$* -f $@ koodnet koodnet-api

build/koodnet-%.zip: build/%/koodnet.exe build/%/koodnet-api.exe
	cd build/$* && zip ../koodnet-$*.zip koodnet.exe koodnet-api.exe

docker/%: build/%/koodnet build/%/koodnet-api
	docker build . $(DOCKER_BUILD_ARGS) -f docker/Dockerfile.koodnet --platform "$(subst -,/,$*)" --tag "${KOODNET_DOCKER_IMAGE_REPO}:${DOCKER_IMAGE_TAG}" --tag "${KOODNET_DOCKER_IMAGE_REPO}:$(BUILD_NUMBER)"
	docker build . $(DOCKER_BUILD_ARGS) -f docker/Dockerfile.koodnet-api --platform "$(subst -,/,$*)" --tag "${KOODNET_API_DOCKER_IMAGE_REPO}:${DOCKER_IMAGE_TAG}" --tag "${KOODNET_API_DOCKER_IMAGE_REPO}:$(BUILD_NUMBER)"

service:
	@echo > $(NULL_FILE)
	$(eval KOODNET_CMD_PATH := "./cmd/koodnet-service")
ifeq ($(words $(MAKECMDGOALS)),1)
	@$(MAKE) service ${.DEFAULT_GOAL} --no-print-directory
endif

bin-docker: bin build/linux-amd64/koodnet build/linux-amd64/koodnet-api

test:
	go test -v ./... -race -cover

setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs

dev:
	@docker compose -f docker-compose.yml up

stop:
	@docker compose -f docker-compose.yml down

.FORCE:
.PHONY: dev stop setup bin release service
.DEFAULT_GOAL := bin
