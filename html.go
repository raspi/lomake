package lomake

import (
	"html/template"
	"bytes"
	"fmt"
	"errors"
	"golang.org/x/text/message"
)

// Template for generating HTML
var HTMLTemplate *template.Template

type decorator struct {
	Item   formFieldDescription
	Parent template.HTML
}

type DecoratorType uint64

//
const (
	DIV            DecoratorType = iota
	LABEL
	TEXTAREA
	INPUT_PASSWORD
	INPUT_TEXT
	INPUT_HIDDEN
)

// Base HTML templates
// .Parent is the parent decorator
var DecoratorTemplates = map[DecoratorType]string{
	DIV:            `<div class="form-group{{ if .Item.Required }} required{{ end }}">{{- .Parent -}}</div>`,
	LABEL:          `<label for="{{- .Item.Name -}}">{{ if .Item.Required }}* {{ end }}{{- T .Item.Description -}}</label>{{- .Parent -}}`,
	TEXTAREA:       `{{- .Parent -}}<textarea class="form-control" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" placeholder="{{- T .Item.Placeholder -}}">{{- .Item.Value -}}</textarea>`,
	INPUT_PASSWORD: `{{- .Parent -}}<input class="form-control" type="password" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" placeholder="{{- T .Item.Placeholder -}}" />`,
	INPUT_TEXT:     `{{- .Parent -}}<input class="form-control" type="text" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" placeholder="{{- T .Item.Placeholder -}}" />`,
	INPUT_HIDDEN:   `{{- .Parent -}}<input type="hidden" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" />`,
}

// Order of how decorators are applied
var DecoratorChains = map[DecoratorType][]DecoratorType{
	INPUT_TEXT:     {INPUT_TEXT, LABEL, DIV},
	INPUT_PASSWORD: {INPUT_PASSWORD, LABEL, DIV},
	INPUT_HIDDEN:   {INPUT_HIDDEN},
	TEXTAREA:       {TEXTAREA, LABEL, DIV},
}

// map string into decorator
var DecoratorMap = map[string]DecoratorType{
	"input.text":     INPUT_TEXT,
	"input.password": INPUT_PASSWORD,
	"input.hidden":   INPUT_HIDDEN,
	"textarea":       TEXTAREA,
}

// Apply decorator to a single field
func applyDecorator(field formFieldDescription) (output []byte, err error) {

	var tmpWriter bytes.Buffer
	var out bytes.Buffer

	var dec decorator
	dec.Parent = ""

	// Apply decorators in order
	for _, decorator := range DecoratorChains[DecoratorMap[field.FieldType]] {
		maintpl, err := HTMLTemplate.Clone()

		if err != nil {
			return nil, err
		}

		// Register template functions
		maintpl = maintpl.Funcs(template.FuncMap{
			"T": func(s string, a ...interface{}) string {
				// Translator
				ref := message.Key(s, fmt.Sprintf(`NOT TRANSLATED: '%v' (lomake)`, s))
				return Translator.Sprintf(ref, a...)
			},
		})

		tpl, err := maintpl.Parse(DecoratorTemplates[decorator])
		if err != nil {
			return nil, errors.New(fmt.Sprintf(`lomake template parsing error: %v`, err))
		}

		tmpWriter.Reset()
		dec.Item = field

		err = tpl.Execute(&tmpWriter, &dec)
		if err != nil {
			return nil, err
		}

		dec.Parent = template.HTML(tmpWriter.String())
	}

	_, err = out.Write(tmpWriter.Bytes())
	if err != nil {
		return nil, err
	}

	return out.Bytes(), err
}

// Apply decorators to all fields
func applyDecorators(fields []formFieldDescription) (output []byte, err error) {
	var out bytes.Buffer

	for _, item := range fields {
		output, err := applyDecorator(item)
		if err != nil {
			return nil, err
		}
		_, err = out.Write(output)
		if err != nil {
			return nil, err
		}

		_, err = out.Write([]byte("\n"))
		if err != nil {
			return nil, err
		}

	}

	return out.Bytes(), nil
}

// Convert struct into HTML form
func New(iface interface{}) (h template.HTML, err error) {
	form, err := readStructDescription(iface)
	if err != nil {
		return h, err
	}

	out, err := applyDecorators(form.fields)
	if err != nil {
		return h, err
	}

	h = template.HTML(out)

	return h, nil
}
