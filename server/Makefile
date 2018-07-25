DEPLOY_TARGET ?= dev
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
	@echo "   - run-arcadeum-server"
	@echo ""
	@echo "   - test"
	@echo "   - test-matcher"
	@echo ""
	@echo " + Builds:"
	@echo "   - build"
	@echo "   - clean"
	@echo "   - generate"
	@echo ""
	@echo "   - build-arcadeum-server"
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
	# http://timmurphy.org/2015/09/27/how-to-get-a-makefile-directory-path/
	test -f ./etc/arcadeum-server.conf || sed -e \
		's|^working_dir   = ""$$|working_dir   = "$(dir $(realpath $(firstword $(MAKEFILE_LIST))))"|' \
		./etc/arcadeum-server.conf.sample > ./etc/arcadeum-server.conf

define run
	@rerun -watch ./ -ignore vendor bin tests -run sh -c 'GOGC=off go build -i -o ./bin/$(1) ./cmd/$(1)/main.go && ./bin/$(1) -config=etc/$(1).conf'
endef

run:
	$(call run,arcadeum-server)

arcadeum-server:
	$(call run,arcadeum-server)

test:
	@$(MAKE) -j1 test-matcher

test-matcher:
	GOGC=off go test $(TEST_FLAGS) ./matcher/tests/...

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
		-ldflags='-X "github.com/horizon-games/arcadeum/server.VERSION=$(VERSION)" -X "github.com/horizon-games/arcadeum/server.GITBRANCH=$(GITBRANCH)" -X "github.com/horizon-games/davatar.GITCOMMIT=$(GITCOMMIT)" -X "github.com/horizon-games/arcadeum/server.GITCOMMITDATE=$(GITCOMMITDATE)" -X "github.com/horizon-games/arcadeum/server.GITCOMMITAUTHOR=$(GITCOMMITAUTHOR)"' \
		$(1)
endef

build: build-arcadeum-server

build-arcadeum-server:
	$(call build, ./cmd/arcadeum-server)

clean:
	rm -rf $$GOPATH/pkg/*/github.com/horizon-games/arcadeum/server{,.*}
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

deploy:
	ansible-playbook -i ../../SkyWeaver/config/ansible.hosts -v playbook.yml -e "host=$(DEPLOY_TARGET)"
