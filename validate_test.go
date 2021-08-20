package validate_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/romdo/go-validate"
	"github.com/stretchr/testify/assert"
)

//
// Test helper types
//

type validatableString string

func (s validatableString) Validate() error {
	if strings.Contains(string(s), " ") {
		return &validate.Error{Msg: "must not contain space"}
	}

	return nil
}

type validatableStruct struct {
	Foo string
	Bar string `json:"bar"`
	Foz string `yaml:"foz"`
	Baz string `form:"baz"`

	f func() error
}

func (s *validatableStruct) Validate() error {
	if s.f == nil {
		return nil
	}

	return s.f()
}

type nestedStruct struct {
	OtherField     *validatableStruct
	OtherFieldJSON *validatableStruct `json:"other_field,omitempty"`
	OtherFieldYAML *validatableStruct `yaml:"otherField,omitempty"`
	OtherFieldFORM *validatableStruct `form:"other-field,omitempty"`

	skippedField     *validatableStruct
	SkippedFieldJSON *validatableStruct `json:"-,omitempty"`
	SkippedFieldYAML *validatableStruct `yaml:"-,omitempty"`
	SkippedFieldFORM *validatableStruct `form:"-,omitempty"`

	OtherArray     [5]*validatableStruct
	OtherArrayJSON [5]*validatableStruct `json:"other_array,omitempty"`
	OtherArrayYAML [5]*validatableStruct `yaml:"otherArray,omitempty"`
	OtherArrayFORM [5]*validatableStruct `form:"other-array,omitempty"`

	OtherSlice     []*validatableStruct
	OtherSliceJSON []*validatableStruct `json:"other_slice,omitempty"`
	OtherSliceYAML []*validatableStruct `yaml:"otherSlice,omitempty"`
	OtherSliceFORM []*validatableStruct `form:"other-slice,omitempty"`

	OtherStringMap     map[string]*validatableStruct
	OtherStringMapJSON map[string]*validatableStruct `json:"other_string_map,omitempty"`
	OtherStringMapYAML map[string]*validatableStruct `yaml:"otherStringMap,omitempty"`
	OtherStringMapFORM map[string]*validatableStruct `form:"other-string-map,omitempty"`

	OtherIntMap     map[int]*validatableStruct
	OtherIntMapJSON map[int]*validatableStruct `json:"other_int_map,omitempty"`
	OtherIntMapYAML map[int]*validatableStruct `yaml:"otherIntMap,omitempty"`
	OtherIntMapFORM map[int]*validatableStruct `form:"other-int-map,omitempty"`

	OtherStruct     *nestedStruct
	OtherStructJSON *nestedStruct `json:"other_struct,omitempty"`
	OtherStructYAML *nestedStruct `yaml:"otherStruct,omitempty"`
	OtherStructFORM *nestedStruct `form:"other-struct,omitempty"`
}

//
// Tests
//

func TestValidate(t *testing.T) {
	mixedValidationErrors := &validatableStruct{
		f: func() error {
			var errs error
			errs = validate.Append(errs, &validate.Error{
				Field: "Foo", Msg: "is required",
				Err: errors.New("oops"),
			})
			errs = validate.Append(errs, errors.New("bar: is missing"))

			return errs
		},
	}

	tests := []struct {
		name     string
		obj      interface{}
		wantErrs []error
	}{
		{
			name:     "nil",
			obj:      nil,
			wantErrs: []error{},
		},
		{
			name:     "no error",
			obj:      &validatableStruct{},
			wantErrs: nil,
		},
		{
			name:     "valid validatable string type",
			obj:      validatableString("hello-world"),
			wantErrs: []error{},
		},
		{
			name: "invalid validatable string type",
			obj:  validatableString("hello world"),
			wantErrs: []error{
				&validate.Error{Msg: "must not contain space"},
			},
		},
		{
			name: "single Go error",
			obj: &validatableStruct{f: func() error {
				return errors.New("foo: is required")
			}},
			wantErrs: []error{
				&validate.Error{Err: errors.New("foo: is required")},
			},
		},
		{
			name: "single *validate.Error",
			obj: &validatableStruct{f: func() error {
				return &validate.Error{
					Field: "foo", Msg: "is required", Err: errors.New("oops"),
				}
			}},
			wantErrs: []error{
				&validate.Error{
					Field: "foo", Msg: "is required", Err: errors.New("oops"),
				},
			},
		},
		{
			name: "multiple Go errors",
			obj: &validatableStruct{f: func() error {
				var errs error
				errs = validate.Append(errs, errors.New("foo: is required"))
				errs = validate.Append(errs, errors.New("bar: is missing"))

				return errs
			}},
			wantErrs: []error{
				&validate.Error{Err: errors.New("foo: is required")},
				&validate.Error{Err: errors.New("bar: is missing")},
			},
		},
		{
			name: "multiple *validate.Error",
			obj: &validatableStruct{f: func() error {
				var errs error
				errs = validate.Append(errs, &validate.Error{
					Field: "foo", Msg: "is required", Err: errors.New("oops"),
				})
				errs = validate.Append(errs, &validate.Error{
					Field: "bar", Msg: "is required", Err: errors.New("whoops"),
				})

				return errs
			}},
			wantErrs: []error{
				&validate.Error{
					Field: "foo", Msg: "is required", Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "bar", Msg: "is required", Err: errors.New("whoops"),
				},
			},
		},
		{
			name: "mix of Go error and *validate.Error",
			obj:  mixedValidationErrors,
			wantErrs: []error{
				&validate.Error{
					Field: "Foo", Msg: "is required", Err: errors.New("oops"),
				},
				&validate.Error{Err: errors.New("bar: is missing")},
			},
		},
		//
		// Field name conversion
		//
		{
			name: "no json, yaml or form field tag",
			obj: &validatableStruct{f: func() error {
				return &validate.Error{Field: "Foo", Msg: "is required"}
			}},
			wantErrs: []error{
				&validate.Error{Field: "Foo", Msg: "is required"},
			},
		},
		{
			name: "converts field name via json field tag",
			obj: &validatableStruct{f: func() error {
				return &validate.Error{Field: "Bar", Msg: "is required"}
			}},
			wantErrs: []error{
				&validate.Error{Field: "bar", Msg: "is required"},
			},
		},
		{
			name: "converts field name via yaml field tag",
			obj: &validatableStruct{f: func() error {
				return &validate.Error{Field: "Foz", Msg: "is required"}
			}},
			wantErrs: []error{
				&validate.Error{Field: "foz", Msg: "is required"},
			},
		},
		{
			name: "converts field name via form field tag",
			obj: &validatableStruct{f: func() error {
				return &validate.Error{Field: "Baz", Msg: "is required"}
			}},
			wantErrs: []error{
				&validate.Error{Field: "baz", Msg: "is required"},
			},
		},
		{
			name: "nested with no validation errors",
			obj: &nestedStruct{
				OtherField: &validatableStruct{},
				OtherArray: [5]*validatableStruct{{}, {}, {}, {}},
				OtherSlice: []*validatableStruct{{}},
				OtherStringMap: map[string]*validatableStruct{
					"hi":  {},
					"bye": {},
				},
				OtherIntMap: map[int]*validatableStruct{42: {}, 64: {}},
				OtherStruct: &nestedStruct{
					OtherField: &validatableStruct{},
				},
			},
			wantErrs: []error{},
		},
		//
		// Nested in a struct field.
		//
		{
			name: "nested in a struct field",
			obj: &nestedStruct{
				OtherField: mixedValidationErrors,
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherField.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherField", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a struct field with json tag",
			obj: &nestedStruct{
				OtherFieldJSON: mixedValidationErrors,
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_field.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_field", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a struct field with yaml tag",
			obj: &nestedStruct{
				OtherFieldYAML: mixedValidationErrors,
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherField.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherField", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a struct field with form tag",
			obj: &nestedStruct{
				OtherFieldFORM: mixedValidationErrors,
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-field.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-field", Err: errors.New("bar: is missing"),
				},
			},
		},
		//
		// Nested in a unexposed/ignored fields.
		//
		{
			name: "nested in a unexposed field",
			obj: &nestedStruct{
				skippedField: mixedValidationErrors,
			},
			wantErrs: []error{},
		},
		{
			name: "nested in a struct field skipped by json tag",
			obj: &nestedStruct{
				SkippedFieldJSON: mixedValidationErrors,
			},
			wantErrs: []error{},
		},
		{
			name: "nested in a struct field skipped by yaml tag",
			obj: &nestedStruct{
				SkippedFieldYAML: mixedValidationErrors,
			},
			wantErrs: []error{},
		},
		{
			name: "nested in a struct field skipped by yaml tag",
			obj: &nestedStruct{
				SkippedFieldFORM: mixedValidationErrors,
			},
			wantErrs: []error{},
		},
		//
		// Nested in an array.
		//
		{
			name: "nested in an array",
			obj: &nestedStruct{
				OtherArray: [5]*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherArray.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherArray.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "OtherArray.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherArray.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in an array with json tag",
			obj: &nestedStruct{
				OtherArrayJSON: [5]*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_array.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_array.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other_array.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_array.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in an array with yaml tag",
			obj: &nestedStruct{
				OtherArrayYAML: [5]*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherArray.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherArray.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "otherArray.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherArray.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in an array with form tag",
			obj: &nestedStruct{
				OtherArrayFORM: [5]*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-array.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-array.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other-array.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-array.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		//
		// Nested in a slice.
		//
		{
			name: "nested in a slice",
			obj: &nestedStruct{
				OtherSlice: []*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherSlice.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherSlice.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "OtherSlice.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherSlice.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a slice with json tag",
			obj: &nestedStruct{
				OtherSliceJSON: []*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_slice.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_slice.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other_slice.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_slice.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a slice with yaml tag",
			obj: &nestedStruct{
				OtherSliceYAML: []*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherSlice.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherSlice.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "otherSlice.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherSlice.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a slice with form tag",
			obj: &nestedStruct{
				OtherSliceFORM: []*validatableStruct{
					mixedValidationErrors,
					mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-slice.0.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-slice.0", Err: errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other-slice.1.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-slice.1", Err: errors.New("bar: is missing"),
				},
			},
		},
		//
		// Nested in an string map.
		//
		{
			name: "nested in a string map",
			obj: &nestedStruct{
				OtherStringMap: map[string]*validatableStruct{
					"hello": mixedValidationErrors,
					"world": mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherStringMap.hello.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherStringMap.hello",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "OtherStringMap.world.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherStringMap.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a string map with json tag",
			obj: &nestedStruct{
				OtherStringMapJSON: map[string]*validatableStruct{
					"hello": mixedValidationErrors,
					"world": mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_string_map.hello.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_string_map.hello",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other_string_map.world.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_string_map.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a string map with yaml tag",
			obj: &nestedStruct{
				OtherStringMapYAML: map[string]*validatableStruct{
					"hello": mixedValidationErrors,
					"world": mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherStringMap.hello.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherStringMap.hello",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "otherStringMap.world.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherStringMap.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a string map with form tag",
			obj: &nestedStruct{
				OtherStringMapFORM: map[string]*validatableStruct{
					"hello": mixedValidationErrors,
					"world": mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-string-map.hello.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-string-map.hello",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other-string-map.world.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-string-map.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		//
		// Nested in an int map.
		//
		{
			name: "nested in a int map",
			obj: &nestedStruct{
				OtherIntMap: map[int]*validatableStruct{
					42: mixedValidationErrors,
					64: mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherIntMap.42.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherIntMap.42",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "OtherIntMap.64.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherIntMap.64",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with json tag",
			obj: &nestedStruct{
				OtherIntMapJSON: map[int]*validatableStruct{
					42: mixedValidationErrors,
					64: mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_int_map.42.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_int_map.42",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other_int_map.64.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_int_map.64",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with yaml tag",
			obj: &nestedStruct{
				OtherIntMapYAML: map[int]*validatableStruct{
					42: mixedValidationErrors,
					64: mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherIntMap.42.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherIntMap.42",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "otherIntMap.64.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherIntMap.64",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with form tag",
			obj: &nestedStruct{
				OtherIntMapFORM: map[int]*validatableStruct{
					42: mixedValidationErrors,
					64: mixedValidationErrors,
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-int-map.42.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-int-map.42",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other-int-map.64.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-int-map.64",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		//
		// Nested in another struct.
		//
		{
			name: "nested in another struct",
			obj: &nestedStruct{
				OtherStruct: &nestedStruct{
					OtherField: mixedValidationErrors,
					OtherStringMap: map[string]*validatableStruct{
						"world": mixedValidationErrors,
					},
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "OtherStruct.OtherField.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherStruct.OtherField",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "OtherStruct.OtherStringMap.world.Foo",
					Msg:   "is required",
					Err:   errors.New("oops"),
				},
				&validate.Error{
					Field: "OtherStruct.OtherStringMap.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with json tag",
			obj: &nestedStruct{
				OtherStructJSON: &nestedStruct{
					OtherFieldJSON: mixedValidationErrors,
					OtherStringMapJSON: map[string]*validatableStruct{
						"world": mixedValidationErrors,
					},
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other_struct.other_field.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other_struct.other_field",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other_struct.other_string_map.world.Foo",
					Msg:   "is required",
					Err:   errors.New("oops"),
				},
				&validate.Error{
					Field: "other_struct.other_string_map.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with yaml tag",
			obj: &nestedStruct{
				OtherStructYAML: &nestedStruct{
					OtherFieldYAML: mixedValidationErrors,
					OtherStringMapYAML: map[string]*validatableStruct{
						"world": mixedValidationErrors,
					},
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "otherStruct.otherField.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "otherStruct.otherField",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "otherStruct.otherStringMap.world.Foo",
					Msg:   "is required",
					Err:   errors.New("oops"),
				},
				&validate.Error{
					Field: "otherStruct.otherStringMap.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
		{
			name: "nested in a int map with form tag",
			obj: &nestedStruct{
				OtherStructFORM: &nestedStruct{
					OtherFieldFORM: mixedValidationErrors,
					OtherStringMapFORM: map[string]*validatableStruct{
						"world": mixedValidationErrors,
					},
				},
			},
			wantErrs: []error{
				&validate.Error{
					Field: "other-struct.other-field.Foo", Msg: "is required",
					Err: errors.New("oops"),
				},
				&validate.Error{
					Field: "other-struct.other-field",
					Err:   errors.New("bar: is missing"),
				},
				&validate.Error{
					Field: "other-struct.other-string-map.world.Foo",
					Msg:   "is required",
					Err:   errors.New("oops"),
				},
				&validate.Error{
					Field: "other-struct.other-string-map.world",
					Err:   errors.New("bar: is missing"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Validate(tt.obj)

			if len(tt.wantErrs) == 0 {
				assert.Nil(t, err, "validation error should be nil")
			}

			got := validate.Errors(err)
			assert.ElementsMatch(t, tt.wantErrs, got)
		})
	}
}
