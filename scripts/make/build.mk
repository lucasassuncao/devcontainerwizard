.PHONY: build build-all release tag install run

build: $(GORELEASER) ## Build binary with goreleaser (current platform only)
	@echo "Building..."
	@$(GORELEASER) build --skip=validate --single-target --snapshot --clean

build-all: $(GORELEASER) ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@$(GORELEASER) build --skip=validate --snapshot --clean

release: $(GORELEASER) ## Create a release with goreleaser
	@echo "Creating release..."
	@$(GORELEASER) release --timeout 360s

# git-cliff writes the notes from cliff.toml into a file, not a pipe: git tag
# accepts an empty message, so a failed git-cliff would still have tagged.
tag: $(GITCLIFF) ## Create and push a git tag annotated with its release notes (usage: make tag VERSION=v1.2.3)
ifndef VERSION
	$(error Usage: make tag VERSION=v1.2.3)
endif
	git diff --exit-code --quiet
	$(GITCLIFF) --unreleased --tag $(VERSION) --output .git/TAG_NOTES
	git tag -a --cleanup=whitespace -F .git/TAG_NOTES $(VERSION)
	git push origin $(VERSION)

install: ## Install binary globally
	@go install

run: ## Run the application
	@go run $(MAIN_PATH)
