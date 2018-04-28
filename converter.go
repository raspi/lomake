package lomake

import (
	"reflect"
	"errors"
	"fmt"
	"strings"
)

type StructureDescription struct {
	Title             string                 `json:",omitempty"`
	Description       string                 `json:",omitempty"`
	InfoIcon          string                 `json:",omitempty"`
	InfoUrl           string                 `json:",omitempty"`
	Fields            []FormFieldDescription `json:","`
	SubmitDescription string                 `json:","`
}

func NewStructureDescription() StructureDescription {
	return StructureDescription{
		Title:             "",
		Description:       "",
		Fields:            []FormFieldDescription{},
		SubmitDescription: "default.form.submit",
	}
}

// Get the tag from struct `tagName:"tagValue"`
func ReadStructTag(tagName string, tag reflect.StructTag) (values []string, err error) {
	value, ok := tag.Lookup(tagName)

	if !ok {
		return nil, errors.New("Not found")
	}

	return strings.Split(value, ","), nil
}

func ConvertStructToFieldDescription(field reflect.StructField) (ffd FormFieldDescription, err error) {
	ffd.Name = field.Name
	ffd.Type = strings.ToLower(field.Type.Name())

	// Get type
	typetagvalues, err := ReadStructTag("type", field.Tag)
	if err != nil {
		return ffd, err
	}

	if len(typetagvalues) == 1 {
		ffd.Type = typetagvalues[0]
	}

	// JSON tag
	jsontagvalues, err := ReadStructTag("json", field.Tag)
	if err != nil {
		return ffd, err
	}
	for idx, item := range jsontagvalues {
		if idx == 0 && item != "" {
			ffd.Name = item
		}

		if item == "omitempty" {
			ffd.Required = true
		}
	}

	return ffd, nil
}

// Convert struct into StructureDescription
func ReadStructDescription(i interface{}) (form StructureDescription, err error) {
	if reflect.TypeOf(i).Elem().Kind() != reflect.Struct {
		return form, errors.New(fmt.Sprintf("Not a struct"))
	}

	structName := reflect.TypeOf(i).Elem()

	var fields []FormFieldDescription

	t := reflect.TypeOf(i).Elem()

	for i := 0; i < t.NumField(); i++ {
		ffd, err := ConvertStructToFieldDescription(t.Field(i))
		if err != nil {
			return form, err
		}

		ffd.Description = fmt.Sprintf("form.%s.%s", structName, ffd.Name)
		ffd.Placeholder = fmt.Sprintf("form.%s.%s.placeholder", structName, ffd.Name)
		ffd.Required = true

		fields = append(fields, ffd)
	}

	form.Fields = fields

	return form, nil
}
