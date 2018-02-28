package stubs

import (
	"github.com/go-openapi/spec"
)

// NOTE: unused
func itemsGenOpts(key string, items *spec.Items) (*simpleOpts, error) {
	var gopts genOpts
	if err := gopts.ExtOverride(items.Extensions); err != nil {
		return nil, err
	}
	return &simpleOpts{
		genOpts:           gopts,
		fieldName:         key,
		CommonValidations: items.CommonValidations,
		SimpleSchema:      items.SimpleSchema,
		required:          true,
	}, nil
}
