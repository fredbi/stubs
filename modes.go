package stubs

// StubMode for generating data
type StubMode uint64

// Has returns true when this mode has the provided flag configured
func (s StubMode) Has(m StubMode) bool {
	return s&m != 0
}

const (
	// Invalid produces a stub which is invalid for a random validation
	Invalid StubMode = 1 << iota
	// InvalidRequired produces a stub which is invalid for required
	InvalidRequired
	// InvalidMaximum produces a stub which is invalid for maximum
	InvalidMaximum
	// InvalidMinimum produces a stub which is invalid for minimum
	InvalidMinimum
	// InvalidMaxLength produces a stub which is invalid for max length
	InvalidMaxLength
	// InvalidMinLength produces a stub which is invalid for min length
	InvalidMinLength
	// InvalidPattern produces a stub which is invalid for pattern
	InvalidPattern
	// InvalidMaxItems produces a stub which is invalid for max items
	InvalidMaxItems
	// InvalidMinItems produces a stub which is invalid for min items
	InvalidMinItems
	// InvalidUniqueItems produces a stub which is invalid for unique items
	InvalidUniqueItems
	// InvalidMultipleOf produces a stub which is invalid for multiple of
	InvalidMultipleOf
	// InvalidEnum produces a stub which is invalid for enum
	InvalidEnum

	// Valid is the default value and generates valid data
	Valid StubMode = 0
)
