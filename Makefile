GO_REPO ?= /golang
BUILDDIR ?= $(GO_REPO)/build
OUTDIR ?= $(CURDIR)/out
RELEASEDIR ?= $(CURDIR)/release
SBINDIR ?= usr/local/zabbix/share/zabbix/alertscripts
ETCDIR ?= etc/
DESTDIR ?= /

GO_ARCH_MAP_x86 := 386
GO_ARCH_MAP_x86_64 := amd64
GO_ARCH_MAP_arm := arm
GO_ARCH_MAP_arm64 := arm64
GO_ARCH_MAP_aarch64 := arm64
GO_ARCH_MAP_mips := mipsx
GO_ARCH_MAP_mips64 := mips64x

export GOOS := linux
export CGO_ENABLED := 0

default: amd64

GOBUILDARCH := $(GO_ARCH_MAP_$(shell uname -m))
GOBUILDOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOBUILDVERSION := 1.14.15
GOBUILDTARBALL := https://dl.google.com/go/go$(GOBUILDVERSION).$(GOBUILDOS)-$(GOBUILDARCH).tar.gz
GOBUILDVERSION_NEEDED := go version go$(GOBUILDVERSION) $(GOBUILDOS)/$(GOBUILDARCH)
export GOPROXY := https://goproxy.io
export GOROOT := $(BUILDDIR)/$(GOBUILDVERSION)/goroot
export GOPATH := $(BUILDDIR)/$(GOBUILDVERSION)/gopath
export PATH := $(GOROOT)/bin:$(PATH)
GOBUILDVERSION_CURRENT := $(shell $(GOROOT)/bin/go version 2>/dev/null)
ifneq ($(GOBUILDVERSION_NEEDED),$(GOBUILDVERSION_CURRENT))
$(shell rm -f $(GOROOT)/bin/go)
endif
$(GOROOT)/bin/go:
	rm -rf "$(GOROOT)"
	mkdir -p "$(GOROOT)"
	curl "$(GOBUILDTARBALL)" | tar -C "$(GOROOT)" --strip-components=1 -xzf - || { rm -rf "$(GOROOT)"; exit 1; }

$(shell test "$$(cat $(BUILDDIR)/.gobuildversion 2>/dev/null)" = "$(GOBUILDVERSION_CURRENT)" || rm -rf "$(OUTDIR)")

amd64: $(GOROOT)/bin/go
	go get -tags linux || { chmod -fR +w "$(GOPATH)/pkg/mod"; rm -rf "$(GOPATH)/pkg/mod"; exit 1; }
	chmod -fR +w "$(GOPATH)/pkg/mod"
	GOARCH=amd64 go build -tags linux -v
	GOARCH=amd64 go install -tags linux -v
	GOARCH=amd64 go build -tags linux -v -o $(OUTDIR)/$@/ytzabbixalert ytzabbixalert.go
	go version > $(BUILDDIR)/.gobuildversion


arm64: $(GOROOT)/bin/go
	go get -tags linux || { chmod -fR +w "$(GOPATH)/pkg/mod"; rm -rf "$(GOPATH)/pkg/mod"; exit 1; }
	chmod -fR +w "$(GOPATH)/pkg/mod"
	GOARCH=arm64 go build -tags linux -v
	GOARCH=arm64 go install -tags linux -v
	GOARCH=arm64 go build -tags linux -v -o $(OUTDIR)/$@/ytzabbixalert ytzabbixalert.go
	go version > $(BUILDDIR)/.gobuildversion

clean:
	rm -rf $(OUTDIR)
	rm -rf $(RELEASEDIR)

install-arm64:
	install -d $(DESTDIR)$(ETCDIR)
	install -d $(DESTDIR)$(SBINDIR)
	install -m 0755 ${OUTDIR}/arm64/ytzabbixalert $(DESTDIR)$(SBINDIR)

install-amd64:
	install -d $(DESTDIR)$(ETCDIR)
	install -d $(DESTDIR)$(SBINDIR)
	install -m 0755 ${OUTDIR}/amd64/ytzabbixalert $(DESTDIR)$(SBINDIR)

uninstall:
	rm -rf $(DESTDIR)$(SBINDIR)ytzabbixalert

env:
	go env