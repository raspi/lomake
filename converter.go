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
	if len(typetagvalues) == 1 {
		ffd.Type = typetagvalues[0]
	}

	return ffd, nil
}

func ReadStructDescription(i interface{}) (form StructureDescription, err error) {
	if reflect.TypeOf(i).Elem().Kind() != reflect.Struct {
		return form, errors.New(fmt.Sprintf("Not a struct"))
	}

	structName := reflect.TypeOf(i).Elem()

	var j []FormFieldDescription

	t := reflect.TypeOf(i).Elem()

	for i := 0; i < t.NumField(); i++ {
		ffd, err := ConvertStructToFieldDescription(t.Field(i))
		if err != nil {
			return form, err
		}

		ffd.Description = fmt.Sprintf("form.%s.%s", structName, ffd.Name)

		j = append(j, ffd)
	}

	form.Fields = j

	return form, nil
}
