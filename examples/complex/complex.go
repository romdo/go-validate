package main

import (
	"fmt"

	"github.com/romdo/go-validate"
)

type Manifest struct {
	Spec *Spec `json:"spec"`
}

func (s *Manifest) Validate() error {
	return validate.RequireField("Spec", s.Spec)
}

type Spec struct {
	Containers []*Container `json:"containers"`
	Images     []*Image     `json:"images"`
}

func (s *Spec) Validate() error {
	var errs error

	if len(s.Containers) == 0 {
		errs = validate.AppendFieldError(errs,
			"Containers", "must contain at least one item",
		)
	} else {
		imgs := map[string]bool{}
		for _, img := range s.Images {
			if img.Name != "" {
				imgs[img.Name] = true
			}
		}
		for i, c := range s.Containers {
			if c.ImageRef != "" && !imgs[c.ImageRef] {
				errs = validate.Append(errs, &validate.Error{
					Field: fmt.Sprintf("containers.%d.imageRef", i),
					Msg: fmt.Sprintf(
						"image with name '%s' not found", c.ImageRef,
					),
				})
			}
		}
	}

	if len(s.Images) == 0 {
		errs = validate.AppendFieldError(errs,
			"Images", "must contain at least one item",
		)
	}

	return errs
}

type Container struct {
	Name     string `json:"name"`
	ImageRef string `json:"imageRef"`
}

func (s *Container) Validate() error {
	var errs error
	errs = validate.Append(errs, validate.RequireField("Name", s.Name))
	errs = validate.Append(errs, validate.RequireField("ImageRef", s.ImageRef))

	return errs
}

type Image struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
	Tag  string `json:"tag"`
}

func (s *Image) Validate() error {
	var errs error
	errs = validate.Append(errs, validate.RequireField("Name", s.Name))
	errs = validate.Append(errs, validate.RequireField("URI", s.URI))
	errs = validate.Append(errs, validate.RequireField("Tag", s.Tag))

	return errs
}

func main() {
	manifest := &Manifest{
		Spec: &Spec{
			Containers: []*Container{
				{
					ImageRef: "server",
				},
				{
					Name:     "worker",
					ImageRef: "myServer",
				},
			},
			Images: []*Image{
				{
					Name: "server",
				},
			},
		},
	}

	errs := validate.Validate(manifest)

	for _, err := range validate.Errors(errs) {
		fmt.Println(err)
	}
}
