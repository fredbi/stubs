package stubs

import (
	"fmt"
	"math"
	"regexp"
	"regexp/syntax"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/asaskevich/govalidator"
	"github.com/davecgh/go-spew/spew"
	"github.com/manveru/faker"
	regen "github.com/zach-klippenstein/goregen"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ValueGenerator represents a function to generate a piece of random data
type ValueGenerator func(GeneratorOpts) (interface{}, error)

type generators struct {
	faker     *faker.Faker
	conv      Converter
	gens      map[string]ValueGenerator
	regenArgs *regen.GeneratorArgs
	entropy   *randomGenerator
}

var (
	// ErrNoValid indicates to the caller that value generator could not abide by validation constraints and generate a valid value
	ErrNoValid error
	// ErrNoInvalid indicates to the caller that value generator could not abide by specified invalid mode flags
	ErrNoInvalid error
)

func init() {
	ErrNoValid = fmt.Errorf("No valid value could be generated")
	ErrNoInvalid = fmt.Errorf("No invalid value could be generated")
}

// newgenerator instantiate a new generator for a specific language.
//
// By default, language is english ("en").
// TODO: replace lang by GeneralOps
func newGenerator(lang string) (*generators, error) {
	if lang == "" {
		lang = "en"
	}
	debugLog("new generator for lang=%s", lang)
	faker, err := faker.New(lang)
	if err != nil {
		return nil, err
	}
	// TODO: set seeds on underlying fakers
	g := &generators{
		faker: faker,
		conv:  defaultConverter{},
		regenArgs: &regen.GeneratorArgs{
			Flags: syntax.Perl,
		},
		entropy: newRandomGenerator(randomOpts{
			// TODO: from opts
			AutoSeed: true,
		},
		),
	}
	g.makeGenerators()
	return g, nil
}

// makeGenerators initializes the map of supported stub generators
func (g *generators) makeGenerators() {
	g.gens = map[string]ValueGenerator{
		"amount":            g.genAmount,
		"small-amount":      g.genSmallAmount,
		"adjective":         g.string(randomdata.Noun),
		"bool":              g.bool,
		"characters":        g.intString(g.faker.Characters),
		"city":              g.string(g.faker.City),
		"city-prefix":       g.string(g.faker.CityPrefix),
		"city-suffix":       g.string(g.faker.CitySuffix),
		"company":           g.string(g.faker.CompanyName),
		"company-bs":        g.string(g.faker.CompanyBs),
		"company-slogan":    g.string(g.faker.CompanyCatchPhrase),
		"company-suffix":    g.string(g.faker.CompanySuffix),
		"country":           g.string(g.faker.Country),
		"credit-card":       g.fromPattern(govalidator.CreditCard),
		"domain":            g.string(g.faker.DomainName),
		"domain-suffix":     g.string(g.faker.DomainSuffix),
		"double":            g.numGenFloat64,
		"email":             g.string(g.faker.Email),
		"first-name":        g.string(g.faker.FirstName),
		"float":             g.numGenFloat32,
		"free-email":        g.string(g.faker.FreeEmail),
		"hexcolor":          g.fromPattern(govalidator.Hexcolor),
		"hostname":          g.string(g.faker.DomainWord),
		"int32":             g.numGenInt32,
		"int64":             g.numGenInt64,
		"ip":                g.altws(randomdata.IpV4Address, randomdata.IpV6Address),
		"ipv4":              g.string(randomdata.IpV4Address),
		"ipv6":              g.string(randomdata.IpV6Address),
		"isbn":              g.altwsp(govalidator.ISBN10, govalidator.ISBN13),
		"isbn10":            g.fromPattern(govalidator.ISBN10),
		"isbn13":            g.fromPattern(govalidator.ISBN13),
		"job-title":         g.string(g.faker.JobTitle),
		"landline":          g.string(g.faker.PhoneNumber),
		"last-name":         g.string(g.faker.LastName),
		"latitude":          g.float(g.faker.Latitude),
		"longitude":         g.float(g.faker.Longitude),
		"mac-address":       g.fromPattern("^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$"),
		"mobile":            g.string(g.faker.CellPhoneNumber),
		"name":              g.string(g.faker.Name),
		"name-prefix":       g.string(g.faker.NamePrefix),
		"name-suffix":       g.string(g.faker.NameSuffix),
		"noun":              g.string(randomdata.Noun),
		"number":            g.numGenFloat64,
		"paragraph":         g.intBoolString(g.faker.Paragraph),
		"paragraphs":        g.intBoolStrings(g.faker.Paragraphs),
		"postcode":          g.string(g.faker.PostCode),
		"pattern":           g.fromPattern(""),
		"rgbcolor":          g.fromPattern(govalidator.RGBcolor),
		"safe-email":        g.string(g.faker.SafeEmail),
		"secondary-address": g.string(g.faker.SecondaryAddress),
		"sentence":          g.intBoolString(g.faker.Sentence),
		"sentences":         g.intBoolStrings(g.faker.Sentences),
		"silly-name":        g.string(randomdata.SillyName),
		"ssn":               g.fromPattern(govalidator.SSN),
		"state":             g.string(g.faker.StateAbbr),
		"state-name":        g.string(g.faker.State),
		"street-address":    g.string(g.faker.StreetAddress),
		"street-name":       g.string(g.faker.StreetName),
		"street-suffix":     g.string(g.faker.StreetSuffix),
		"uint32":            g.numGenUint32,
		"uint64":            g.numGenUint64,
		"user-name":         g.string(g.faker.UserName),
		"uuid":              g.fromPattern(strfmt.UUIDPattern),
		"uuid3":             g.fromPattern(strfmt.UUID3Pattern),
		"uuid4":             g.fromPattern(strfmt.UUID4Pattern),
		"uuid5":             g.fromPattern(strfmt.UUID5Pattern),
		"word":              g.string(func() string { return g.faker.Words(1, false)[0] }),
		"words":             g.intBoolStrings(g.faker.Words),
		"date":              g.dateGen,
		"datetime":          g.dateTimeGen,
		"duration":          g.durationGen,
	}

	/* TODO:
	* add near date
	* add near date-time, timestamps, update-time...
	* add short duration
	* add slices
	 */
}

func normalizeGeneratorName(str string) string {
	kn := strings.ToLower(str)
	if k, ok := generatorAliases[kn]; ok {
		return k
	}
	return kn
}

// For selects an appropriate value generator for the generation option.
func (g *generators) For(opts GeneratorOpts) (ValueGenerator, bool) {
	debugLog("looking for valueGenerator for option: %s", opts.Name())
	if gen, ok := g.gens[normalizeGeneratorName(opts.Name())]; ok {
		return gen, true
	}
	if gen, ok := g.gens[normalizeGeneratorName(swag.ToCommandName(opts.FieldName()))]; ok {
		return gen, true
	}
	return nil, false
}

// altws returns a values generator which chooses randomly among a list of generating functions
//
// Example:
//  altws(randomdata.IpV4Address, randomdata.IpV6Address)
func (g *generators) altws(fns ...func() string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		idx := g.entropy.IntN(len(fns))
		return fns[idx](), nil
	}
}

// altwsp returns a values generator which chooses randomly among a list of generating patterns.
// Patterns are generated using regen.Generate().
//
// Example:
//  altwsp(govalidator.ISBN10, govalidator.ISBN13)
func (g *generators) altwsp(patterns ...string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		idx := g.entropy.IntN(len(patterns))
		return regen.Generate(patterns[idx])
	}
}

// fromPattern returns a value generator based on pattern generation.
// Patterns are generated using regen.Generate().
func (g *generators) fromPattern(pattern string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		if pattern == "" {
			// use option-defined pattern
			p, defined := opts.Pattern()
			if defined {
				pattern = p
			}
		}
		// check pattern is valid
		if _, err := regexp.Compile(pattern); err != nil {
			return nil, err
		}
		debugLog("calling regen with pattern: %q", pattern)
		rexgen, err := regen.NewGenerator(pattern, g.regenArgs)
		if err != nil {
			return nil, err
		}
		return rexgen.Generate(), nil
	}
}

func (g *generators) stringError(fn func() (string, error)) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn()
	}
}

func (g *generators) string(fn func() string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) stringer(fn func() fmt.Stringer) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn().String(), nil
	}
}

func (g *generators) float(fn func() float64) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) integer(fn func() int64) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		return fn(), nil
	}
}

func (g *generators) intString(fn func(int) string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		var count int = StubsDefaultStringLength
		if opts != nil {
			args := opts.Args()
			if args.Length != 0 {
				count = args.Length
			}
		}
		return fn(count), nil
	}
}

func (g *generators) intBoolString(fn func(int, bool) string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		var count int = StubsDefaultWordCount
		var supplemental bool = StubsDefaultSupplemental
		if opts != nil {
			args := opts.Args()
			if args.Words != 0 {
				count = args.Words
			}
			supplemental = args.Supplemental
		}
		return fn(count, supplemental), nil
	}
}

func (g *generators) intBoolStrings(fn func(int, bool) []string) ValueGenerator {
	return func(opts GeneratorOpts) (interface{}, error) {
		var count int = StubsDefaultWordCount
		var supplemental bool = StubsDefaultSupplemental
		if opts != nil {
			args := opts.Args()
			if args.Words != 0 {
				count = args.Words
			}
			supplemental = args.Supplemental
		}
		return fn(count, supplemental), nil
	}
}

func (g *generators) bool(opts GeneratorOpts) (interface{}, error) {
	return g.entropy.Bool(), nil
}

func (g *generators) genAmount(opts GeneratorOpts) (interface{}, error) {
	args := opts.Args()
	if args.Max == nil {
		args.Max = swag.Float64(StubsDefaultMaxAmount)
	}
	if args.Min == nil || swag.Float64Value(args.Min) < 0 {
		args.Min = swag.Float64(StubsDefaultMinAmount)
	}
	opts.SetArgs(args)
	spew.Dump(opts)
	res, err := g.numGenFloat64(opts)
	spew.Dump(res)
	if err != nil {
		return 0, err
	}
	f := res.(float64)
	return math.Trunc(f*float64(100)) / float64(100), nil
}

func (g *generators) genSmallAmount(opts GeneratorOpts) (interface{}, error) {
	args := opts.Args()
	if args.Max == nil {
		args.Max = swag.Float64(StubsDefaultMaxSmallAmount)
	}
	if args.Min == nil || swag.Float64Value(args.Min) < 0 {
		args.Min = swag.Float64(StubsDefaultMinSmallAmount)
	}
	opts.SetArgs(args)
	res, err := g.numGenFloat64(opts)
	if err != nil {
		return 0, err
	}
	f := res.(float64)
	return math.Trunc(f*float64(100)) / float64(100), nil
}

func (g *generators) numGenFloat64(opts GeneratorOpts) (interface{}, error) {
	var min, max, multipleOf float64 = defaultMinFloat64, defaultMaxFloat64, float64(0)
	var mode StubMode
	returnValid := true
	hasMultiple := false
	if opts != nil {
		// check options
		args := opts.Args()
		if args.Min != nil {
			min = swag.Float64Value(args.Min)
		}
		if args.Max != nil {
			max = swag.Float64Value(args.Max)
		}
		if args.MultipleOf != nil {
			multipleOf = swag.Float64Value(args.MultipleOf)
			hasMultiple = true
		}

		// check validations
		minCheck, _, definedMin := opts.Minimum()
		if definedMin && minCheck > min {
			min = minCheck
		}
		maxCheck, _, definedMax := opts.Maximum()
		if definedMax && maxCheck < max {
			max = maxCheck
		}

		multipleOfCheck, definedMult := opts.MultipleOf()
		if definedMult {
			multipleOf = multipleOfCheck
			hasMultiple = true
			if multipleOf != 0 && multipleOf != 1 {
				min = min / multipleOf
				max = max / multipleOf
			}
		}

		mode = opts.Mode()

		if mode.Has(InvalidMinimum) {
			returnValid = false
		}

		// clear edge cases
		if mode.Has(InvalidMinimum) && !definedMin { // Safeguard
			return 0, ErrNoInvalid
		} else if mode.Has(InvalidMaximum) && !definedMax { // Safeguard
			return 0, ErrNoInvalid
		} else if mode.Has(InvalidMultipleOf) && !definedMult { // Safeguard
			return 0, ErrNoInvalid

			// mutually exclusive failures
		} else if mode.Has(InvalidMinimum) {
			returnValid = false
			// TODO: check if boundaries included or not
			max = minCheck
			if hasMultiple && multipleOf != 0 {
				max = max / multipleOf
			}
			if max < min {
				return 0, ErrNoInvalid
			}
		} else if mode.Has(InvalidMaximum) {
			returnValid = false
			min = maxCheck
			if hasMultiple && multipleOf != 0 {
				min = min / multipleOf
			}
			if max < min {
				return 0, ErrNoInvalid
			}
		}
	}

	// return valid
	if returnValid {
		if hasMultiple && multipleOf == 0 {
			// only valid value
			return 0, nil
		}
		if hasMultiple {
			return multipleOf * math.Trunc(g.entropy.Float64(min, max)), nil
		}
		return g.entropy.Float64(min, max), nil
	}

	// return invalid
	res := g.entropy.Float64(min, max)

	if mode.Has(InvalidMultipleOf) {
		// loop until we get a validation error. This should not take long...
		for validate.MultipleOfNativeType("", "", res, multipleOf) == nil {
			res = g.entropy.Float64(min, max)
		}
	}
	return res, nil
}

func (g *generators) numGenFloat32(opts GeneratorOpts) (interface{}, error) {
	var min, max, multipleOf float32 = defaultMinFloat32, defaultMaxFloat32, float32(0)
	var mode StubMode
	returnValid := true
	hasMultiple := false
	if opts != nil {
		// check options
		args := opts.Args()
		if args.Min != nil {
			min = float32(swag.Float64Value(args.Min))
		}
		if args.Max != nil {
			max = float32(swag.Float64Value(args.Max))
		}
		if args.MultipleOf != nil {
			multipleOf = float32(swag.Float64Value(args.MultipleOf))
			hasMultiple = true
		}

		// check validations
		minCheck, _, definedMin := opts.Minimum()
		if definedMin && minCheck > float64(min) {
			min = float32(minCheck)
		}
		maxCheck, _, definedMax := opts.Maximum()
		if definedMax && maxCheck < float64(max) {
			max = float32(maxCheck)
		}

		multipleOfCheck, definedMult := opts.MultipleOf()
		if definedMult {
			multipleOf = float32(multipleOfCheck)
			hasMultiple = true
			if multipleOf != 0 && multipleOf != 1 {
				min = min / multipleOf
				max = max / multipleOf
			}
		}

		mode = opts.Mode()

		if mode.Has(InvalidMinimum) {
			returnValid = false
		}

		// clear edge cases
		if mode.Has(InvalidMinimum) && !definedMin { // Safeguard
			return 0, ErrNoInvalid
		} else if mode.Has(InvalidMaximum) && !definedMax { // Safeguard
			return 0, ErrNoInvalid
		} else if mode.Has(InvalidMultipleOf) && !definedMult { // Safeguard
			return 0, ErrNoInvalid

			// mutually exclusive failures
		} else if mode.Has(InvalidMinimum) {
			returnValid = false
			max = float32(minCheck)
			if hasMultiple && multipleOf != 0 {
				max = max / multipleOf
			}
			if max < min {
				return 0, ErrNoInvalid
			}
		} else if mode.Has(InvalidMaximum) {
			returnValid = false
			min = float32(maxCheck)
			if hasMultiple && multipleOf != 0 {
				min = min / multipleOf
			}
			if max < min {
				return 0, ErrNoInvalid
			}
		}
	}

	// return valid
	if returnValid {
		if hasMultiple && multipleOf == 0 {
			// only valid value
			return 0, nil
		}
		if hasMultiple {
			return multipleOf * float32(math.Trunc(float64(g.entropy.Float32(min, max)))), nil
		}
		return g.entropy.Float32(min, max), nil
	}

	// return invalid
	res := g.entropy.Float32(min, max)

	if mode.Has(InvalidMultipleOf) {
		// loop until we get a validation error. This should not take long...
		for validate.MultipleOfNativeType("", "", res, float64(multipleOf)) == nil {
			res = g.entropy.Float32(min, max)
		}
	}
	return res, nil
}

func (g *generators) numGenInt64(opts GeneratorOpts) (interface{}, error) {
	return g.entropy.Int64(0, 0), nil
}

func (g *generators) numGenInt32(opts GeneratorOpts) (interface{}, error) {
	return g.entropy.Int32(0, 0), nil
}

func (g *generators) numGenUint32(opts GeneratorOpts) (interface{}, error) {
	return g.entropy.Uint32(0, 0), nil
}

func (g *generators) numGenUint64(opts GeneratorOpts) (interface{}, error) {
	return g.entropy.Uint64(0, 0), nil
}

func (g *generators) dateGen(opts GeneratorOpts) (interface{}, error) {
	offset := g.entropy.Int64(-10000, +10000) // TODO: options to qualify dates
	t := time.Now().AddDate(0, 0, int(offset))
	var d strfmt.Date = strfmt.Date(t)
	return d.String(), nil
}

func (g *generators) dateTimeGen(opts GeneratorOpts) (interface{}, error) {
	offset := g.entropy.Int64(-10000, +10000) // TODO: options to qualify dates
	t := time.Now().AddDate(0, 0, int(offset))
	t.Add(time.Duration(offset))
	offset = g.entropy.Int64(-1000000000, +1000000000) // TODO: options to qualify date-times
	var dt strfmt.DateTime = strfmt.DateTime(t)
	return dt.String(), nil
}

func (g *generators) durationGen(opts GeneratorOpts) (interface{}, error) {
	offset := g.entropy.Int64(0, 0)
	return strfmt.Duration(offset).String(), nil
}
