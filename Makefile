export CFLAGS=`go env GOGCCFLAGS`

test:
	gofmt -e -s -w .
	go tool vet .
	LD_LIBRARY_PATH=`pwd` go test -v -gocheck.v

common:
	cd mruby && git clean -xdf && make
	go get -u launchpad.net/gocheck

linux: common
	ld --whole-archive -shared -o libmruby.so mruby/build/host/lib/libmruby.a

mac: common
	# Explodes. -force_load? -all_load?
	ld -dylib -o libmruby.dylib mruby/build/host/lib/libmruby.a
