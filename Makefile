default: build

.PHONY: \
	build \
	gentls \
	help \
	version

gentls:
	mkdir -p cert
	openssl genrsa -out cert/server.key 2048
	openssl ecparam -genkey -name secp384r1 -out cert/server.key
	openssl req -new -x509 -sha256 -key cert/server.key -out cert/server.crt -days 3650

build: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

build-linux-amd64:
	@echo "\033[33mBuilding linux amd 64 \033[0m"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o bin/ftf-linux-amd-64
	@echo "\033[32mBuild successfully ftf-linux-amd-64 \033[0m\n"

build-linux-arm64:
	@echo "\033[33mBuilding linux arm 64"
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-w -s' -o bin/ftf-linux-arm-64
	@echo "\033[32mBuild successfully ftf-linux-arm-64 \033[0m\n"

build-darwin-amd64:
	@echo "\033[33mBuilding darwin amd 64"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o bin/ftf-darwin-amd-64
	@echo "\033[32mBuild successfully ftf-darwin-amd-64 \033[0m\n"

build-darwin-arm64:
	@echo "\033[33mBuilding darwin arm 64 (m1)"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags '-w -s' -o bin/ftf-darwin-arm-64
	@echo "\033[32mBuild successfully ftf-darwin-arm-64 (m1) \033[0m\n"

build-windows-amd64:
	@echo "\033[33mBuilding windows amd 64"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o bin/ftf-windows-amd-64.exe
	@echo "\033[32mBuild successfully ftf-windows-amd-64.exe \033[0m\n"

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen'
	@echo '    build              Compile a program into an executable file'
	@echo '    clean              Clean all executable files'
	@echo '    version            Display Go version'
	@echo ''
	@echo 'Targets run by default is: build'
	@echo ''

version:
	@go version