package validate

import (
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
