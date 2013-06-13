test:
	gofmt -e -s -w .
	go tool vet .
	go test -v -gocheck.v

prepare:
	cd mruby && git reset --hard && git clean -xdf && make
	cp mruby/build/host/lib/libmruby.a .
	go get -u launchpad.net/gocheck
