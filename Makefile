all: libmruby.a
	gofmt -e -s -w .
	go tool vet .
	go test -v -gocheck.v

libmruby.a:
	git submodule update --init
	cd mruby && git reset --hard
	cd mruby && git clean -xdf && make
	cp mruby/build/host/lib/libmruby.a .
	go get launchpad.net/gocheck
	-go get code.google.com/p/go.tools/cmd/vet

clean:
	-rm libmruby.a
