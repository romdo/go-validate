package validate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Field string
		Msg   string
		Err   error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   "unknown error",
		},
		{
			name: "field only",
			fields: fields{
				Field: "spec.images.0.name",
			},
			want: "spec.images.0.name: unknown error",
		},
		{
			name: "msg only",
			fields: fields{
				Msg: "flux capacitor is missing",
			},
			want: "flux capacitor is missing",
		},
		{
			name: "err only",
			fields: fields{
				Err: errors.New("flux capacitor is king"),
			},
			want: "flux capacitor is king",
		},
		{
			name: "field and msg",
			fields: fields{
				Field: "spec.images.0.name",
				Msg:   "is required",
			},
			want: "spec.images.0.name: is required",
		},
		{
			name: "field and err",
			fields: fields{
				Field: "spec",
				Err:   errors.New("something is wrong"),
			},
			want: "spec: something is wrong",
		},
		{
			name: "msg and err",
			fields: fields{
				Msg: "flux capacitor is missing",
				Err: errors.New("flux capacitor is king"),
			},
			want: "flux capacitor is missing",
		},
		{
			name: "field, msg, and err",
			fields: fields{
				Field: "spec.images.0.name",
				Msg:   "is required",
				Err:   errors.New("something is wrong"),
			},
			want: "spec.images.0.name: is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{
				Field: tt.fields.Field,
				Msg:   tt.fields.Msg,
				Err:   tt.fields.Err,
			}

			got := err.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Is(t *testing.T) {
	errTest1 := errors.New("errtest1")
	errTest2 := errors.New("errtest2")

	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		target error
		want   bool
	}{
		{
			name:   "empty",
			fields: fields{},
			target: errTest1,
			want:   false,
		},
		{
			name:   "Err and target match",
			fields: fields{Err: errTest1},
			target: errTest1,
			want:   true,
		},
		{
			name:   "Err and target do not match",
			fields: fields{Err: errTest2},
			target: errTest1,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{Err: tt.fields.Err}

			got := errors.Is(err, tt.target)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	errTest1 := errors.New("errtest1")
	errTest2 := errors.New("errtest2")

	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   nil,
		},
		{
			name:   "Err test1",
			fields: fields{Err: errTest1},
			want:   errTest1,
		},
		{
			name:   "Err test2",
			fields: fields{Err: errTest2},
			want:   errTest2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{Err: tt.fields.Err}

			got := err.Unwrap()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAppend(t *testing.T) {
	type args struct {
		errs error
		err  error
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "append nil to nil",
			args: args{
				errs: nil,
				err:  nil,
			},
			want: []error{},
		},
		{
			name: "append nil to err",
			args: args{
				errs: errors.New("foo"),
				err:  nil,
			},
			want: []error{
				errors.New("foo"),
			},
		},
		{
			name: "append nil to multi err",
			args: args{
				errs: multierr.Combine(errors.New("foo"), errors.New("bar")),
				err:  nil,
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
		},
		{
			name: "append err to nil",
			args: args{
				errs: nil,
				err:  errors.New("foo"),
			},
			want: []error{
				errors.New("foo"),
			},
		},
		{
			name: "append err to err",
			args: args{
				errs: errors.New("foo"),
				err:  errors.New("bar"),
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
		},
		{
			name: "append err to multi err",
			args: args{
				errs: multierr.Combine(errors.New("foo"), errors.New("bar")),
				err:  errors.New("baz"),
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
				errors.New("baz"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Append(tt.args.errs, tt.args.err)

			if len(tt.want) == 0 {
				assert.Nil(t, err)
			}

			got := multierr.Errors(err)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestAppendError(t *testing.T) {
	type args struct {
		errs error
		msg  string
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "append empty msg to nil",
			args: args{
				errs: nil,
				msg:  "",
			},
			want: []error{
				&Error{},
			},
		},
		{
			name: "append msg to nil",
			args: args{
				errs: nil,
				msg:  "foo",
			},
			want: []error{
				&Error{Msg: "foo"},
			},
		},
		{
			name: "append msg to err",
			args: args{
				errs: errors.New("foo"),
				msg:  "bar",
			},
			want: []error{
				errors.New("foo"),
				&Error{Msg: "bar"},
			},
		},
		{
			name: "append msg to multi err",
			args: args{
				errs: multierr.Combine(errors.New("foo"), errors.New("bar")),
				msg:  "baz",
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
				&Error{Msg: "baz"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendError(tt.args.errs, tt.args.msg)

			if len(tt.want) == 0 {
				assert.Nil(t, err)
			}

			got := multierr.Errors(err)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestAppendFieldError(t *testing.T) {
	type args struct {
		errs  error
		field string
		msg   string
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "append empty field and msg to nil",
			args: args{
				errs:  nil,
				field: "",
				msg:   "",
			},
			want: []error{
				&Error{},
			},
		},
		{
			name: "append msg to nil",
			args: args{
				errs:  nil,
				field: "Type",
				msg:   "foo",
			},
			want: []error{
				&Error{Field: "Type", Msg: "foo"},
			},
		},
		{
			name: "append msg to err",
			args: args{
				errs:  errors.New("foo"),
				field: "Type",
				msg:   "bar",
			},
			want: []error{
				errors.New("foo"),
				&Error{Field: "Type", Msg: "bar"},
			},
		},
		{
			name: "append msg to multi err",
			args: args{
				errs:  multierr.Combine(errors.New("foo"), errors.New("bar")),
				field: "Type",
				msg:   "baz",
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
				&Error{Field: "Type", Msg: "baz"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendFieldError(tt.args.errs, tt.args.field, tt.args.msg)

			if len(tt.want) == 0 {
				assert.Nil(t, err)
			}

			got := multierr.Errors(err)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestErrors(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want []error
	}{
		{
			name: "nil",
			args: args{err: nil},
			want: nil,
		},
		{
			name: "single error",
			args: args{err: errors.New("foo")},
			want: []error{errors.New("foo")},
		},
		{
			name: "multi error with one error",
			args: args{
				err: multierr.Combine(errors.New("foo")),
			},
			want: []error{
				errors.New("foo"),
			},
		},
		{
			name: "multi error with two errors",
			args: args{
				err: multierr.Combine(errors.New("foo"), errors.New("bar")),
			},
			want: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Errors(tt.args.err)

			assert.Equal(t, tt.want, got)
		})
	}
}
