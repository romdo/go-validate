package validate

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func TestRequireField(t *testing.T) {
	var nilMapString map[string]string
	emptyMapString := map[string]string{}
	mapString := map[string]string{"foo": "bar"}
	type testStruct struct {
		Name string
	}

	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				field: "Title",
				value: nil,
			},
			want: &Error{Field: "Title", Msg: "is required"},
		},
		{
			name: "nil pointer",
			args: args{
				field: "Title",
				value: &nilMapString,
			},
			want: &Error{Field: "Title", Msg: "is required"},
		},
		{
			name: "true boolean",
			args: args{
				field: "Book",
				value: true,
			},
			want: nil,
		},
		{
			name: "false boolean",
			args: args{
				field: "Book",
				value: false,
			},
			want: &Error{Field: "Book", Msg: "is required"},
		},
		{
			name: "int",
			args: args{
				field: "Count",
				value: int(834),
			},
			want: nil,
		},
		{
			name: "zero int",
			args: args{
				field: "Count",
				value: int(0),
			},
			want: &Error{Field: "Count", Msg: "is required"},
		},
		{
			name: "int8",
			args: args{
				field: "Ticks",
				value: int8(3),
			},
			want: nil,
		},
		{
			name: "zero int8",
			args: args{
				field: "Ticks",
				value: int8(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "int16",
			args: args{
				field: "Ticks",
				value: int16(3),
			},
			want: nil,
		},
		{
			name: "zero int16",
			args: args{
				field: "Ticks",
				value: int16(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "int32",
			args: args{
				field: "Ticks",
				value: int32(3),
			},
			want: nil,
		},
		{
			name: "zero int32",
			args: args{
				field: "Ticks",
				value: int32(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "int64",
			args: args{
				field: "Ticks",
				value: int64(3),
			},
			want: nil,
		},
		{
			name: "zero int64",
			args: args{
				field: "Ticks",
				value: int64(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "zero uint",
			args: args{
				field: "Count",
				value: uint(0),
			},
			want: &Error{Field: "Count", Msg: "is required"},
		},
		{
			name: "uint8",
			args: args{
				field: "Ticks",
				value: uint8(3),
			},
			want: nil,
		},
		{
			name: "zero uint8",
			args: args{
				field: "Ticks",
				value: uint8(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "uint16",
			args: args{
				field: "Ticks",
				value: uint16(3),
			},
			want: nil,
		},
		{
			name: "zero uint16",
			args: args{
				field: "Ticks",
				value: uint16(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "uint32",
			args: args{
				field: "Ticks",
				value: uint32(3),
			},
			want: nil,
		},
		{
			name: "zero uint32",
			args: args{
				field: "Ticks",
				value: uint32(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "uint64",
			args: args{
				field: "Ticks",
				value: uint64(3),
			},
			want: nil,
		},
		{
			name: "zero uint64",
			args: args{
				field: "Ticks",
				value: uint64(0),
			},
			want: &Error{Field: "Ticks", Msg: "is required"},
		},
		{
			name: "complex64",
			args: args{
				field: "Offset",
				value: complex64(3),
			},
			want: nil,
		},
		{
			name: "zero complex64",
			args: args{
				field: "Offset",
				value: complex64(0),
			},
			want: &Error{Field: "Offset", Msg: "is required"},
		},
		{
			name: "complex128",
			args: args{
				field: "Offset",
				value: complex128(3),
			},
			want: nil,
		},
		{
			name: "zero complex128",
			args: args{
				field: "Offset",
				value: complex128(0),
			},
			want: &Error{Field: "Offset", Msg: "is required"},
		},
		{
			name: "array",
			args: args{
				field: "List",
				value: [3]string{"foo", "bar", "baz"},
			},
			want: nil,
		},
		{
			name: "empty array",
			args: args{
				field: "List",
				value: [3]string{},
			},
			want: &Error{Field: "List", Msg: "is required"},
		},
		{
			name: "chan",
			args: args{
				field: "Comms",
				value: make(chan int),
			},
			want: nil,
		},
		{
			name: "func",
			args: args{
				field: "Callback",
				value: func() error { return nil },
			},
			want: nil,
		},
		{
			name: "map",
			args: args{
				field: "Lookup",
				value: map[string]string{"foo": "bar"},
			},
			want: nil,
		},
		{
			name: "map pointer",
			args: args{
				field: "Lookup",
				value: &mapString,
			},
			want: nil,
		},
		{
			name: "empty map",
			args: args{
				field: "Lookup",
				value: map[string]string{},
			},
			want: &Error{Field: "Lookup", Msg: "is required"},
		},
		{
			name: "empty map pointer",
			args: args{
				field: "Lookup",
				value: &emptyMapString,
			},
			want: &Error{Field: "Lookup", Msg: "is required"},
		},
		{
			name: "nil map",
			args: args{
				field: "Lookup",
				value: nilMapString,
			},
			want: &Error{Field: "Lookup", Msg: "is required"},
		},
		{
			name: "slice",
			args: args{
				field: "List",
				value: []string{"foo", "bar", "baz"},
			},
			want: nil,
		},
		{
			name: "empty slice",
			args: args{
				field: "List",
				value: []string{},
			},
			want: &Error{Field: "List", Msg: "is required"},
		},
		{
			name: "string",
			args: args{
				field: "Book",
				value: "foo",
			},
			want: nil,
		},
		{
			name: "string pointer",
			args: args{
				field: "Book",
				value: stringPtr("foo"),
			},
			want: nil,
		},
		{
			name: "empty string",
			args: args{
				field: "Book",
				value: "",
			},
			want: &Error{Field: "Book", Msg: "is required"},
		},
		{
			name: "empty string pointer",
			args: args{
				field: "Book",
				value: stringPtr(""),
			},
			want: &Error{Field: "Book", Msg: "is required"},
		},
		{
			name: "struct",
			args: args{
				field: "Thing",
				value: testStruct{Name: "hi"},
			},
			want: nil,
		},
		{
			name: "struct pointer",
			args: args{
				field: "Thing",
				value: &testStruct{Name: "hi"},
			},
			want: nil,
		},
		{
			name: "empty struct",
			args: args{
				field: "Thing",
				value: testStruct{},
			},
			want: &Error{Field: "Thing", Msg: "is required"},
		},
		{
			name: "empty struct pointer",
			args: args{
				field: "Thing",
				value: &testStruct{},
			},
			want: &Error{Field: "Thing", Msg: "is required"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequireField(tt.args.field, tt.args.value)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInRange(t *testing.T) {
	type args struct {
		field string
		value interface{}
		min   float64
		max   float64
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "int in range",
			args: args{
				field: "Age",
				value: int(25),
				min:   18,
				max:   65,
			},
			want: nil,
		},
		{
			name: "int below range",
			args: args{
				field: "Age",
				value: int(15),
				min:   18,
				max:   65,
			},
			want: &Error{
				Field: "Age",
				Msg:   "must be in range [18.000000, 65.000000]",
			},
		},
		{
			name: "int above range",
			args: args{
				field: "Age",
				value: int(70),
				min:   18,
				max:   65,
			},
			want: &Error{
				Field: "Age",
				Msg:   "must be in range [18.000000, 65.000000]",
			},
		},
		{
			name: "float in range",
			args: args{
				field: "Rating",
				value: float64(4.5),
				min:   1,
				max:   5,
			},
			want: nil,
		},
		{
			name: "float below range",
			args: args{
				field: "Rating",
				value: float64(0.5),
				min:   1,
				max:   5,
			},
			want: &Error{
				Field: "Rating",
				Msg:   "must be in range [1.000000, 5.000000]",
			},
		},
		{
			name: "float above range",
			args: args{
				field: "Rating",
				value: float64(5.5),
				min:   1,
				max:   5,
			},
			want: &Error{
				Field: "Rating",
				Msg:   "must be in range [1.000000, 5.000000]",
			},
		},
		{
			name: "unsupported type",
			args: args{
				field: "Tags",
				value: []string{"tag1", "tag2"},
				min:   1,
				max:   5,
			},
			want: &Error{
				Field: "Tags",
				Msg:   "unsupported type []string for InRange",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InRange(
				tt.args.field,
				tt.args.value,
				tt.args.min,
				tt.args.max,
			)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMinLength(t *testing.T) {
	type args struct {
		field     string
		value     interface{}
		minLength int
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				field:     "Title",
				value:     nil,
				minLength: 5,
			},
			want: &Error{Field: "Title", Msg: "cannot be nil"},
		},
		{
			name: "negative minLength",
			args: args{
				field:     "Title",
				value:     "hello",
				minLength: -1,
			},
			want: &Error{Field: "Title", Msg: "minLength must be non-negative"},
		},
		{
			name: "string valid",
			args: args{
				field:     "Title",
				value:     "hello",
				minLength: 5,
			},
			want: nil,
		},
		{
			name: "string invalid",
			args: args{
				field:     "Title",
				value:     "hello",
				minLength: 6,
			},
			want: &Error{
				Field: "Title",
				Msg:   "must have a minimum length of 6",
			},
		},
		{
			name: "slice valid",
			args: args{
				field:     "Tags",
				value:     []string{"tag1", "tag2"},
				minLength: 1,
			},
			want: nil,
		},
		{
			name: "slice invalid",
			args: args{
				field:     "Tags",
				value:     []string{"tag1", "tag2"},
				minLength: 3,
			},
			want: &Error{Field: "Tags", Msg: "must have a minimum length of 3"},
		},
		{
			name: "map valid",
			args: args{
				field:     "Lookup",
				value:     map[string]string{"foo": "bar"},
				minLength: 1,
			},
			want: nil,
		},
		{
			name: "map invalid",
			args: args{
				field:     "Lookup",
				value:     map[string]string{"foo": "bar"},
				minLength: 2,
			},
			want: &Error{
				Field: "Lookup",
				Msg:   "must have a minimum length of 2",
			},
		},
		{
			name: "unsupported type",
			args: args{
				field:     "Number",
				value:     123,
				minLength: 2,
			},
			want: &Error{
				Field: "Number",
				Msg:   "unsupported type int for MinLength",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinLength(tt.args.field, tt.args.value, tt.args.minLength)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMaxLength(t *testing.T) {
	type args struct {
		field     string
		value     interface{}
		maxLength int
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				field:     "Title",
				value:     nil,
				maxLength: 5,
			},
			want: &Error{Field: "Title", Msg: "cannot be nil"},
		},
		{
			name: "negative maxLength",
			args: args{
				field:     "Title",
				value:     "hello",
				maxLength: -1,
			},
			want: &Error{Field: "Title", Msg: "maxLength must be non-negative"},
		},
		{
			name: "string valid",
			args: args{
				field:     "Title",
				value:     "hello",
				maxLength: 5,
			},
			want: nil,
		},
		{
			name: "string invalid",
			args: args{
				field:     "Title",
				value:     "hello",
				maxLength: 4,
			},
			want: &Error{
				Field: "Title",
				Msg:   "must have a maximum length of 4",
			},
		},
		{
			name: "slice valid",
			args: args{
				field:     "Tags",
				value:     []string{"tag1", "tag2"},
				maxLength: 2,
			},
			want: nil,
		},
		{
			name: "slice invalid",
			args: args{
				field:     "Tags",
				value:     []string{"tag1", "tag2"},
				maxLength: 1,
			},
			want: &Error{Field: "Tags", Msg: "must have a maximum length of 1"},
		},
		{
			name: "map valid",
			args: args{
				field:     "Lookup",
				value:     map[string]string{"foo": "bar"},
				maxLength: 1,
			},
			want: nil,
		},
		{
			name: "map invalid",
			args: args{
				field:     "Lookup",
				value:     map[string]string{"foo": "bar"},
				maxLength: 0,
			},
			want: &Error{
				Field: "Lookup",
				Msg:   "must have a maximum length of 0",
			},
		},
		{
			name: "unsupported type",
			args: args{
				field:     "Number",
				value:     42,
				maxLength: 5,
			},
			want: &Error{
				Field: "Number",
				Msg:   "unsupported type int for MaxLength",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxLength(tt.args.field, tt.args.value, tt.args.maxLength)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatchRegexp(t *testing.T) {
	usernameRegexp := regexp.MustCompile(`^[a-z]+\d+$`)
	passwordRegexp := regexp.MustCompile(`^.*[A-Z]+.*$`)

	type args struct {
		field   string
		value   interface{}
		pattern *regexp.Regexp
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "string matches pattern",
			args: args{
				field:   "username",
				value:   "johndoe123",
				pattern: usernameRegexp,
			},
			want: nil,
		},
		{
			name: "string does not match pattern",
			args: args{
				field:   "username",
				value:   "JohnDoe123",
				pattern: usernameRegexp,
			},
			want: &Error{
				Field: "username",
				Msg:   "does not match pattern '^[a-z]+\\d+$': 'JohnDoe123'",
			},
		},
		{
			name: "byte slice matches pattern",
			args: args{
				field:   "username",
				value:   []byte("johndoe123"),
				pattern: usernameRegexp,
			},
			want: nil,
		},
		{
			name: "byte slice does not match pattern",
			args: args{
				field:   "username",
				value:   []byte("JohnDoe123"),
				pattern: usernameRegexp,
			},
			want: &Error{
				Field: "username",
				Msg:   "does not match pattern '^[a-z]+\\d+$': 'JohnDoe123'",
			},
		},
		{
			name: "unsupported type",
			args: args{
				field:   "password",
				value:   123456,
				pattern: passwordRegexp,
			},
			want: &Error{
				Field: "password",
				Msg:   "unsupported type int for MatchRegexp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchRegexp(tt.args.field, tt.args.value, tt.args.pattern)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNotNil(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil pointer",
			args: args{
				field: "Name",
				value: (*string)(nil),
			},
			want: &Error{
				Field: "Name",
				Msg:   "must not be nil",
			},
		},
		{
			name: "non-nil pointer",
			args: args{
				field: "Name",
				value: new(string),
			},
			want: nil,
		},
		{
			name: "nil slice",
			args: args{
				field: "Tags",
				value: []string(nil),
			},
			want: &Error{
				Field: "Tags",
				Msg:   "must not be nil",
			},
		},
		{
			name: "non-nil slice",
			args: args{
				field: "Tags",
				value: []string{},
			},
			want: nil,
		},
		{
			name: "nil map",
			args: args{
				field: "Metadata",
				value: map[string]string(nil),
			},
			want: &Error{
				Field: "Metadata",
				Msg:   "must not be nil",
			},
		},
		{
			name: "non-nil map",
			args: args{
				field: "Metadata",
				value: map[string]string{},
			},
			want: nil,
		},
		{
			name: "nil interface",
			args: args{
				field: "Data",
				value: fmt.Stringer(nil),
			},
			want: &Error{Field: "Data", Msg: "must not be nil"},
		},
		{
			name: "non-nil interface",
			args: args{
				field: "Data",
				// Using *bytes.Buffer as an example of a non-nil fmt.Stringer.
				value: fmt.Stringer(bytes.NewBuffer(nil)),
			},
			want: nil,
		},
		{
			name: "nil function",
			args: args{
				field: "Callback",
				value: (func())(nil),
			},
			want: &Error{
				Field: "Callback",
				Msg:   "must not be nil",
			},
		},
		{
			name: "non-nil function",
			args: args{
				field: "Callback",
				value: func() {},
			},
			want: nil,
		},
		{
			name: "non-nil struct",
			args: args{
				field: "Person",
				value: struct{ Name string }{Name: "Alice"},
			},
			want: nil,
		},
		{
			name: "pointer to nil pointer",
			args: args{
				field: "Data",
				value: pointerToPointer(nil),
			},
			want: &Error{
				Field: "Data",
				Msg:   "must not be nil",
			},
		},
		{
			name: "pointer to non-nil pointer",
			args: args{
				field: "Data",
				value: pointerToPointer(new(int)),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NotNil(tt.args.field, tt.args.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func pointerToPointer(v *int) **int {
	return &v
}
