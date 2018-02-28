package stubs

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// Generator generates a stub for a descriptor.
// A descriptor can either be a parameter, response header or json schema
type Generator struct {
	Language string
	// Args represents general settings for the generator. These may be overriden by local x-datagen extensions.
	//Args genTag
}

//TODO: replace key by unstructured hint

// Generate a stub from swagger spec constructs into the opts.Target
func (s *Generator) Generate(key string, descriptor interface{}) (interface{}, error) {
	switch desc := descriptor.(type) {
	case *spec.Parameter:
		return s.GenParameter(key, desc)
	case spec.Parameter:
		return s.GenParameter(key, &desc)
	case *spec.Header:
		return s.GenHeader(key, desc)
	case spec.Header:
		return s.GenHeader(key, &desc)
	case *spec.Schema:
		return s.GenSchema(key, desc)
	case spec.Schema:
		return s.GenSchema(key, &desc)
	case *spec.Response:
		return s.GenResponse(key, desc)
	case spec.Response:
		return s.GenResponse(key, &desc)
	default:
		return nil, fmt.Errorf("%T is unsupported for Generator", descriptor)
	}
}

// GenResponse generates a random value for a response
func (s *Generator) GenResponse(key string, response *spec.Response) (interface{}, error) {
	debugLog("generation for response: %q", response.Description)
	// TODO: push downstream description and code (?) to help fuzzying
	return s.GenSchema(key, response.Schema)
}

// GenParameter generates a random value for a parameter
func (s *Generator) GenParameter(key string, param *spec.Parameter) (interface{}, error) {
	debugLog("generation for parameter: %s", param.Name)
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	if param.Schema != nil {
		// TODO: push downstream name and description to help fuzzying
		return s.GenSchema(key, param.Schema)
	}
	debugLog("resolving generation options for parameter: %s, type: %s, format: %s, isSchema: %t", param.Name, param.Type, param.Format, param.Schema != nil)
	gopts, err := paramGenOpts(key, param)
	if err != nil {
		return nil, err
	}

	// infer options
	gopts.Infer()

	//if args.WithEdgeCase.Enabled && (len(args.WithEdgeCase.Args)==0 || args.WithEdgeCase.Args.Contains("standard")) {
	// Standard edge cases are boundary values
	// Edge cases managed at a higher level
	//}
	debugLog("generation options determined for parameter: %#v", gopts)

	// TODO : slices
	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for parameter [%s]", param.Name)
	}

	return datagen(gopts)
}

// GenHeader generates a random value for a header
func (s *Generator) GenHeader(key string, header *spec.Header) (interface{}, error) {
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	gopts, err := headerGenOpts(key, header)
	if err != nil {
		return nil, err
	}

	gopts.Infer()

	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for header [%s]", key)
	}

	return datagen(gopts)
}

// GenSchema generates a random value for a schema
func (s *Generator) GenSchema(key string, schema *spec.Schema) (interface{}, error) {
	generator, err := newGenerator(s.Language)
	if err != nil {
		return nil, err
	}

	gopts, err := schemaGenOpts(key, true, schema)
	if err != nil {
		return nil, err
	}

	datagen, found := generator.For(gopts)
	if !found {
		return nil, fmt.Errorf("no generator found for schema [%s]", key)
	}

	return datagen(gopts)
}
