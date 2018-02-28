package stubs

import (
	"github.com/mitchellh/mapstructure"
)

const (
	// XdataGen is the extension tag to be used in specs and hint the stub generator
	XdataGen = "x-datagen"
	// StubsDefaultWordCount is the default maximum word count in sentences. It may be overriden with the "words" argument.
	StubsDefaultWordCount = 10
	// StubsDefaultStringLength is the default length for generated strings
	StubsDefaultStringLength = 30
	// StubsDefaultSupplemental defines if "supplemental" words from extra dictionary in faker is enabled
	StubsDefaultSupplemental = true
	// StubsDefaultMaxAmount defines the default max amount for prices and other currency amounts
	StubsDefaultMaxAmount = 1000000
	// StubsDefaultMinAmount defines the default min amount for prices and other currency amounts
	StubsDefaultMinAmount = 100
	// StubsDefaultMaxSmallAmount defines the default max amount for small prices and other small currency amounts
	StubsDefaultMaxSmallAmount = 1000
	// StubsDefaultMinSmallAmount defines the default min amount for small prices and other small currency amounts
	StubsDefaultMinSmallAmount = 10
)

// basicGeneratorOpts publishes basic operations for GeneratorOpts
type basicGeneratorOpts interface {
	// value generator name
	Name() string

	// Args for the value generator (eg. number of words in a sentence)
	// TODO(fredbi): the following may not be true now
	// Arguments here are used to generate a valid value when no validations are specified.
	// Args are used as default setting but validations can override the args should that be necessary.
	Args() *argTags

	// Mode which kind of random data to return and to indicate which validation(s) should fail.
	// This is a bitmask so it allows for combinations of invalid values.
	Mode() StubMode

	// Infer deduces the generator to be used from inference rules
	Infer()

	// ExtOverride captures the arguments under the x-datagen extension
	ExtOverride(map[string]interface{}) error
	// SetArgs overrides Args
	SetArgs(*argTags)
}

// GeneratorOpts interface to capture various types that can get data generated for them.
type GeneratorOpts interface {
	seeder
	basicGeneratorOpts

	// FieldName for the value generator, this is mostly used as an alternative to the name
	// for inferring which value generator to use
	// TODO(fredbi): atm unused
	FieldName() string

	// Type for the value generator to return, aids in infering the name of the value generator
	Type() string

	// Format for the value generator to return, aids in infering the name of the value generator
	Format() string

	// Maximum a numeric value can have, returns value, exclusive, defined
	Maximum() (float64, bool, bool)

	// Minimum a numeric value can have, returns value, exclusive, defined
	Minimum() (float64, bool, bool)

	// MaxLength a string can have, returns value, defined
	MaxLength() (int64, bool)

	// MinLength a string can have, returns value, defined
	MinLength() (int64, bool)

	// Pattern a string should match, returns value, defined
	Pattern() (string, bool)

	// MaxItems a collection of values can contain, returns length, defined
	MaxItems() (int64, bool)

	// MinItems a collection of values must contain, returns length, defined
	MinItems() (int64, bool)

	// UniqueItems when true the collection can't contain duplicates
	UniqueItems() bool

	// MultipleOf a numeric value should be divisible by this value, returns value, defined
	MultipleOf() (float64, bool)

	// Enum a list of acceptable values for a value, returns value, defined
	Enum() ([]interface{}, bool)

	// Items options for the members of a collection
	Items() (GeneratorOpts, error)

	// Required when true the property can't be nil
	Required() bool
}

type boolOrMap struct {
	Enabled bool
	Args    map[string]string
}

type boolOrSlice struct {
	Enabled bool
	Args    []string
}

func (b *boolOrSlice) Contains(arg string) bool {
	for _, v := range b.Args {
		if v == arg {
			return true
		}
	}
	return false
}

// argTags describe all override options that may be set in generation options
type argTags struct {
	Valid                   boolOrMap   `mapstructure:"valid"`
	Invalid                 boolOrMap   `mapstructure:"invalid"`
	Lang                    []string    `mapstructure:"lang"`
	Length                  int         `mapstructure:"length"`
	Words                   int         `mapstructure:"words"`
	Supplemental            bool        `mapstructure:"supplemental"`
	Max                     *float64    `mapstructure:"max"`
	Min                     *float64    `mapstructure:"min"`
	MultipleOf              *float64    `mapstructure:"multipleOf"`
	Precision               *int64      `mapstructure:"precision"`
	WithExample             boolOrMap   `mapstructure:"withExample"`
	WithDefault             boolOrMap   `mapstructure:"withDefault"`
	WithEdgeCase            boolOrSlice `mapstructure:"withEdgeCase"`
	WithTilting             boolOrMap   `mapstructure:"withTilting"`
	WithAllValidationChecks boolOrMap   `mapstructure:"withAllValidationChecks"`
	SkipTilting             bool        `mapstructure:"skipTilting"`
	SkipFuzzying            bool        `mapstructure:"skipFuzzying"`
}

// genTag describe the structure of a x-datagen hint in the swagger spec
type genTag struct {
	Name string   `mapstructure:"name"` // name refers to the stub generator's key
	Args argTags  `mapstructure:"args"` // args define specific behavior expected from the generator
	Mode StubMode `mapstructure:"mode"` // mode selects the stub generation strategy
}

// genOpts is the concrete type for basic common generation options
type genOpts struct {
	name  string   // name refers to the stub generator's key
	args  argTags  // args define specific behavior expected from the generator
	mode  StubMode // mode selects the stub generation strategy
	rules []ruler  // rules is the ordered list of rules used to infer the generator options
}

func (g *genOpts) Name() string {
	return g.name
}

func (g *genOpts) Args() *argTags {
	return &g.args
}

func (g *genOpts) Mode() StubMode {
	return g.mode
}

func (g *genOpts) SetArgs(args *argTags) {
	g.args = *args
}

// Infer chains inference rules to take a decision about the generator option to set
func (g *genOpts) Infer() {
	var decisions = make([]basicGeneratorOpts, 0, 50)
	for _, rule := range g.rules {
		decision := rule.Decide()
		if decision != nil {
			debugLog("decision: %v", decision)
			decisions = append(decisions, rule.Decide())
		}
	}
	g.merge(decisions)
}

// merge merges decisions from several rulers
func (g *genOpts) merge(decisions []basicGeneratorOpts) {
	// Merge decisions from inference rules
	for _, d := range decisions {
		debugLog("merging decision: %v", d)
		if d.Name() != "" {
			g.name = d.Name()
		}
		// other args...
	}
	return
}

// ExtOverride captures the arguments under the x-datagen extension
func (g *genOpts) ExtOverride(extensions map[string]interface{}) error {
	if ext, ok := extensions[XdataGen]; ok {
		debugLog("found extension")
		// TODO: factorize in loadOpts()
		// whenever a x-datagen extension is specified, take its directives for granted
		tag := genTag{}
		if err := mapstructure.WeakDecode(ext, &tag); err != nil {
			return err
		}
		g.name = tag.Name
		g.args = tag.Args
		g.mode = tag.Mode
		return nil
	}
	return nil
}
