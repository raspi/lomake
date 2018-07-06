package lomake

import (
	"reflect"
	"strings"
)

// Generic structure
type formFieldDescription struct {
	Name        string `json:"name,"`                 // Field's name, for example "phonenumber"
	Description string `json:"description,"`          // Field's description, for example "Phone number"
	FieldType   string `json:"type,"`                 // Field's type, for example "input"
	Required    bool   `json:"req"`                   // Is field required y/n?
	Placeholder string `json:"placeholder,omitempty"` // Field's placeholder description, for example "Your phone number"
	Value       string `json:"value,omitempty"`       // Field's value, for example initial value "555-1234"
}

func newFormFieldDescription() formFieldDescription {
	return formFieldDescription{
		FieldType:   `input`,
		Required:    true,
		Placeholder: ``,
		Value:       ``,
	}
}

func overrideFieldTypes(replaceMap map[string]string, sd structureDescription) structureDescription {
	for idx, item := range sd.fields {
		sd.fields[idx] = overrideFieldType(replaceMap, item)
	}

	return sd
}

func overrideFieldType(replaceMap map[string]string, ffd formFieldDescription) formFieldDescription {
	ffd.FieldType = replaceMap[ffd.FieldType]
	return ffd
}

func (ffd *formFieldDescription) readJsonTag(tag reflect.StructTag) {
	// Parse json tag
	jsonTagValues, err := readStructTag("json", tag)
	if err == nil {
		for idx, t := range jsonTagValues {
			t = strings.TrimSpace(t)

			if t == "" {
				continue
			}

			if t == "omitempty" && ffd.Required {
				ffd.Required = false
			}

			if idx == 0 {
				ffd.Name = t
			}
		}
	}

}
