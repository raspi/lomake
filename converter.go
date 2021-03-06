package lomake

import (
	"reflect"
	"errors"
	"fmt"
	"strings"
	"golang.org/x/text/message"
)

type structureDescription struct {
	title             string                 `json:",omitempty"`
	description       string                 `json:",omitempty"`
	infoIcon          string                 `json:",omitempty"`
	infoUrl           string                 `json:",omitempty"`
	fields            []formFieldDescription `json:","`
	submitDescription string                 `json:","`
}

const (
	LOMAKE_TYPE_FIELD = `lomaketype`
)

// Translate form labels, placeholders, ..
var Translator *message.Printer

func NewStructureDescription() structureDescription {
	return structureDescription{
		title:             "",
		description:       "",
		fields:            []formFieldDescription{},
		submitDescription: "default.form.submit",
	}
}

// Get the tag from struct `tagName:"tagValue"`
func readStructTag(tagName string, tag reflect.StructTag) (values []string, err error) {
	value, ok := tag.Lookup(tagName)

	if !ok {
		return nil, errors.New(fmt.Sprintf(`Not found: %v`, tagName))
	}

	return strings.Split(value, ","), nil
}

func convertStructToFieldDescription(field reflect.StructField) (ffd formFieldDescription, err error) {
	ffd.Name = field.Name
	ffd.FieldType = strings.ToLower(field.Type.Name())

	// Get type
	typetagvalues, err := readStructTag(LOMAKE_TYPE_FIELD, field.Tag)
	if err != nil {
		return ffd, err
	}

	if len(typetagvalues) == 1 {
		ftype := typetagvalues[0]
		_, ok := DecoratorMap[ftype]

		if !ok {
			return ffd, errors.New(fmt.Sprintf(`unknown lomake field type: '%v'`, ftype))
		}

		ffd.FieldType = ftype
	}

	// JSON tag
	jsontagvalues, err := readStructTag("json", field.Tag)
	if err != nil {
		return ffd, err
	}

	for idx, item := range jsontagvalues {
		// Foo string `json:"foo,omitempty"` -> name = foo instead of Foo
		if idx == 0 && item != "" {
			ffd.Name = item
		}

		// foo string `json:",omitempty"` -> required
		if item == "omitempty" {
			ffd.Required = true
		}
	}

	return ffd, nil
}

// Convert struct into StructureDescription
func readStructDescription(iface interface{}) (form structureDescription, err error) {
	if reflect.TypeOf(iface).Elem().Kind() != reflect.Struct {
		return form, errors.New(fmt.Sprintf(fmt.Sprintf(`Not a struct: %[1]T %#[1]v`, iface)))
	}

	structName := strings.Split(reflect.TypeOf(iface).Elem().String(), `.`)[0]

	var fields []formFieldDescription

	t := reflect.TypeOf(iface).Elem()
	v := reflect.ValueOf(iface).Elem()

	for i := 0; i < t.NumField(); i++ {
		ffd, err := convertStructToFieldDescription(t.Field(i))
		if err != nil {
			return form, err
		}

		ffd.Description = fmt.Sprintf(`lomake.%v.%v`, structName, ffd.Name)
		ffd.Placeholder = fmt.Sprintf(`lomake.%v.%v.placeholder`, structName, ffd.Name)
		ffd.Required = true
		ffd.Value = v.Field(i).String()

		fields = append(fields, ffd)
	}

	form.fields = fields

	return form, nil
}
