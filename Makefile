CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

EXE =
ifeq ($(shell go env GOOS),windows)
EXE = .exe
endif

## The following tasks delegate to `script/build.go` so they can be run cross-platform.

.PHONY: bin/goctl$(EXE)
bin/goctl$(EXE): script/build$(EXE)
	@script/build$(EXE) $@

script/build$(EXE): script/build.go
ifeq ($(EXE),)
	GOOS= GOARCH= GOARM= GOFLAGS= CGO_ENABLED= go build -o $@ $<
else
	go build -o $@ $<
endif

.PHONY: clean
clean: script/build$(EXE)
	@$< $@

.PHONY: manpages
manpages: script/build$(EXE)
	@$< $@

.PHONY: completions
completions: bin/goctl$(EXE)
	mkdir -p ./share/bash-completion/completions ./share/fish/vendor_completions.d ./share/zsh/site-functions
	bin/goctl$(EXE) completion -s bash > ./share/bash-completion/completions/goctl
	bin/goctl$(EXE) completion -s fish > ./share/fish/vendor_completions.d/goctl.fish
	bin/goctl$(EXE) completion -s zsh > ./share/zsh/site-functions/_goctl

# just a convenience task around `go test`
.PHONY: test
test:
	go test ./...

## Site-related tasks are exclusively intended for use by the GitHub CLI team and for our release automation.

site:
	git clone https://github.com/github/cli.github.com.git "$@"

.PHONY: site-docs
site-docs: site
	git -C site pull
	git -C site rm 'manual/gh*.md' 2>/dev/null || true
	go run ./cmd/gen-docs --website --doc-path site/manual
	rm -f site/manual/*.bak
	git -C site add 'manual/gh*.md'
	git -C site commit -m 'update help docs' || true

.PHONY: site-bump
site-bump: site-docs
ifndef GITHUB_REF
	$(error GITHUB_REF is not set)
endif
	sed -i.bak -E 's/(assign version = )".+"/\1"$(GITHUB_REF:refs/tags/v%=%)"/' site/index.html
	rm -f site/index.html.bak
	git -C site commit -m '$(GITHUB_REF:refs/tags/v%=%)' index.html

## Install/uninstall tasks are here for use on *nix platform. On Windows, there is no equivalent.

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin
datadir := ${prefix}/share
mandir  := ${datadir}/man

.PHONY: install
install: bin/goctl manpages completions
	install -d ${DESTDIR}${bindir}
	install -m755 bin/goctl ${DESTDIR}${bindir}/
	install -d ${DESTDIR}${mandir}/man1
	install -m644 ./share/man/man1/* ${DESTDIR}${mandir}/man1/
	install -d ${DESTDIR}${datadir}/bash-completion/completions
	install -m644 ./share/bash-completion/completions/goctl ${DESTDIR}${datadir}/bash-completion/completions/goctl
	install -d ${DESTDIR}${datadir}/fish/vendor_completions.d
	install -m644 ./share/fish/vendor_completions.d/goctl.fish ${DESTDIR}${datadir}/fish/vendor_completions.d/goctl.fish
	install -d ${DESTDIR}${datadir}/zsh/site-functions
	install -m644 ./share/zsh/site-functions/_goctl ${DESTDIR}${datadir}/zsh/site-functions/_goctl

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/goctl ${DESTDIR}${mandir}/man1/gh.1 ${DESTDIR}${mandir}/man1/goctl-*.1
	rm -f ${DESTDIR}${datadir}/bash-completion/completions/goctl
	rm -f ${DESTDIR}${datadir}/fish/vendor_completions.d/goctl.fish
	rm -f ${DESTDIR}${datadir}/zsh/site-functions/_goctl
