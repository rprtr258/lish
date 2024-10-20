# run main binary
run:
	./bin/main.ink translate-full test/cases/*.ink

# run all tests under test/
.PHONY: test
test:
	./bin/main.ink run test/cases/*.ink

fmt:
	inkfmt fix bin/*.ink src/*.ink test/*.ink test/cases/*.ink

fmt-check:
	inkfmt bin/*.ink src/*.ink test/*.ink test/cases/*.ink
