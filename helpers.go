package validate

import (
	"reflect"
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
