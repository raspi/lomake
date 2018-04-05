package lomake

import (
	"reflect"
	"strings"
)

// Generic structure
type FormFieldDescription struct {
	Name        string `json:"name,"` // Field's name, for example "phonenumber"
	Description string `json:"description,"` // Field's description, for example "Phone number"
	Type        string `json:"type,"` // Field's type, for example "input"
	Required    bool   `json:"req"` // Is field required y/n?
	Placeholder string `json:"placeholder,omitempty"` // Field's placeholder description, for example "Your phone number"
	Value       string `json:"value,omitempty"` // Field's value, for example initial value "555-1234"
}

func NewFormFieldDescription() FormFieldDescription {
	return FormFieldDescription{
		Type:        "input",
		Required:    true,
		Placeholder: "",
		Value:       "",
	}
}

func OverrideFieldTypes(replaceMap map[string]string, sd StructureDescription) StructureDescription {
	for idx, item := range sd.Fields {
		sd.Fields[idx] = OverrideFieldType(replaceMap, item)
	}

	return sd
}

func OverrideFieldType(replaceMap map[string]string, ffd FormFieldDescription) FormFieldDescription {
	ffd.Type = replaceMap[ffd.Type]
	return ffd
}

func (ffd *FormFieldDescription) ReadJsonTag(tag reflect.StructTag) {
	// Parse json tag
	jsontagvalues, err := ReadStructTag("json", tag)
	if err == nil {
		for idx, t := range jsontagvalues {
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
