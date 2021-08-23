package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/multierr"
)

// FieldNameFunc is a function which converts a given reflect.StructField to a
// string. The default will lookup json, yaml, and form field tags.
type FieldNameFunc func(reflect.StructField) string

// FieldJoinFunc joins a path slice with a given field. Both path and field may
// be empty values.
type FieldJoinFunc func(path []string, field string) string

// Validator validates Validatable objects.
type Validator struct {
	fieldName FieldNameFunc
	fieldJoin FieldJoinFunc
}

// New creates a new Validator.
func New() *Validator {
	return &Validator{}
}

// Validate will validate the given object. Structs, maps, slices, and arrays
// will have each of their fields/items validated, effectively performing a
// deep-validation.
func (s *Validator) Validate(data interface{}) error {
	if s.fieldName == nil {
		s.fieldName = DefaultFieldName
	}

	if s.fieldJoin == nil {
		s.fieldJoin = DefaultFieldJoin
	}

	return s.validate(nil, data)
}

// FieldNameFunc allows setting a custom FieldNameFunc method. It receives a
// reflect.StructField, and must return a string for the name of that field. If
// the returned string is empty, validation will not run against the field's
// value, or any nested data within.
func (s *Validator) FieldNameFunc(f FieldNameFunc) {
	s.fieldName = f
}

// FieldJoinFunc allows setting a custom FieldJoinFunc method. It receives a
// string slice of parent fields, and a string of the field name the error is
// reported against. All parent paths, must be joined with the current.
func (s *Validator) FieldJoinFunc(f FieldJoinFunc) {
	s.fieldJoin = f
}

func (s *Validator) validate(path []string, data interface{}) error {
	var errs error
	if data == nil {
		return nil
	}
	d := reflect.ValueOf(data)
	if d.Kind() == reflect.Ptr {
		if d.IsNil() {
			return nil
		}
		d = d.Elem()
	}

	if v, ok := data.(Validatable); ok {
		verrs := v.Validate()
		for _, err := range multierr.Errors(verrs) {
			// Create a new Error for all errors returned by Validate function
			// to correctly resolve field name, and also field path in relation
			// to parent objects being validated.
			newErr := &Error{}

			e := &Error{}
			if ok := errors.As(err, &e); ok {
				field := e.Field
				if field != "" && d.Kind() == reflect.Struct {
					if sf, ok := d.Type().FieldByName(e.Field); ok {
						field = s.fieldName(sf)
					}
				}
				newErr.Field = s.fieldJoin(path, field)
				newErr.Msg = e.Msg
				newErr.Err = e.Err
			} else {
				newErr.Field = s.fieldJoin(path, "")
				newErr.Err = err
			}

			errs = multierr.Append(errs, newErr)
		}
	}

	switch d.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Array:
		for i := 0; i < d.Len(); i++ {
			v := d.Index(i)
			err := s.validate(append(path, strconv.Itoa(i)), v.Interface())
			errs = multierr.Append(errs, err)
		}
	case reflect.Map:
		for _, k := range d.MapKeys() {
			v := d.MapIndex(k)
			err := s.validate(append(path, fmt.Sprintf("%v", k)), v.Interface())
			errs = multierr.Append(errs, err)
		}
	case reflect.Struct:
		for i := 0; i < d.NumField(); i++ {
			v := d.Field(i)
			fldName := s.fieldName(d.Type().Field(i))
			if v.CanSet() && fldName != "" {
				err := s.validate(append(path, fldName), v.Interface())
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}

// DefaultFieldName is the default FieldNameFunc used by Validator.
//
// Uses json, yaml, and form field tags to lookup field name first.
func DefaultFieldName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "" {
		name = strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]
	}

	if name == "" {
		name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
	}

	if name == "-" {
		return ""
	}

	if name == "" {
		return fld.Name
	}

	return name
}

// DefaultFieldJoin is the default FieldJoinFunc used by Validator.
func DefaultFieldJoin(path []string, field string) string {
	if field != "" {
		path = append(path, field)
	}

	return strings.Join(path, ".")
}
