package validate

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Test helper types
//

type testStruct struct {
	Foo string `json:"foo"`

	f func() error
}

func (s *testStruct) Validate() error {
	if s.f == nil {
		return nil
	}

	return s.f()
}

type testNestedStruct struct {
	OtherField *testStruct `yaml:"other_field"`
}

type MyStruct struct {
	Name string
	Kind string
}

//
// Tests
//

func TestNew(t *testing.T) {
	got := New()

	assert.IsType(t, &Validator{}, got)
}

func TestValidator_FieldNameFunc(t *testing.T) {
	v := New()
	v.FieldNameFunc(func(sf reflect.StructField) string {
		return "<" + strings.ToUpper(sf.Name) + ">"
	})
	err := v.Validate(&testNestedStruct{
		OtherField: &testStruct{f: func() error {
			return &Error{Field: "Foo", Msg: "oops"}
		}},
	})

	got := Errors(err)

	assert.ElementsMatch(t, []error{
		&Error{Field: "<OTHERFIELD>.<FOO>", Msg: "oops"},
	}, got)
}

func TestValidator_FieldJoinFunc(t *testing.T) {
	v := New()
	v.FieldJoinFunc(func(path []string, field string) string {
		if field != "" {
			path = append(path, field)
		}

		return "[" + strings.Join(path, "][") + "]"
	})
	err := v.Validate(&testNestedStruct{
		OtherField: &testStruct{f: func() error {
			return &Error{Field: "Foo", Msg: "oops"}
		}},
	})

	got := Errors(err)

	assert.ElementsMatch(t, []error{
		&Error{Field: "[other_field][foo]", Msg: "oops"},
	}, got)
}
