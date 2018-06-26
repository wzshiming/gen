package parser

import (
	"github.com/wzshiming/gen/spec"
	"github.com/wzshiming/gotype"
)

var kindMapping = map[gotype.Kind]spec.Kind{
	gotype.Bool:       spec.Bool,
	gotype.Int:        spec.Int,
	gotype.Int8:       spec.Int8,
	gotype.Int16:      spec.Int16,
	gotype.Int32:      spec.Int32,
	gotype.Int64:      spec.Int64,
	gotype.Uint:       spec.Uint,
	gotype.Uint8:      spec.Uint8,
	gotype.Uint16:     spec.Uint16,
	gotype.Uint32:     spec.Uint32,
	gotype.Uint64:     spec.Uint64,
	gotype.Uintptr:    spec.Uintptr,
	gotype.Float32:    spec.Float32,
	gotype.Float64:    spec.Float64,
	gotype.Complex64:  spec.Complex64,
	gotype.Complex128: spec.Complex128,
	gotype.String:     spec.String,
	gotype.Byte:       spec.Byte,
	gotype.Rune:       spec.Rune,
	gotype.Error:      spec.Error,
	gotype.Array:      spec.Array,
	gotype.Chan:       spec.Chan,
	gotype.Func:       spec.Func,
	gotype.Interface:  spec.Interface,
	gotype.Map:        spec.Map,
	gotype.Ptr:        spec.Ptr,
	gotype.Slice:      spec.Slice,
	gotype.Struct:     spec.Struct,
}
