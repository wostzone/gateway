# Makefile to build and test the WoST Hub
# To build including generating certificates: make all
DIST_FOLDER=./dist
PKG_NAME=wost-hub.tgz
.DEFAULT_GOAL := help

.PHONY: 

all: hub gencerts mosqauth mosquittomgr idprov-pb ## Build hub and core apps

install:  ## Install the hub into ~/bin/wost/bin and config
	mkdir -p ~/bin/wost/bin
	mkdir -p ~/bin/wost/config
	mkdir -p ~/bin/wost/logs
	cp $(DIST_FOLDER)/bin/* ~/bin/wost/bin/
	cp $(DIST_FOLDER)/config/mosquitto.conf.template ~/bin/wost/config/  
	cp -n $(DIST_FOLDER)/config/* ~/bin/wost/config/  

dist: clean   ## Build binary distribution tarball 
		tar -czf $(PKG_NAME) -C $(DIST_FOLDER) .

test: .PHONY ## Run tests (todo fix this)
	go test -failfast -v ./...

clean: ## Clean distribution files
	go clean
	go mod tidy
	rm -f test/certs/*
	rm -f test/logs/*
	rm -f test/config/mosquitto.conf
	rm -f $(DIST_FOLDER)/certs/*
	rm -f $(DIST_FOLDER)/logs/*
	rm -f $(DIST_FOLDER)/bin/*
	rm -f $(DIST_FOLDER)/config/mosquitto.conf
	rm -f debug $(PKG_NAME)
	mkdir -p $(DIST_FOLDER)/bin
	mkdir -p $(DIST_FOLDER)/certs
	mkdir -p $(DIST_FOLDER)/config
	mkdir -p $(DIST_FOLDER)/logs

#deps: ## Build GO dependencies 
#		go get

#upgrade: ## Upgrade the dependencies to the latest version. Use with care.
#		go fix

#prof: ## Run application with CPU and memory profiling
#	  go run main.go -cpuprofile=cpu.prof -memprofile=mem.prof

idprov-pb: ## Build idprov-pb plugin
	go build -o $(DIST_FOLDER)/bin/$@ ./cmd/$@/main.go

mosquittomgr: mosqauth ## Build mosquitto configuration manager
	go build -o $(DIST_FOLDER)/bin/$@ ./cmd/$@/main.go
	
mosqauth: ## Build mosquitto auth plugin for use by the Hub
	cd cmd/mosqauth && make
	@echo "> SUCCESS. The executable '$@' can be found in $(DIST_FOLDER)/bin/$@"

gencerts: ## Build gencerts utility to generate self-signed certificates with CA
	go build -o $(DIST_FOLDER)/bin/$@ ./cmd/$@/main.go
	@echo "> SUCCESS. The executable '$@' can be found in $(DIST_FOLDER)/bin/$@"

hub: ## Build WoST Hub
	go build -o $(DIST_FOLDER)/bin/$@ ./cmd/$@/main.go
	@echo "> SUCCESS. The executable '$@' can be found in $(DIST_FOLDER)/bin/$@"


#docker: ## Build hub for Docker target (TODO untested)
#		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_AMD64)" -v

help: ## Show this help
		@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

