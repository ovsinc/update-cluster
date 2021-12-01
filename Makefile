_CURDIR := `git rev-parse --show-toplevel 2>/dev/null | sed -e 's/(//'`
_BUILD_DIR := "${_CURDIR}/build"


REGISTRY_URL = "127.0.0.1:5000"


API_IMG := "test-api"
API_IMG_TAG := "v1"
API_IMAGE := "${API_IMG}:${API_IMG_TAG}"


BACKEND_IMG := "test-backend"
BACKEND_IMG_TAG := "v1"
BACKEND_IMAGE := "${BACKEND_IMG}:${BACKEND_IMG_TAG}"

AGENT_OUT := "${_CURDIR}/build/backend"
API_OUT := "${_CURDIR}/build/api"




.PHONY: all
all: help


.PHONY: build
build: build_code build_docker ## Build all


.PHONY: build_code
build_code: build_backend build_api ## Build all code

.PHONY: build_backend
build_backend: ## Build the BACKEND server
	@CGO_ENABLED=0 \
		go build \
		-v -mod=vendor -installsuffix cgo \
		-o ${AGENT_OUT} \
		${_CURDIR}/cmd/backend

.PHONY: build_api
build_api: ## Build the API server
	@CGO_ENABLED=0 \
		go build \
		-v -mod=vendor -installsuffix cgo \
		-o ${API_OUT} \
		${_CURDIR}/cmd/api

.PHONY: build_docker
build_docker: ## Build and push the Docker images
	@docker build \
		--force-rm \
		--tag "${BACKEND_IMAGE}" \
		--no-cache \
		--file "${_CURDIR}/docker/backend.docker" \
		"${_BUILD_DIR}"
	@docker tag "${BACKEND_IMAGE}" "${BACKEND_IMG}:latest"
	@docker tag "${BACKEND_IMAGE}" "${REGISTRY_URL}/${BACKEND_IMAGE}"
	@docker tag "${BACKEND_IMAGE}" "${REGISTRY_URL}/${BACKEND_IMG}:latest"
	@echo "PUSH '${BACKEND_IMG}' -> '${REGISTRY_URL}'"
	@docker push "${REGISTRY_URL}/${BACKEND_IMAGE}"
	@docker push "${REGISTRY_URL}/${BACKEND_IMG}:latest"

	@docker build \
		--force-rm \
		--tag "${API_IMAGE}" \
		--no-cache \
		--file "${_CURDIR}/docker/api.docker" \
		"${_BUILD_DIR}"
	@docker tag "${API_IMAGE}" "${API_IMG}:latest"
	@docker tag "${API_IMAGE}" "${REGISTRY_URL}/${API_IMG}:latest"
	@docker tag "${API_IMAGE}" "${REGISTRY_URL}/${API_IMAGE}"
	@echo "PUSH '${API_IMAGE}' -> '${REGISTRY_URL}'"
	@docker push "${REGISTRY_URL}/${API_IMAGE}"
	@docker push "${REGISTRY_URL}/${API_IMG}:latest"



.PHONY: registry
registry: ## Run local registry
	@docker stack deploy -c "${_CURDIR}/docker/registry.yml" registry


.PHONY: start
start: ## Run services with TF
	@pushd "${_CURDIR}/tf" &>/dev/null && terraform apply -auto-approve  || popd &>/dev/null

.PHONY: stop
stop: ## Stop services with TF
	@pushd "${_CURDIR}/tf" &>/dev/null && terraform destroy -auto-approve || popd &>/dev/null


.PHONY: clean
clean: ## Clean
	@go clean
	@docker rmi "${BACKEND_IMAGE}" "${API_IMAGE}"



.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'