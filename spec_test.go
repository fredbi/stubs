// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stubs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/stretchr/testify/assert"
)

var (
	// This debug environment variable allows to report and capture actual validation messages
	// during testing. It should be disabled (undefined) during CI tests.
	DebugTest = os.Getenv("SWAGGER_DEBUG_TEST") != ""
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func TestGenerator_Spec(t *testing.T) {
	errs := testStubsGenerator(t, true, true) /* set haltOnErrors=true to iterate spec by spec */
	assert.Zero(t, errs, "Generator testing didn't match expectations")
}

func testStubsGenerator(t *testing.T, haltOnErrors bool, continueOnErrors bool) (errs int) {
	err := filepath.Walk(filepath.Join("fixtures", "specs"),
		func(path string, info os.FileInfo, err error) error {
			basename := info.Name()
			if !info.IsDir() {
				t.Logf("Generating stubs for spec: %s", basename)
				doc, err := loads.Spec(path)
				if assert.NoError(t, err, "Expected this spec to load properly") {
					// Validate the spec document
					validator := validate.NewSpecValidator(doc.Schema(), strfmt.Default)
					validator.SetContinueOnErrors(continueOnErrors)
					res, _ := validator.Validate(doc)
					// Check specs with load errors (error is located in pkg loads or spec)
					if assert.Truef(t, res.IsValid(), "Expected this spec to be valid: %v", res.Errors) {
						// Now walks the spec to generate stubs
						analyzer := analysis.New(doc.Spec())
						for method, pathItem := range analyzer.Operations() {
							if pathItem != nil { // Safeguard
								for path := range pathItem {
									params := []spec.Parameter{}
									// parameters
									for _, ppr := range analyzer.ParamsFor(method, path) {
										// Expand params (TODO: update with current expand method)
										pr := ppr
										sw := doc.Spec()
										for pr.Ref.String() != "" {
											obj, _, _ := pr.Ref.GetPointer().Get(sw)
											pr = obj.(spec.Parameter)
										}
										params = append(params, pr)
									}
									// TODO: should sort params to get a repeatable random generation
									for _, param := range params {
										//Debug = true
										t.Logf("Parameter: [name: %s, in:%s]", param.Name, param.In)
										fixture, err := generateStubsForParam(t, path, method, param)
										assert.NoError(t, err)
										t.Logf("Stub: %v", fixture)
									}
								}
							}
						}
					} else {
						errs++
					}

					//
				} else {
					errs++
				}
			}
			if haltOnErrors && errs > 0 {
				return fmt.Errorf("Test halted: stop testing on stubs generation")
			}
			return nil
		})
	if err != nil {
		t.Logf("%v", err)
		errs++
	}
	return
}

type FixtureParam struct {
	Name    string
	In      string
	Example interface{}
}

type GeneratedFixture struct {
	Params       []*FixtureParam
	ParamsAsJSON json.RawMessage
}

func generateStubsForParam(t *testing.T, path, method string, param spec.Parameter) (*FixtureParam, error) {
	gen := Generator{"en"}
	//strings.Join([]string{method, path}, ":")
	result, err := gen.Generate("", &param)
	assert.NoError(t, err)
	return &FixtureParam{Name: param.Name, In: param.In, Example: result}, err
}

// Test unitary fixture for dev and bug fixing
func Test_SingleFixture(t *testing.T) {
	t.SkipNow()
	path := "fixtures/validation/gentest3.yaml"
	doc, err := loads.Spec(path)
	if assert.NoError(t, err) {
		validator := validate.NewSpecValidator(doc.Schema(), strfmt.Default)
		validator.SetContinueOnErrors(true)
		res, _ := validator.Validate(doc)
		t.Log("Returned errors:")
		for _, e := range res.Errors {
			t.Logf("%v", e)
		}
		t.Log("Returned warnings:")
		for _, e := range res.Warnings {
			t.Logf("%v", e)
		}

	} else {
		t.Logf("Load error: %v", err)
	}
}
