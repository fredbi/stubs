package stubs

import (
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
)

func schemaGenOpts(key string, required bool, schema *spec.Schema) (*schemaOpts, error) {
	var gopts genOpts
	if err := gopts.ExtOverride(schema.Extensions); err != nil {
		debugLog("extension error on %s: %v", gopts.Name(), err)
		return nil, err
	}
	debugLog("generator override: %s", gopts.Name())

	if gopts.Name() == "" {
		// register rules to infer options
		gopts.rules = make([]ruler, 0, 50)

		gopts.rules = append(gopts.rules, newSchemaRulerFor(schema))
	}
	return &schemaOpts{
		genOpts:   gopts,
		fieldName: key,
		schema:    schema,
		required:  required,
	}, nil
}

type schemaOpts struct {
	genOpts
	defaultSeeder
	schema *spec.Schema

	fieldName string
	required  bool
}

func (s *schemaOpts) FieldName() string {
	return s.fieldName
}
func (s *schemaOpts) Maximum() (float64, bool, bool) {
	return swag.Float64Value(s.schema.Maximum), s.schema.ExclusiveMaximum, s.schema.Maximum != nil
}
func (s *schemaOpts) Minimum() (float64, bool, bool) {
	return swag.Float64Value(s.schema.Minimum), s.schema.ExclusiveMinimum, s.schema.Minimum != nil
}
func (s *schemaOpts) MaxLength() (int64, bool) {
	return swag.Int64Value(s.schema.MaxLength), s.schema.MaxLength != nil
}
func (s *schemaOpts) MinLength() (int64, bool) {
	return swag.Int64Value(s.schema.MinLength), s.schema.MinLength != nil
}
func (s *schemaOpts) Pattern() (string, bool) {
	return s.schema.Pattern, s.schema.Pattern != ""
}
func (s *schemaOpts) MaxItems() (int64, bool) {
	mx := s.schema.MaxItems
	return swag.Int64Value(mx), mx != nil
}
func (s *schemaOpts) MinItems() (int64, bool) {
	mn := s.schema.MinItems
	return swag.Int64Value(mn), mn != nil
}
func (s *schemaOpts) UniqueItems() bool {
	return s.schema.UniqueItems
}
func (s *schemaOpts) MultipleOf() (float64, bool) {
	mo := s.schema.MultipleOf
	return swag.Float64Value(mo), mo != nil
}
func (s *schemaOpts) Enum() ([]interface{}, bool) {
	enm := s.schema.Enum
	return enm, len(enm) > 0
}
func (s *schemaOpts) Type() string {
	if len(s.schema.Type) == 0 {
		return "object"
	}
	// NOTE: does not support multiple types
	return s.schema.Type[0]
}
func (s *schemaOpts) Format() string {
	return s.schema.Format
}
func (s *schemaOpts) Items() (GeneratorOpts, error) {
	return schemaGenOpts(s.fieldName+".items", false, s.schema.Items.Schema)
}
func (s *schemaOpts) Required() bool {
	return s.required
}
