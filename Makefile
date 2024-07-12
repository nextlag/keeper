ALL_TARGETS := client app
GO := go
LINTER := staticcheck

define build
	${GO} build -o ./cmd/$(1)/ -v ./cmd/$(1)/
endef

define clean
	rm -f ./cmd/$(1)/$(1)
endef

.PHONY: all build build-% clean clean-% test lint generate test_with_coverage coverage_total coverage_html

all: build

build: $(ALL_TARGETS:%=build-%)

build-%:
	@echo === Building $*
	$(call build,$*)

clean: $(ALL_TARGETS:%=clean-%)

clean-%:
	@echo === Cleaning $*
	$(call clean,$*)

test:
	@echo === Running tests
	${GO} test -count=1 -v -cover ./...

lint:
	@echo === Running linter
	$(LINTER) --version
	$(LINTER) ./...

generate:
	${GO} generate ./...

test_with_coverage:
	${GO} test -count=1 -coverprofile=coverage.out ./...

coverage_total: test_with_coverage
	@echo Total coverage: $$( ${GO} tool cover -func=coverage.out | awk '/total:/ {print $$3}' )

coverage_html: test_with_coverage
	${GO} tool cover -html=coverage.out
