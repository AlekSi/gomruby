#include <stdlib.h>

#include <mruby.h>
#include <mruby/compile.h>
#include <mruby/string.h>
#include <mruby/array.h>
#include <mruby/hash.h>

#if defined(MRB_NAN_BOXING) || defined(MRB_WORD_BOXING)
#   error "gomruby: boxing is not supported."
#endif

#if defined(MRB_INT16)
#   error "gomruby: code assumes mrb_int == int32_t or mrb_int == int64_t"
#endif

// mrb_fixnum is define, can't access it with cgo
inline static int _gomruby_fixnum(mrb_value o) {
	return mrb_fixnum(o);
}

// mrb_float is define, can't access it with cgo
inline static double _gomruby_float(mrb_value o) {
	return mrb_float(o);
}

// mrb_nil_p is define, can't access it with cgo
inline static int _gomruby_is_nil(mrb_value o) {
	return mrb_nil_p(o);
}
