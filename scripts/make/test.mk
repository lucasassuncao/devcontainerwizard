# Plain mkdir as a directory target: on Windows the recipe may run under cmd.exe,
# where `mkdir -p` creates a directory named "-p", and $(wildcard) in an ifeq
# expands at parse time, before `make clean test-coverage` has deleted it.
$(COVERAGE_DIR):
	@mkdir $(COVERAGE_DIR)

.PHONY: test test-watch test-coverage fuzz

# -shuffle=on randomises test order so a test that only passes because an
# earlier one ran first fails here rather than in someone else's branch.
test: $(GOTESTSUM) ## Run tests with gotestsum (testdox format)
	@$(GOTESTSUM) --format testdox -- -race -shuffle=on ./...

test-watch: $(GOTESTSUM) ## Run tests in watch mode (reruns on file changes)
	@$(GOTESTSUM) --format testdox --watch -- -race ./...

# Order-only (after the |): writing a report changes the directory's timestamp.
test-coverage: $(GOTESTSUM) $(GOCOBERTURA) | $(COVERAGE_DIR) ## Run tests with coverage (HTML + Cobertura XML)
	@$(GOTESTSUM) --format testdox -- -race -shuffle=on -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	@go tool cover -func=$(COVERAGE_OUT) | tail -1
	@go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@$(GOCOBERTURA) < $(COVERAGE_OUT) > $(COVERAGE_XML)
	@echo "Reports: $(COVERAGE_HTML) | $(COVERAGE_XML)"

# The config file is the only thing a user hand-writes, so it is this tool's
# untrusted input. FUZZTIME overrides the duration: make fuzz FUZZTIME=5m
FUZZTIME ?= 60s
fuzz: ## Fuzz the config parse pipeline (override with FUZZTIME=5m)
	@go test ./internal/devcontainer/ -run '^$$' -fuzz FuzzConfigPipeline -fuzztime $(FUZZTIME)
