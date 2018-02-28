package stubs

import (
	"strings"

	"github.com/go-openapi/spec"
)

// ruler knows how to determine stub generation rules in a chain of rules
type ruler interface {
	Decide() basicGeneratorOpts
}

// typeRuler infers options from type
type typeRuler struct {
	ruler
	Type string
}

// newTypeRuler instantiate a new ruler for type-based decisions
func newTypeRulerFor(v interface{}) *typeRuler {
	t := &typeRuler{}
	switch tv := v.(type) {
	case *spec.Parameter:
		t.Type = tv.Type
	case *spec.Header:
		t.Type = tv.Type
	case *spec.Items:
		t.Type = tv.Type
	case *spec.Schema:
		if tv.Type != nil && !tv.Type.Contains("object") && !tv.Type.Contains("array") {
			// TODO: multiple types
			t.Type = tv.Type[0]
		}
	default:
		return nil
	}
	if t.Type == "" {
		return nil
	}
	return t
}

// Decide takes a decision according to Type
func (t *typeRuler) Decide() basicGeneratorOpts {
	debugLog("typeRuler.Decide()")
	if t == nil {
		return nil
	}
	g := &genOpts{}
	switch t.Type {
	case "string":
		g.name = "sentence"
	case "boolean":
		g.name = "bool"
	case "integer":
		g.name = "int64"
	case "number":
		g.name = "float64"
	// TODO: binary
	case "array":
		fallthrough
	default:
		return nil
	}
	debugLog("typeRuler decides: %s", g.name)
	return g
}

// formatRuler infers options from format
type formatRuler struct {
	ruler
	Format string
}

// newFormatRuler instantiate a new ruler for format-based decisions
func newFormatRulerFor(v interface{}) *formatRuler {
	f := &formatRuler{}
	switch tv := v.(type) {
	case *spec.Parameter:
		f.Format = tv.Format
	case *spec.Header:
		f.Format = tv.Format
	case *spec.Items:
		f.Format = tv.Format
	case *spec.Schema:
		f.Format = tv.Format
	default:
		return nil
	}
	if f.Format == "" {
		return nil
	}
	return f
}

// Decide takes a decision according to Format
func (f *formatRuler) Decide() basicGeneratorOpts {
	debugLog("formatRuler.Decide()")
	if f == nil {
		return nil
	}
	debugLog("formatRuler checks for %s given format %s", normalizeGeneratorName(f.Format), f.Format)
	if _, found := generatorAliases[normalizeGeneratorName(f.Format)]; found {
		g := &genOpts{}
		g.name = f.Format
		debugLog("formatRuler decides: %s", g.name)
		return g
	}
	return nil
}

// patternRuler infers options from pattern
type patternRuler struct {
	ruler
	Pattern string
}

// newPatternRuler instantiate a new ruler for pattern-based decisions
func newPatternRulerFor(v interface{}) *patternRuler {
	p := &patternRuler{}
	switch tv := v.(type) {
	case *spec.Parameter:
		p.Pattern = tv.Pattern
	case *spec.Header:
		p.Pattern = tv.Pattern
	case *spec.Items:
		p.Pattern = tv.Pattern
	case *spec.Schema:
		p.Pattern = tv.Pattern
	default:
		return nil
	}
	if p.Pattern == "" {
		return nil
	}
	return p
}

// Decide takes a decision according to Pattern
func (p *patternRuler) Decide() basicGeneratorOpts {
	debugLog("patternRuler.Decide()")
	if p == nil || p.Pattern == "" {
		return nil
	}
	g := &genOpts{}
	g.name = "pattern"
	debugLog("patternRuler decides: %s", g.name)
	return g
}

// fuzzyRuler infers options from title and description
type fuzzyRuler struct {
	ruler
	Title       string
	Description string
}

// newFuzzyRuler instantiate a new ruler for fuzzy-based decisions
// TODO: inference could be helped with default or example values
// TODO: header inference could be helped with the header key in Response
func newFuzzyRulerFor(v interface{}) *fuzzyRuler {
	f := &fuzzyRuler{}
	switch tv := v.(type) {
	case *spec.Parameter:
		f.Title = tv.Name
		f.Description = tv.Description
	case *spec.Header:
		f.Description = tv.Description
	case *spec.Schema:
		f.Title = tv.Title
		f.Description = tv.Description
	default:
		return nil
	}
	if f.Title == "" && f.Description == "" {
		return nil
	}
	return f
}

// Decide takes a decision according to probable matches in Title and Description
func (f *fuzzyRuler) Decide() basicGeneratorOpts {
	debugLog("fuzzyRuler.Decide()")
	if f == nil || (f.Title == "" && f.Description == "") {
		return nil
	}
	scoreTitle := scoreAgainstProposals(f.Title, generatorAliases)
	scoreDescription := scoreAgainstProposals(f.Description, generatorAliases)
	key := mergeAndElectProposal(scoreTitle, scoreDescription)
	if key == "" {
		return nil
	}
	g := &genOpts{}
	g.name = key
	return g
}

// itemsRuler infers options from items type
type itemsRuler struct {
	ruler
	byType    *typeRuler
	byFormat  *formatRuler
	byPattern *patternRuler
	Child     *itemsRuler
	// TODO: keep CollectionFormat to allow for ready-to-parse fixtures (option)
}

// newItemsRuler instantiate a new ruler for items-based decisions
// TODO: inference could be helped with default or example values
func newItemsRulerFor(v interface{}) *itemsRuler {
	i := &itemsRuler{}
	switch tv := v.(type) {
	case *spec.Items:
		// simple items
		if tv.Items != nil {
			i.Child = newItemsRulerFor(tv.Items)
		} else {
			i.byType = newTypeRulerFor(tv)
			i.byFormat = newFormatRulerFor(tv)
			i.byPattern = newPatternRulerFor(tv)
		}
	default:
		return nil
	}
	return i
}

// Decide takes a decision for items
func (i *itemsRuler) Decide() basicGeneratorOpts {
	debugLog("itemsRuler.Decide()")
	if i == nil {
		return nil
	}
	current := i
	for current.Child != nil {
		current = current.Child
	}
	decisions := []basicGeneratorOpts{
		current.byType.Decide(),
		current.byFormat.Decide(),
		current.byPattern.Decide(),
	}
	g := &genOpts{}
	g.merge(decisions)
	return g
}

// schema infers options for schemas
type schemaRuler struct {
	ruler
	Schema *spec.Schema
}

// newSchemaRuler instantiate a new ruler for schema-based decisions
func newSchemaRulerFor(v interface{}) *schemaRuler {
	// TODO
	i := &schemaRuler{}
	return i
}

func (s *schemaRuler) Decide() basicGeneratorOpts {
	debugLog("schemaRuler.Decide()")
	// TODO
	g := &genOpts{}
	return g
}

// scoreAgainstProposals computes the similarity score for a text against the proposed aliased keys
func scoreAgainstProposals(text string, proposals map[string]string) (score map[string]float64) {
	score = make(map[string]float64, 300)
	// trivial implem for demo purpose
	for alias, generator := range proposals {
		if strings.Contains(text, alias) {
			// TODO: here should be a better appreciation on how close we are from the key
			score[generator] = score[generator] + 1
		} else {
			score[generator] = 0
		}
	}
	return
}

// mergeAndElectProposal elects the best proposal from a slice of scoring maps
func mergeAndElectProposal(scores ...map[string]float64) (key string) {
	var merger = make(map[string]float64, 300)
	for _, scoringMap := range scores {
		for k, score := range scoringMap {
			merger[k] = merger[k] + score
		}
	}
	best := float64(0)
	for k, s := range merger {
		if s > best {
			key = k
		}
	}
	return
}
