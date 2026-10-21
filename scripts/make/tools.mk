# Tool versions
GOLANGCI_LINT_VERSION := v2.13.2
GORELEASER_VERSION    := v2.15.2
GOTESTSUM_VERSION     := v1.13.0
GOSEC_VERSION         := v2.29.0
GOCOBERTURA_VERSION   := latest
GOMARKDOC_VERSION     := latest
SYFT_VERSION          := v1.51.1
GRYPE_VERSION         := v0.118.0
GOVULNCHECK_VERSION   := v1.7.0
GIT_CLIFF_VERSION     := v2.14.2

# Tools live in ./.gobin (gitignored) so every target runs the pinned version,
# not whatever PATH has. Each is a file target, installed on first use; bumping
# a version above does not reinstall it, `make clean-tools` does.

GORELEASER  := $(GOBIN_DIR)/goreleaser$(EXE)
GOLANGCI    := $(GOBIN_DIR)/golangci-lint$(EXE)
GOTESTSUM   := $(GOBIN_DIR)/gotestsum$(EXE)
GOSEC       := $(GOBIN_DIR)/gosec$(EXE)
GOCOBERTURA := $(GOBIN_DIR)/gocover-cobertura$(EXE)
GOMARKDOC   := $(GOBIN_DIR)/gomarkdoc$(EXE)
SYFT        := $(GOBIN_DIR)/syft$(EXE)
GRYPE       := $(GOBIN_DIR)/grype$(EXE)
GOVULNCHECK := $(GOBIN_DIR)/govulncheck$(EXE)
GITCLIFF    := $(GOBIN_DIR)/git-cliff$(EXE)

TOOLS := $(GORELEASER) $(GOLANGCI) $(GOTESTSUM) $(GOSEC) $(GOCOBERTURA) $(GOMARKDOC) $(SYFT) $(GRYPE) $(GOVULNCHECK) $(GITCLIFF)

# `go install` has no destination flag, only GOBIN. Exported for the tool rules
# alone: file-wide it would also redirect `install`, which must put devcontainerwizard in
# the user's own GOBIN.
$(TOOLS): export GOBIN = $(GOBIN_DIR)

.PHONY: tools clean-tools

$(GORELEASER):
	@echo "Installing goreleaser $(GORELEASER_VERSION)..."
	@go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION)

$(GOLANGCI):
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GOTESTSUM):
	@echo "Installing gotestsum $(GOTESTSUM_VERSION)..."
	@go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

$(GOSEC):
	@echo "Installing gosec $(GOSEC_VERSION)..."
	@go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)

$(GOCOBERTURA):
	@echo "Installing gocover-cobertura $(GOCOBERTURA_VERSION)..."
	@go install github.com/t-yuki/gocover-cobertura@$(GOCOBERTURA_VERSION)

$(GOMARKDOC):
	@echo "Installing gomarkdoc $(GOMARKDOC_VERSION)..."
	@go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@$(GOMARKDOC_VERSION)

$(SYFT):
	@echo "Installing syft $(SYFT_VERSION)..."
	@go install github.com/anchore/syft/cmd/syft@$(SYFT_VERSION)

$(GRYPE):
	@echo "Installing grype $(GRYPE_VERSION)..."
	@go install github.com/anchore/grype/cmd/grype@$(GRYPE_VERSION)

$(GOVULNCHECK):
	@echo "Installing govulncheck $(GOVULNCHECK_VERSION)..."
	@go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

# git-cliff is Rust, so go install cannot build it: the rule unpacks the pinned
# release for this OS and CPU. Windows curl and tar go by path, as make looks
# beside itself first, where Gow keeps ones too old for GitHub and for .zip.
ifeq ($(OS),Windows_NT)
GIT_CLIFF_TARGET = $(if $(filter ARM64,$(or $(PROCESSOR_ARCHITEW6432),$(PROCESSOR_ARCHITECTURE))),aarch64,x86_64)-pc-windows-msvc.zip
# Git Bash hands make SystemRoot as SYSTEMROOT.
SYSTEM32 := $(subst \,/,$(or $(SystemRoot),$(SYSTEMROOT)))/System32
CURL := $(SYSTEM32)/curl.exe
TAR  := $(SYSTEM32)/tar.exe
else
GIT_CLIFF_TARGET = $(if $(filter arm64 aarch64,$(shell uname -m)),aarch64,x86_64)-$(if $(filter Darwin,$(shell uname -s)),apple-darwin,unknown-linux-musl).tar.gz
CURL := curl
TAR  := tar
endif

$(GITCLIFF):
	@echo "Installing git-cliff $(GIT_CLIFF_VERSION)..."
	@$(CURL) -sSfL --create-dirs -o $(GOBIN_DIR)/git-cliff.archive https://github.com/orhun/git-cliff/releases/download/$(GIT_CLIFF_VERSION)/git-cliff-$(GIT_CLIFF_VERSION:v%=%)-$(GIT_CLIFF_TARGET)
	@$(TAR) -xf $(GOBIN_DIR)/git-cliff.archive -C $(GOBIN_DIR) --strip-components=1 git-cliff-$(GIT_CLIFF_VERSION:v%=%)/git-cliff$(EXE)
	@rm -f $(GOBIN_DIR)/git-cliff.archive

tools: $(TOOLS) ## Install every pinned tool into ./.gobin
	@echo "Tools installed in $(GOBIN_DIR)"

clean-tools: ## Remove ./.gobin so the next target reinstalls the pinned tools
	@echo "Removing $(GOBIN_DIR)..."
	@rm -rf $(GOBIN_DIR)
