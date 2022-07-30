.PHONY: help
help: # show list of all commands
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: # run lish
	cargo run

.PHONY: test
test: # run unit tests
	cargo test

.PHONY: lint
lint: # run cargo check
	cargo check

.PHONY: precommit
precommit: # run precommit checks
	$(MAKE) check
	$(MAKE) test

.PHONY: todo
todo: # show list of all todos left in code
	@rg 'TODO' --glob '**/*.rs' || echo 'All done!'

.PHONY: build
build: # build executable
	cargo build
