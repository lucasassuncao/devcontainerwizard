.PHONY: docs generate-docs show-examples

docs: $(GOMARKDOC) generate-docs ## Generate the reference docs and the Go package docs
	@echo "Generating package docs..."
	@$(GOMARKDOC) -e \
		--repository.url https://github.com/lucasassuncao/devcontainerwizard \
		--repository.default-branch main \
		--repository.path / \
		-o '{{.Dir}}/README.md' ./internal/...

generate-docs: ## Generate the devcontainer reference docs, examples, and JSON Schema
	@go run $(MAIN_PATH) generate-docs

show-examples: ## Browse the built-in presets in the terminal
	@go run $(MAIN_PATH) show-examples
