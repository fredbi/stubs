package stubs

import (
	"github.com/go-openapi/spec"
)

// headerGenOpts generates stubs for a swagger header in response
func headerGenOpts(key string, header *spec.Header) (*simpleOpts, error) {
	var gopts genOpts
	if header.Type == "" { // Safeguard
		return nil, nil
	}
	if err := gopts.ExtOverride(header.Extensions); err != nil {
		debugLog("extension error on %s: %v", gopts.Name(), err)
		return nil, err
	}
	debugLog("generator override: %s", gopts.Name())
	if gopts.Name() == "" {
		// register rules to infer options
		gopts.rules = make([]ruler, 0, 50)

		gopts.rules = append(gopts.rules, newTypeRulerFor(header))

		if header.Type == "string" {
			gopts.rules = append(gopts.rules, newFuzzyRulerFor(header))
		}

		if header.Format != "" {
			gopts.rules = append(gopts.rules, newFormatRulerFor(header))
		}

		if header.Pattern != "" {
			gopts.rules = append(gopts.rules, newPatternRulerFor(header))
		}

		if header.Items != nil {
			gopts.rules = append(gopts.rules, newItemsRulerFor(header.Items))
		}
	}

	if key == "" {
		key = "header"
	}

	return &simpleOpts{
		genOpts:           gopts,
		fieldName:         key,
		CommonValidations: header.CommonValidations,
		SimpleSchema:      header.SimpleSchema,
		required:          true,
	}, nil
}
