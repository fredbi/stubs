package stubs

import (
	"log"
	"testing"
)

func TestGenerators_Structs(t *testing.T) {
	o := simpleOpts{}
	log.Printf("%#v", o)
}
