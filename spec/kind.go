package spec

//go:generate stringer -type Kind kind.go
type Kind uint8

const (
	Invalid Kind = iota

	predeclaredTypesBeg
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	Byte
	Rune
	predeclaredTypesEnd

	Error

	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	Struct

	Time
	Duration
)

func (k Kind) IsPredeclared() bool {
	return k > predeclaredTypesBeg && k < predeclaredTypesEnd
}
