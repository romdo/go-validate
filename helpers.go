package validate

import (
	"fmt"
	"reflect"
	"regexp"
)

// RequireField returns a Error type for the given field if provided value is
// empty/zero.
func RequireField(field string, value interface{}) error {
	err := &Error{Field: field, Msg: "is required"}
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return err
		}
		v = v.Elem()
	}

	if !v.IsValid() {
		return err
	}

	switch v.Kind() { //nolint:exhaustive
	case reflect.Map, reflect.Slice:
		if v.Len() == 0 {
			return err
		}
	default:
		if v.IsZero() {
			return err
		}
	}

	return nil
}

func InRange(field string, value interface{}, min, max float64) error {
	if value == nil {
		return &Error{Field: field, Msg: "cannot be nil"}
	}

	kind := reflect.TypeOf(value).Kind()
	var floatValue float64

	switch kind { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		floatValue = float64(reflect.ValueOf(value).Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		floatValue = float64(reflect.ValueOf(value).Uint())
	case reflect.Float32, reflect.Float64:
		floatValue = reflect.ValueOf(value).Float()
	default:
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("unsupported type %T for InRange", value),
		}
	}

	if floatValue < min || floatValue > max {
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("must be in range [%.6f, %.6f]", min, max),
		}
	}

	return nil
}

func MinLength(field string, value interface{}, minLength int) error {
	if value == nil {
		return &Error{Field: field, Msg: "cannot be nil"}
	}

	if minLength < 0 {
		return &Error{Field: field, Msg: "minLength must be non-negative"}
	}

	kind := reflect.TypeOf(value).Kind()
	var length int

	switch kind { //nolint:exhaustive
	case reflect.String:
		length = len(reflect.ValueOf(value).String())
	case reflect.Slice, reflect.Array, reflect.Map:
		length = reflect.ValueOf(value).Len()
	default:
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("unsupported type %T for MinLength", value),
		}
	}

	if length < minLength {
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("must have a minimum length of %d", minLength),
		}
	}

	return nil
}

func MaxLength(field string, value interface{}, maxLength int) error {
	if value == nil {
		return &Error{Field: field, Msg: "cannot be nil"}
	}

	if maxLength < 0 {
		return &Error{Field: field, Msg: "maxLength must be non-negative"}
	}

	kind := reflect.TypeOf(value).Kind()
	var length int

	switch kind { //nolint:exhaustive
	case reflect.String:
		length = len(reflect.ValueOf(value).String())
	case reflect.Slice, reflect.Array, reflect.Map:
		length = reflect.ValueOf(value).Len()
	default:
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("unsupported type %T for MaxLength", value),
		}
	}

	if length > maxLength {
		return &Error{
			Field: field,
			Msg:   fmt.Sprintf("must have a maximum length of %d", maxLength),
		}
	}

	return nil
}

// MatchesRegexp checks if the value of a field matches the specified regular
// expression. It returns an error if the value doesn't match the pattern.
func MatchRegexp(
	field string,
	value interface{},
	pattern *regexp.Regexp,
) error {
	if pattern == nil {
		return &Error{Field: field, Msg: "pattern cannot be nil"}
	}

	switch v := value.(type) {
	case string:
		if !pattern.MatchString(v) {
			return &Error{
				Field: field,
				Msg: fmt.Sprintf(
					"does not match pattern '%s': '%s'", pattern, v,
				),
			}
		}
	case []byte:
		if !pattern.Match(v) {
			return &Error{
				Field: field,
				Msg: fmt.Sprintf(
					"does not match pattern '%s': '%s'", pattern, string(v),
				),
			}
		}
	default:
		return &Error{
			Field: field,
			Msg: fmt.Sprintf(
				"unsupported type %T for MatchRegexp", value,
			),
		}
	}

	return nil
}

func NotNil(field string, value interface{}) error {
	val := reflect.ValueOf(value)
	kind := val.Kind()
	isNil := false

	if kind == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
		kind = val.Kind()
	}

	switch kind { //nolint:exhaustive
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Slice,
		reflect.Ptr:
		if val.IsNil() {
			isNil = true
		}
	case reflect.Invalid:
		isNil = true
	default:
	}

	if isNil {
		return &Error{Field: field, Msg: "must not be nil"}
	}

	return nil
}
