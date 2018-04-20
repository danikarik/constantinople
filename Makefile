APP ?= cbkz-server

M = $(shell printf "\033[34;1m▶\033[0m")
HAS_DEP := $(shell command -v dep;)

all: vendor test build done

done:
	$(info $(M) done.)

.PHONY: vendor
vendor: prepare-dep ## Install dependicies
	$(info $(M) installing dependencies...)
	dep ensure

.PHONY: prepare-dep
prepare-dep: ## Install dep package manager
ifndef HAS_DEP
	$(info $(M) installing dep...)
	go get -u -v -d github.com/golang/dep/cmd/dep && \
	go install -v github.com/golang/dep/cmd/dep
endif
	
.PHONY: clean
clean:
	$(info $(M) cleaning build...)
	@rm -f bin/${APP}

.PHONY: build
build: clean ## Build program binary
	$(info $(M) building program...)
	go build -o bin/${APP} ./cmd/server...

.PHONY: run
run: ## Run in debug mode
	$(info $(M) running server...)
	bin/$(APP)

.PHONY: test
test: ## Run test
	$(info $(M) running test...)
	go test ./...

.PHONY: docker-build
docker-build: ## Build docker image.
	$(info $(M) building docker image...)
	docker-compose build

.PHONY: docker-up
docker-up: ## Run docker container.
	$(info $(M) running docker container...)
	docker-compose up -d

.PHONY: docker-start
docker-start: ## Start docker container.
	$(info $(M) starting docker container...)
	docker-compose start

.PHONY: docker-stop
docker-stop: ## Stop docker conrainer.
	$(info $(M) stopping docker container...)
	docker-compose stop

.PHONY: docker-down
docker-down: ## Remove docker image and container.
	$(info $(M) removing docker image...)
	docker-compose down --rmi all

.PHONY: docker-host
docker-host: ## Define docker host address.
	$(info $(M) finding host machine ip...)
	ip route | awk '/default/ { print $3 }'

.PHONY: help
help: ## Show usage
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
