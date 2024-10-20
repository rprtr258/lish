SEPTEMBER := ink /home/rprtr258/pr/september/bin/main.ink

all: ci

# run app server
run:
	ink src/main.ink

# build app clients
build:
	cat static/js/ink.js \
		static/js/torus.min.js \
		> static/ink/vendor.js
	${SEPTEMBER} translate \
		lib/stub.ink \
		vendor/std.ink \
		vendor/str.ink \
		vendor/quicksort.ink \
		lib/reader.ink \
		lib/md.ink \
		lib/torus.js.ink \
		src/app.js.ink \
		| tee /dev/stderr > static/ink/common.js
	${SEPTEMBER} translate src/config.js.ink \
		> static/ink/config.js
	${SEPTEMBER} translate src/static-config.js.ink \
		> static/ink/static-config.js
	cat \
		static/ink/vendor.js \
		static/ink/config.js \
		static/ink/common.js \
		> static/ink/bundle.js
	cat \
		static/ink/vendor.js \
		static/ink/static-config.js \
		static/ink/common.js \
		> static/ink/static-bundle.js

# build whenever Ink sources change
watch:
	ls lib/* src/* | entr make build

# run all tests under test/
.PHONY: test
test:
	ink ./test/main.ink

fmt:
	inkfmt fix lib/*.ink src/*.ink test/*.ink

# like test, but runs in CI with the bundled Ink interpreter
ci:
	./bin/ink ./test/main.ink
