export CFLAGS=`go env GOGCCFLAGS`

test:
	gofmt -e -s -w .
	go tool vet .
	LD_LIBRARY_PATH=`pwd` go test -v -gocheck.v

rebuild_mruby:
	cd mruby && git clean -xdf && make

linux: rebuild_mruby
	ld --whole-archive -shared -o libmruby.so mruby/build/host/lib/libmruby.a

mac: rebuild_mruby
	# Explodes. -force_load? -all_load?
	ld -dylib -o libmruby.dylib mruby/build/host/lib/libmruby.a
