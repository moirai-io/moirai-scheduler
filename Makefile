
# Image URL to use all building/pushing image targets
CONTROLLER_IMG ?= ghcr.io/rudeigerc/moirai-controller:latest
IMG ?= ghcr.io/rudeigerc/moirai-scheduler:latest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.24.1

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build-controller build-scheduler build-cli

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./pkg/... ./controllers/... -coverprofile cover.out

E2ETEST_ASSETS_DIR=$(shell pwd)/testdir
.PHONY: e2e-test
e2e-test: kustomize ## Run e2e tests.
	mkdir -p $(E2ETEST_ASSETS_DIR)
	$(KUSTOMIZE) build config/crd > ${E2ETEST_ASSETS_DIR}/moirai-scheduler.crds.yaml
	$(KUSTOMIZE) build config/default > ${E2ETEST_ASSETS_DIR}/moirai-scheduler.yaml
	IMG=${IMG} E2ETEST_ASSETS_DIR=${E2ETEST_ASSETS_DIR} go test -tags=e2e -v ./test/e2e/controller
	@rm -rf ${E2ETEST_ASSETS_DIR}

ETCD_PATH=$(shell pwd)/third_party/etcd
.PHONY: test-scheduler-perf
test-scheduler-perf: ## Run Kubernetes scheduler performance tests.
ifeq (,$(wildcard ${ETCD_PATH}))
	@echo "Installing etcd"
	./hack/install-etcd.sh
endif
	PATH="${ETCD_PATH}:${PATH}" go test ./test/integration/scheduler_perf -alsologtostderr=false -logtostderr=false -run=^$$ -benchtime=1ns -bench=BenchmarkPerfScheduling/SchedulingBasic/5000Nodes/5000InitPods/1000PodsToSchedule

##@ Build

LDFLAGS=-ldflags '\
	-X k8s.io/component-base/version.gitVersion=v1.24.0-moirai-scheduler-$(shell date +%Y%m%d) \
	-X k8s.io/component-base/version.gitCommit=$(shell git rev-parse HEAD) \
	-X k8s.io/component-base/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
	'

.PHONY: build-controller
build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: build-scheduler
build-scheduler: ## Build Moirai scheduler binary.
	go build ${LDFLAGS} -o bin/moirai-scheduler cmd/scheduler/main.go

.PHONY: build-cli
build-cli: ## Build Moirai command line tool binary.
	go build -o bin/moiraictl cmd/cli/main.go

.PHONY: run-controller
run-controller: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: run-scheduler
run-scheduler: manifests generate fmt vet ## Run a scheduler from your host.
	go run ${LDFLAGS} ./cmd/scheduler/main.go --config manifests/scheduler/scheduler-config.yaml -v 5

.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker build -f ./build/scheduler/Dockerfile -t ${IMG} .
	docker build -f ./build/controller/Dockerfile -t ${CONTROLLER_IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}
	docker push ${CONTROLLER_IMG}

.PHONY: clean
clean:
	rm -rf ./bin

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN): ## Ensure that the directory exists
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.9.0

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
