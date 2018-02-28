package stubs

import conv "github.com/cstockton/go-conv"

// Converter provides functions to convert values from other types
type Converter interface {
	Bool(interface{}) (bool, error)
	Int(interface{}) (int, error)
}

type defaultConverter struct {
}

func (d defaultConverter) Bool(from interface{}) (bool, error) {
	return conv.Bool(from)
}

func (d defaultConverter) Int(from interface{}) (int, error) {
	return conv.Int(from)
}
