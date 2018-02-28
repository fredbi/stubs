package stubs

import (
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
)

type simpleOpts struct {
	genOpts
	defaultSeeder
	spec.CommonValidations
	spec.SimpleSchema

	fieldName string
	required  bool
}

func (g *simpleOpts) FieldName() string {
	return g.fieldName
}
func (g *simpleOpts) Maximum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Maximum), g.CommonValidations.ExclusiveMaximum, g.CommonValidations.Maximum != nil
}
func (g *simpleOpts) Minimum() (float64, bool, bool) {
	return swag.Float64Value(g.CommonValidations.Minimum), g.CommonValidations.ExclusiveMinimum, g.CommonValidations.Minimum != nil
}
func (g *simpleOpts) MaxLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MaxLength), g.CommonValidations.MaxLength != nil
}
func (g *simpleOpts) MinLength() (int64, bool) {
	return swag.Int64Value(g.CommonValidations.MinLength), g.CommonValidations.MinLength != nil
}
func (g *simpleOpts) Pattern() (string, bool) {
	return g.CommonValidations.Pattern, g.CommonValidations.Pattern != ""
}
func (g *simpleOpts) MaxItems() (int64, bool) {
	mx := g.CommonValidations.MaxItems
	return swag.Int64Value(mx), mx != nil
}
func (g *simpleOpts) MinItems() (int64, bool) {
	mn := g.CommonValidations.MinItems
	return swag.Int64Value(mn), mn != nil
}
func (g *simpleOpts) UniqueItems() bool {
	return g.CommonValidations.UniqueItems
}
func (g *simpleOpts) MultipleOf() (float64, bool) {
	mo := g.CommonValidations.MultipleOf
	return swag.Float64Value(mo), mo != nil
}
func (g *simpleOpts) Enum() ([]interface{}, bool) {
	enm := g.CommonValidations.Enum
	return enm, len(enm) > 0
}
func (g *simpleOpts) Type() string {
	return g.SimpleSchema.Type
}
func (g *simpleOpts) Format() string {
	return g.SimpleSchema.Format
}

// TODO: remove this, rulers replace it now
func (g *simpleOpts) Items() (GeneratorOpts, error) {
	return itemsGenOpts(g.name+".items", g.SimpleSchema.Items)
}
func (g *simpleOpts) Required() bool {
	return g.required
}
