all: libmruby.a
	gofmt -e -s -w .
	go tool vet .
	go test -v -gocheck.v

libmruby.a:
	git submodule update --init
	cd mruby && git reset --hard && git clean -xdf && make
	cp mruby/build/host/lib/libmruby.a .
	go get launchpad.net/gocheck

clean:
	rm libmruby.a
