package gomruby

/*
#include <gomruby.h>
#cgo CFLAGS: -I. -Imruby/include
#cgo LDFLAGS: libmruby.a -lm
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	inspectCS, argvCS *C.char
)

func init() {
	inspectCS = C.CString("inspect")
	argvCS = C.CString("ARGV")
}

// mruby VM.
type MRuby struct {
	state *C.mrb_state
}

// Creates new mruby VM. Panics if it's not possible.
func New() *MRuby {
	state := C.mrb_open()
	if state == nil {
		panic(errors.New("gomruby bug: failed to create mruby state"))
	}
	return &MRuby{state}
}

// Deletes mruby VM.
func (m *MRuby) Delete() {
	if m.state != nil {
		C.mrb_close(m.state)
		m.state = nil
	}
}

// Converts Go value to mruby value.
func (m *MRuby) mrubyValue(i interface{}) C.mrb_value {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Invalid:
		return C.mrb_nil_value() // mrb_undef_value() explodes
	case reflect.Bool:
		b := v.Bool()
		if b {
			return C.mrb_true_value()
		} else {
			return C.mrb_false_value()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return C.mrb_fixnum_value(C.mrb_int(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return C.mrb_fixnum_value(C.mrb_int(v.Uint()))
	case reflect.Float32, reflect.Float64:
		return C.mrb_float_value((*C.struct_mrb_state)(m.state), C.mrb_float(v.Float()))
	case reflect.String:
		cs := C.CString(v.String())
		defer C.free(unsafe.Pointer(cs))
		return C.mrb_str_new_cstr(m.state, cs)
	case reflect.Array, reflect.Slice:
		l := v.Len()
		res := C.mrb_ary_new_capa(m.state, C.mrb_int(l))
		for i := 0; i < l; i++ {
			C.mrb_ary_set(m.state, res, C.mrb_int(i), m.mrubyValue(v.Index(i).Interface()))
		}
		return res
	case reflect.Map:
		l := v.Len()
		res := C.mrb_hash_new_capa(m.state, C.int(l))
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			C.mrb_hash_set(m.state, res, m.mrubyValue(key.Interface()), m.mrubyValue(val.Interface()))
		}
		return res
	}

	panic(fmt.Errorf("gomruby bug: failed to convert Go value %#v (%T) to mruby value", i, i))
}

// Converts mruby value to Go value.
func (m *MRuby) goValue(v C.mrb_value) interface{} {
	switch v.tt {
	case C.MRB_TT_UNDEF: // for example, result of syntax error
		return nil
	case C.MRB_TT_TRUE:
		return true
	case C.MRB_TT_FALSE:
		if C._gomruby_is_nil(v) == 0 {
			return false
		}
		return nil
	case C.MRB_TT_FIXNUM:
		return int(C._gomruby_fixnum(v))
	case C.MRB_TT_FLOAT:
		return float64(C._gomruby_float(v))
	case C.MRB_TT_STRING:
		cs := C.mrb_string_value_ptr(m.state, v)
		return C.GoString(cs)
	case C.MRB_TT_ARRAY:
		l := int(C.mrb_ary_len(m.state, v))
		res := make([]interface{}, l)
		for i := 0; i < l; i++ {
			res[i] = m.goValue(C.mrb_ary_ref(m.state, v, C.mrb_int(i)))
		}
		return res
	case C.MRB_TT_HASH:
		keys := C.mrb_hash_keys(m.state, v)
		l := int(C.mrb_ary_len(m.state, keys))
		res := make(map[interface{}]interface{}, l)
		for i := 0; i < l; i++ {
			key := C.mrb_ary_ref(m.state, keys, C.mrb_int(i))
			val := C.mrb_hash_get(m.state, v, key)
			res[m.goValue(key)] = m.goValue(val)
		}
		return res
	}

	panic(fmt.Errorf("gomruby bug: failed to convert mruby value %#v to Go value", v))
}

// v.inspect()
func (m *MRuby) inspect(v C.mrb_value) string {
	v = C.mrb_funcall_argv(m.state, v, C.mrb_intern_cstr(m.state, inspectCS), 0, nil)
	cs := C.mrb_string_value_ptr(m.state, v)
	return C.GoString(cs)
}

// mruby VM load context.
type LoadContext struct {
	context *C.mrbc_context
	m       *MRuby
}

// Creates new load context. Panics if it's not possible. Filename is used in error messages.
func (m *MRuby) NewLoadContext(filename string) (context *LoadContext) {
	ctx := C.mrbc_context_new(m.state)
	if ctx == nil {
		panic(errors.New("gomruby bug: failed to create mruby context"))
	}

	if filename != "" {
		fn := C.CString(filename)
		C.mrbc_filename(m.state, ctx, fn)
		C.free(unsafe.Pointer(fn))
	}

	context = &LoadContext{ctx, m}
	return
}

// Deletes load context.
func (c *LoadContext) Delete() {
	if c.context != nil {
		C.mrbc_context_free(c.m.state, c.context)
		c.m = nil
		c.context = nil
	}
}

// Loads mruby code. Arguments are exposed as ARGV array.
func (c *LoadContext) Load(code string, args ...interface{}) (res interface{}, err error) {
	l := len(args)
	ARGV := C.mrb_ary_new_capa(c.m.state, C.mrb_int(l))
	for i := 0; i < l; i++ {
		ii := C.mrb_int(i)
		C.mrb_ary_set(c.m.state, ARGV, ii, c.m.mrubyValue(args[ii]))
	}
	C.mrb_define_global_const(c.m.state, argvCS, ARGV)

	codeC := C.CString(code)
	defer C.free(unsafe.Pointer(codeC))
	v := C.mrb_load_string_cxt(c.m.state, codeC, c.context)
	res = c.m.goValue(v)
	if c.m.state.exc != nil {
		v = C.mrb_obj_value(unsafe.Pointer(c.m.state.exc))
		err = errors.New(c.m.inspect(v))
	}
	return
}
