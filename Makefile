GOPATH    ?= $(HOME)/go
GOBIN     ?= $(GOPATH)/bin
GOIMPORTS = $(GOBIN)/goimports
CMD       ?= bin/addon
IMG       ?= quay.io/konveyor/tackle2-addon-discovery:latest

PKG = ./cmd/...
PKGDIR = $(subst /...,,$(PKG))


cmd: fmt vet
	go build -ldflags="-w -s" -o ${CMD} github.com/konveyor/tackle2-addon-discovery/cmd

image-docker:
	docker build -t ${IMG} .

image-podman:
	podman build -t ${IMG} .

fmt: $(GOIMPORTS)
	$(GOIMPORTS) -w $(PKGDIR)

vet:
	go vet $(PKG)

# Ensure goimports installed.
$(GOIMPORTS):
	go install golang.org/x/tools/cmd/goimports@v0.24.0
