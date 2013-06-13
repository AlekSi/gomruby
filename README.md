GoMRuby[![Build Status](https://secure.travis-ci.org/AlekSi/gomruby.png)](https://travis-ci.org/AlekSi/gomruby) [![Is maintained?](http://stillmaintained.com/AlekSi/gomruby.png)](http://stillmaintained.com/AlekSi/gomruby)
=======

Package gomruby embeds mruby (mini Ruby) VM into Go.

[Documentation](http://godoc.org/github.com/AlekSi/gomruby).

Installation
------------
It's slightly more than just `go get`:

    go get -d github.com/AlekSi/gomruby
    cd $GOPATH/src/github.com/AlekSi/gomruby
    make

mruby is built statically, use gomruby as typical Go package.
