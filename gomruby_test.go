package gomruby_test

import (
	. "."
	"errors"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"math"
	"testing"
)

func TestGoMRuby(t *testing.T) { TestingT(t) }

type F struct {
	m *MRuby
	c *LoadContext
}

var _ = Suite(&F{})

func (f *F) SetUpTest(c *C) {
	f.m = New()
	f.c = f.m.NewLoadContext("test.rb")
}

func (f *F) TearDownTest(c *C) {
	f.c.Delete()
	f.m.Delete()
}

func (f *F) TestLoad(c *C) {
	must := func(res interface{}, err error) interface{} {
		c.Check(err, IsNil)
		return res
	}

	c.Check(must(f.c.Load(`nil`)), Equals, nil)
	c.Check(must(f.c.Load(`true`)), Equals, true)
	c.Check(must(f.c.Load(`false`)), Equals, false)
	c.Check(must(f.c.Load(`1 + 1`)), Equals, 2)
	c.Check(must(f.c.Load(`1 - 1`)), Equals, 0)
	c.Check(must(f.c.Load(`1 - 2`)), Equals, -1)
	c.Check(must(f.c.Load(`3.14 + 42`)), Equals, 45.14)

	// assumes typedef int32_t mrb_int
	c.Check(must(f.c.Load(fmt.Sprintf("%d", int16(math.MaxInt16)))), Equals, math.MaxInt16)
	c.Check(must(f.c.Load(fmt.Sprintf("%d", uint16(math.MaxUint16)))), Equals, math.MaxUint16)
	c.Check(must(f.c.Load(fmt.Sprintf("%d", int32(math.MaxInt32)))), Equals, math.MaxInt32)
	c.Check(must(f.c.Load(fmt.Sprintf("%d", uint32(math.MaxUint32)))), Equals, float64(math.MaxUint32))
	// c.Check(must(f.c.Load(fmt.Sprintf("%d", int64(math.MaxInt64)))), Equals, float64(math.MaxInt64)) // FIXME
	c.Check(must(f.c.Load(fmt.Sprintf("%d", uint64(math.MaxUint64)))), Equals, float64(math.MaxUint64))

	c.Check(must(f.c.Load(`domain = "express" + "42" + ".com"`)), Equals, "express42.com")
	c.Check(must(f.c.Load(`domain`)), Equals, "express42.com")
	c.Check(must(f.c.Load(`""`)), Equals, "")

	slice := []interface{}{1, 3.14, "foo"}
	hash := map[interface{}]interface{}{"foo": 1, 3.14: "bar"}
	mix := []interface{}{42, map[interface{}]interface{}{3.14: []interface{}{"bar"}, "foo": 1}}
	c.Check(must(f.c.Load(`[1, 3.14, "foo"]`)), DeepEquals, slice)
	c.Check(must(f.c.Load(`{3.14=>"bar", "foo"=>1}`)), DeepEquals, hash)
	c.Check(must(f.c.Load(`[42, {3.14=>["bar"], "foo"=>1}]`)), DeepEquals, mix)

	c.Check(must(f.c.Load("ARGV.inspect")), Equals, `[]`)
	c.Check(must(f.c.Load("ARGV.inspect", nil, true, false)), Equals, `[nil, true, false]`)
	c.Check(must(f.c.Load("ARGV.inspect", 1, 3.14, "foo")), Equals, `[1, 3.14, "foo"]`)
	c.Check(must(f.c.Load("ARGV.inspect", slice)), Equals, `[[1, 3.14, "foo"]]`)
	c.Check(must(f.c.Load("ARGV.inspect", hash)), Equals, `[{3.14=>"bar", "foo"=>1}]`)
	c.Check(must(f.c.Load("ARGV.inspect", mix)), Equals, `[[42, {3.14=>["bar"], "foo"=>1}]]`)

	res, err := f.c.Load(`foo`)
	c.Check(res, Equals, nil)
	c.Check(err, DeepEquals, errors.New(`test.rb:1: undefined method 'foo' for main (NoMethodError)`))

	res, err = f.c.Load(`
begin
  foo
rescue => e
  e.inspect
  e.inspect
`)
	c.Check(res, Equals, nil)
	c.Check(err, DeepEquals, errors.New("SyntaxError: syntax error"))
}

func (f *F) TestLoadMore(c *C) {
	res, err := f.c.Load("ARGV.map { |x| x * x }", 1, 2, 3)
	c.Check(err, IsNil)
	c.Check(res, DeepEquals, []interface{}{1, 4, 9})
}

func (f *F) TestDefineGoConst(c *C) {
	f.m.DefineGoConst("MY_CONST", 42)
	res, err := f.c.Load("Go::MY_CONST")
	c.Check(err, IsNil)
	c.Check(res, Equals, 42)
}

func (f *F) TestLoadContext(c *C) {
	res, err := f.c.Load(`$global = 1; local = 2`)
	c.Check(res, Equals, 2)
	c.Check(err, IsNil)

	// state is preserved
	res, err = f.c.Load(`$global`)
	c.Check(res, Equals, 1)
	c.Check(err, IsNil)
	res, err = f.c.Load(`local`)
	c.Check(res, Equals, 2)
	c.Check(err, IsNil)

	c2 := f.m.NewLoadContext("test2.rb")
	defer c2.Delete()

	// global variable is accessible from other context
	res, err = c2.Load(`$global`)
	c.Check(res, Equals, 1)
	c.Check(err, IsNil)

	// local is not
	res, err = c2.Load(`local`)
	c.Check(res, Equals, nil)
	c.Check(err, DeepEquals, errors.New("test2.rb:1: undefined method 'local' for main (NoMethodError)"))
}

func (f *F) TestBugRb(c *C) {
	b, err := ioutil.ReadFile("test.rb")
	c.Assert(err, IsNil)
	res, err := f.c.Load(string(b))
	c.Logf("%#v", res)
	c.Check(err, IsNil)
}
