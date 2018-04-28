package lomake

import (
	"html/template"
	"bytes"
)

type Decorator struct {
	Item   FormFieldDescription
	Parent template.HTML
}

type DecoratorType uint64

const (
	DIV            DecoratorType = iota
	LABEL
	TEXTAREA
	INPUT_PASSWORD
	INPUT_TEXT
	INPUT_HIDDEN
)

//type DecoratorTemplate string

var DecoratorTemplates = map[DecoratorType]string{
	DIV:            `<div class="form-group{{ if .Item.Required }} required{{ end }}">{{- .Parent -}}</div>`,
	LABEL:          `<label for="{{- .Item.Name -}}">{{ if .Item.Required }}* {{ end }}{{- .Item.Description -}}</label>{{- .Parent -}}`,
	TEXTAREA:       `{{- .Parent -}}<textarea class="form-control" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" placeholder="{{- .Item.Placeholder -}}">{{- .Item.Value -}}</textarea>`,
	INPUT_PASSWORD: `{{- .Parent -}}<input class="form-control" type="password" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" placeholder="{{- .Item.Placeholder -}}" />`,
	INPUT_TEXT:     `{{- .Parent -}}<input class="form-control" type="text" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" placeholder="{{- .Item.Placeholder -}}" />`,
	INPUT_HIDDEN:   `{{- .Parent -}}<input type="hidden" name="{{- .Item.Name -}}" id="{{- .Item.Name -}}" value="{{- .Item.Value -}}" />`,
}

var DecoratorChains = map[DecoratorType][]DecoratorType{
	INPUT_TEXT:     {INPUT_TEXT, LABEL, DIV},
	INPUT_PASSWORD: {INPUT_PASSWORD, LABEL, DIV},
	INPUT_HIDDEN:   {INPUT_HIDDEN},
	TEXTAREA:       {TEXTAREA, LABEL, DIV},
}

var DecoratorMap = map[string]DecoratorType{
	"input.text":     INPUT_TEXT,
	"input.password": INPUT_PASSWORD,
	"input.hidden":   INPUT_HIDDEN,
	"textarea":       TEXTAREA,
}

// Apply decorator to a single field
func ApplyDecorator(field FormFieldDescription, decorators []DecoratorType) (output []byte, err error) {

	var tmpWriter bytes.Buffer
	var out bytes.Buffer

	var dec Decorator
	dec.Parent = ""

	for _, decorator := range decorators {
		tpl, err := template.New("").Parse(DecoratorTemplates[decorator])
		if err != nil {
			return nil, err
		}

		tmpWriter.Reset()
		dec.Item = field

		err = tpl.Execute(&tmpWriter, dec)
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
func ApplyDecorators(fields []FormFieldDescription, tpl *template.Template, decorators map[DecoratorType][]DecoratorType) (output []byte, err error) {

	if len(decorators) == 0 {
		decorators = DecoratorChains
	}

	var out bytes.Buffer

	for _, item := range fields {

		output, err := ApplyDecorator(item, decorators[DecoratorMap[item.Type]])
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
