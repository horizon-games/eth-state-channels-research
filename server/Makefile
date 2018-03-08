SHELL=bash -o pipefail
TEST_FLAGS ?=

all:
	@echo "***************************************************************************"
	@echo "**                      DGame Server Build Tool                          **"
	@echo "***************************************************************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo ""
	@echo " + Development:"
	@echo "   - tools"
	@echo "   - bootstrap"
	@echo ""
	@echo "   - run"
	@echo "   - run-wsrelay"
	@echo ""
	@echo "   - test"
	@echo "   - test-wsrelay"
	@echo ""
	@echo " + Builds:"
	@echo "   - build"
	@echo "   - clean"
	@echo "   - generate"
	@echo ""
	@echo "   - build-wsrelay"
	@echo ""
	@echo "   - dist"
	@echo "   - docs"
	@echo "   - todo"
	@echo ""
	@echo " + Dep management:"
	@echo "   - dep"
	@echo "   - dep-ensure"
	@echo ""

##
## Tools
##
tools:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/pkieltyka/tools/cmd/goimports
	go get -u github.com/VojtechVitek/rerun/cmd/rerun
	go get -u github.com/pressly/sup/cmd/sup


##
## Development
##
bootstrap: tools copy-sample-config

copy-sample-config:
	for i in wsrelay; do \
		test -f ./etc/$$i.conf || cp ./etc/$$i.conf{.sample,}; \
	done

define run
	@rerun -watch ./ -ignore vendor bin tests -run sh -c 'GOGC=off go build -i -o ./bin/$(1) ./cmd/$(1)/main.go && ./bin/$(1) -config=etc/$(1).conf'
endef

run:
	$(call run,wsrelay)

run-wsrelay:
	$(call run,wsrelay)

test: 
	@$(MAKE) -j1 test-wsrelay

test-wsrelay:
	GOGC=off go test $(TEST_FLAGS) ./services/wsrelay/tests/...

test-pkgs: generate
	GOGC=off go test $(TEST_FLAGS) $$(go list ./... | grep -v -e /vendor/ -e /tests -e /cmd/) | grep -v '\[no test files\]'

define run-test
	rerun -watch ./ -ignore vendor bin cmd -run go test -v -run=$(TEST) $(1)
endef


##
## Builds
##
GITTAG := $(shell git describe --exact-match --tags HEAD 2>/dev/null || :)
GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || :)
LONGVERSION := $(shell git describe --tags --long --abbrev=8 --always HEAD)$(echo -$GITBRANCH | tr / - | grep -v '\-master' || :)
VERSION := $(if $(GITTAG),$(GITTAG),$(LONGVERSION))
GITCOMMIT := $(shell git log -1 --date=iso --pretty=format:%H)
GITCOMMITDATE := $(shell git log -1 --date=iso --pretty=format:%cd)
GITCOMMITAUTHOR := $(shell git log -1 --date=iso --pretty="format:%an")

define build
	GOGC=off GOBIN=$$PWD/bin \
	go install -v \
		-tags='$(BUILDTAGS)' \
		-gcflags='-e' \
		-ldflags='-X "github.com/horizon-games/dgame-server.VERSION=$(VERSION)" -X "github.com/horizon-games/dgame-server.GITBRANCH=$(GITBRANCH)" -X "github.com/horizon-games/davatar.GITCOMMIT=$(GITCOMMIT)" -X "github.com/horizon-games/dgame-server.GITCOMMITDATE=$(GITCOMMITDATE)" -X "github.com/horizon-games/dgame-server.GITCOMMITAUTHOR=$(GITCOMMITAUTHOR)"' \
		$(1)
endef

build: build-wsrelay

build-wsrelay:
	$(call build, ./cmd/wsrelay)

clean:
	rm -rf $$GOPATH/pkg/*/github.com/horizon-games/dgame-server{,.*}
	rm -rf ./bin/*

dist: BUILDTAGS += production
dist: clean build

docs:
	@echo "TODO"

todo:
	@git grep TODO -- './*' ':!./vendor/' ':!./Makefile' || :


##
## Dependency mgmt
##
.PHONY: dep
$(GOPATH)/bin/dep:
	go get -u github.com/golang/dep/cmd/dep

dep: $(GOPATH)/bin/dep
	dep ensure

dep-status: $(GOPATH)/bin/dep
	dep status
