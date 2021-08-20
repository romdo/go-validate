// Package validate is yet another Go struct/object validation package, with a
// focus on simplicity, flexibility, and full control over validation logic.
//
// Interface
//
// To add validation to any type, simply implement the Validatable interface:
//
//  type Validatable interface {
//      Validate() error
//  }
//
// To mark a object as failing validation, the Validate method simply needs to
// return a error.
//
// When validating array, slice, map, and struct types each item and/or field
// that implements Validatable will be validated, meaning deeply nested structs
// can be fully validated, and the nested path to each object is tracked and
// reported back any validation errors.
//
// Multiple Errors
//
// Multiple errors can be reported from the Validate method using one of the
// available Append helper functions which append errors together. Under the
// hood the go.uber.org/multierr package is used to represent multiple errors as
// a single error return type, and you can in fact just directly use multierr in
// the a type's Validate method.
//
// Structs and Field-specific Errors
//
// When validating a struct, you are likely to have multiple errors for multiple
// fields. To specify which field on the struct the error relates to, you have
// to return a *validate.Error instead of a normal Go error type. For example:
//
//  type Book struct {
//      Title  string
//      Author string
//  }
//
//  func (s *Book) Validate() error {
//      var errs error
//
//      if s.Title == "" {
//          errs = validate.Append(errs, &validate.Error{
//              Field: "Title", Msg: "is required",
//          })
//      }
//
//      if s.Author == "" {
//          // Yields the same result as the Title field check above.
//          errs = validate.AppendFieldError(errs, "Author", "is required")
//      }
//
//      return errs
//  }
//
// With the above example, if you validate a empty *Book:
//
//  err := validate.Validate(&Book{})
//  for _, e := range validate.Errors(err) {
//      fmt.Println(e.Error())
//  }
//
// The following errors would be printed:
//
//  Title: is required
//  Kind: is required
//
// Error type
//
// All errors will be wrapped in a *Error before being returned, which is used
// to keep track of the path and field the error relates to. There are various
// helpers available to create Error instances.
//
// Handling Validation Errors
//
// As mentioned above, multiple errors are wrapped up into a single error return
// value using go.uber.org/multierr. You can access all errors individually with
// Errors(), which accepts a single error, and returns []error. The Errors()
// function is just wrapper around multierr.Errors(), so you could use that
// instead if you prefer.
//
// Struct Field Tags
//
// Fields on a struct which customize the name via a json, yaml, or form field
// tag, will automatically have the field name converted to the name in the tag
// in returned *Error types with a non-empty Field value.
//
// You can customize the field name conversion logic by creating a custom
// Validator instance, and calling FieldNameFunc() on it.
//
// Nested Validatable Objects
//
// All items/fields on any structs, maps, slices or arrays which are encountered
// will be validated if they implement the Validatable interface. While
// traversing nested data structures, a path list tracks the location of the
// current object being validation in relation to the top-level object being
// validated. This path is used within the field in the final output errors.
//
// By default path components are joined with a dot, but this can be customized
// when using a custom Validator instance and calling FieldJoinFunc() passing in
// a custom function to handle path joining.
//
// As an example, if our Book struct from above is nested within the following
// structs:
//
//  type Order struct {
//      Items []*Item `json:"items"`
//  }
//
//  type Item struct {
//      Book *Book `json:"book"`
//  }
//
// And we have a Order where the book in the second Item has a empty Author
// field:
//
//  err := validate.Validate(&Order{
//      Items: []*Item{
//          {Book: &Book{Title: "The Traveler", Author: "John Twelve Hawks"}},
//          {Book: &Book{Title: "The Firm"}},
//      },
//  })
//  for _, e := range validate.Errors(err) {
//      fmt.Println(e.Error())
//  }
//
// Then we would get the following error:
//
//  items.1.book.Author: is required
//
// Note how both "items" and "book" are lower cased thanks to the json tags on
// the struct fields, while our Book struct does not have a json tag for the
// Author field.
//
// Also note that the error message does not start with "Order". The field path
// is relative to the object being validated, hence the top-level object is not
// part of the returned field path.
package validate

// global is a private instance of Validator to enable the package root-level
// Validate() function.
var global = New()

// Validate will validate the given object. Structs, maps, slices, and arrays
// will have each of their fields/items validated, effectively performing a
// deep-validation.
func Validate(v interface{}) error {
	return global.Validate(v)
}

// Validatable is the primary interface that a object needs to implement to be
// validatable with Validator.
//
// Validation errors are reported by returning a error from the Validate
// method. Multiple errors can be combined into a single error to return with
// Append() and related functions, or via go.uber.org/multierr.
//
// For validatable structs, the field the validation error relates to can be
// specified by returning a *Error type with the Field value specified.
type Validatable interface {
	Validate() error
}
