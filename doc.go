// Package gomruby embeds mruby (mini Ruby) VM into Go.
//
// Type conversions:
//   * nil, true, false, string - as expected
//   * Float <-> FIXME tricky
//   * Fixnum <-> FIXME tricky
//   * Symbol <-> gomruby.Symbol
//   * Array <-> []interface{}
//   * Hash <-> map[interface{}]interface{}
package gomruby
