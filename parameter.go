package stubs

import (
	"github.com/go-openapi/spec"
)

// paramGenOpts generates stubs for a swagger simple parameter
func paramGenOpts(key string, param *spec.Parameter) (*simpleOpts, error) {
	var gopts genOpts
	if param.Schema != nil || param.Type == "" { // Safeguard
		// parameters with schema must be handled by schemaGenOpts
		return nil, nil
	}
	if err := gopts.ExtOverride(param.Extensions); err != nil {
		debugLog("extension error on %s: %v", gopts.Name(), err)
		return nil, err
	}
	debugLog("generator override: %s", gopts.Name())
	if gopts.Name() == "" {
		// register rules to infer options
		gopts.rules = make([]ruler, 0, 50)

		gopts.rules = append(gopts.rules, newTypeRulerFor(param))

		if param.Type == "string" {
			gopts.rules = append(gopts.rules, newFuzzyRulerFor(param))
		}

		if param.Format != "" {
			gopts.rules = append(gopts.rules, newFormatRulerFor(param))
		}

		if param.Pattern != "" {
			gopts.rules = append(gopts.rules, newPatternRulerFor(param))
		}

		if param.Items != nil {
			gopts.rules = append(gopts.rules, newItemsRulerFor(param.Items))
		}
	}

	if key == "" {
		key = param.Name
	}

	return &simpleOpts{
		genOpts:           gopts,
		fieldName:         key,
		CommonValidations: param.CommonValidations,
		SimpleSchema:      param.SimpleSchema,
		required:          param.Required,
	}, nil
}
