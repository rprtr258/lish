CMD = ./main.go
RUN = go run -race ${CMD}
LDFLAGS = -ldflags="-s -w"

all: run test

# run standard samples
run:
	go build -o ink -race ${CMD}
	./ink samples/fizzbuzz.ink
	./ink samples/graph.ink
	./ink samples/basic.ink
	./ink samples/kv.ink
	./ink samples/html.ink
	./ink samples/fib.ink
	./ink samples/collatz.ink
	./ink samples/newton.ink
	./ink samples/pi.ink
	./ink samples/prime.ink
	./ink samples/quicksort.ink
	./ink samples/pingpong.ink
	./ink samples/undefinedme.ink || true
	./ink samples/error.ink || true
	./ink samples/exec.ink
	echo 'Linus' | ./ink samples/prompt.ink # we echo in some input for prompt.ink testing stdin
	./ink samples/io_test.ink
	rm ./ink

repl:
	${RUN}

# run just the minimal test suite
test-mini:
	${RUN} samples/test.ink

# run standard test suites
test:
	go test ./...
	go build -o ink -race ${CMD}
	./ink samples/mangled.ink
	./ink samples/test.ink
	# run I/O test under isolated mode -- all ops should still return valid responses
	# We copy the file in question -- eval.go -- to a temporary location, since
	# no-read and no-write I/O operations will delete the file.
	./ink samples/io.ink
	./ink samples/pingpong.ink
	./ink samples/exec.ink
	./ink samples/grapher.ink
	./ink samples/graph.ink
	./ink samples/rainbow.ink 'YAY SHALAVI'
	./ink samples/img.ink
	# test -eval flag
	./ink -eval "log:=import('samples/logging.ink').log,f:=x=>()=>log('Eval test: '+x),f('passed!.ink')()"
	rm ./ink

# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build ${LDFLAGS} -o ink-$* ${CMD}

# build for all OS targets, useful for releases
build: build-linux build-darwin build-windows build-openbsd

# install on host system
install:
	#cp utils/ink.vim ~/.vim/syntax/ink.vim
	go install ${LDFLAGS} .
	ls -l `which ink`

# pre-commit hook
precommit:
	go vet .
	go fmt .

# clean any generated files
clean:
	rm -rvf *.bmp ink ink-*
