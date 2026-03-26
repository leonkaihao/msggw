package model

// Symbol record type and value
// Type  	|  	Value
// ---------------------
// raw   	|  	string
// metadata |  	msg metadata field string
// prop     |   config prop field string
// func     | 	builtin Func interface
// topic    |   no value
// operator |   builtin Operator interface
// keyword  |   KW_.. const
// mix      |   no value, but has children Symbols
type Symbol interface {
	Type() string
	Value() interface{}
	Children() []Symbol
	Format(ctx MsgContext) (string, error)
}
