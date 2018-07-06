package _example

import (
	"html/template"
	"log"
	"fmt"
	"reflect"
	"bytes"
	"golang.org/x/text/message"
	"golang.org/x/text/language"
	"github.com/raspi/lomake"
)

type RegisterForm struct {
	Username     string `json:"," lomaketype:"input.text"`
	EmailAddress string `json:"," lomaketype:"input.text"`
	Password     string `json:"," lomaketype:"input.password"`
	Password2    string `json:"," lomaketype:"input.password"`
}

var pageTemplate = `
<form id="pageform" action="/api/account/register.json" method="post">
	{{- .Form -}}
    <input id="formsubmit" type="submit" class="btn-primary" value="{{ T "form.submit" }}" />
</form>
`

// Get HTML
func (f RegisterForm) HTML() template.HTML {
	// TODO cache
	out, err := lomake.New(&f)
	if err != nil {
		log.Fatalf(`error=%v`, err)
		return ``
	}

	return out
}

func main() {

	// -- Global:
	translator := message.NewPrinter(language.Finnish)

	t := template.New("")

	t = t.Funcs(template.FuncMap{
		"T": func(s string, a ...interface{}) string {
			ref := message.Key(s, fmt.Sprintf(`NOT TRANSLATED: '%v'`, s))
			return translator.Sprintf(ref, a...)
		},
	})

	t.Parse(pageTemplate)

	lomake.Translator = translator
	lomake.HTMLTemplate = t

	// -- Page (view):
	var buf bytes.Buffer
	page, err := t.Clone()

	// Render form
	var form RegisterForm

	out, err := lomake.New(&form)

	if err != nil {
		log.Fatalf(`error: %v`, err)
	}

	view := struct {
		Form template.HTML
	}{
		Form: out,
	}

	page.Execute(&buf, &view)

	fmt.Println(buf.String())
}

// Ignore, used to get the package name
type Empty struct{}

// Initialize translations
func init() {
	pkgName := reflect.TypeOf(Empty{}).PkgPath()

	prefix := fmt.Sprintf(`lomake.%v`, pkgName)

	t := make(map[string]map[language.Tag]string)

	t[fmt.Sprintf(`%v.Username`, prefix)] = map[language.Tag]string{
		language.English: `User name`,
		language.Finnish: `Käyttäjätunnus`,
	}

	t[fmt.Sprintf(`%v.Username.placeholder`, prefix)] = map[language.Tag]string{
		language.English: `Enter your user name`,
		language.Finnish: `Syötä käyttäjätunnuksesi`,
	}


	t[fmt.Sprintf(`%v.EmailAddress`, prefix)] = map[language.Tag]string{
		language.English: `E-mail address`,
		language.Finnish: `Sähköpostiosoite`,
	}

	t[fmt.Sprintf(`%v.EmailAddress.placeholder`, prefix)] = map[language.Tag]string{
		language.English: `Enter your e-mail`,
		language.Finnish: `Syötä sähköpostiosoitteesi`,
	}

	t[fmt.Sprintf(`%v.Password`, prefix)] = map[language.Tag]string{
		language.English: `Password`,
		language.Finnish: `Salasana`,
	}

	t[fmt.Sprintf(`%v.Password2`, prefix)] = map[language.Tag]string{
		language.English: `Password (again)`,
		language.Finnish: `Salasana (uudestaan)`,
	}

	t[fmt.Sprintf(`%v.Password.placeholder`, prefix)] = map[language.Tag]string{
		language.English: `%%mYs3cr37p455w0rd`,
		language.Finnish: `sAl454N4%%`,
	}

	t[fmt.Sprintf(`%v.Password2.placeholder`, prefix)] = map[language.Tag]string{
		language.English: t[fmt.Sprintf(`%v.Password.placeholder`, prefix)][language.English],
		language.Finnish: t[fmt.Sprintf(`%v.Password.placeholder`, prefix)][language.Finnish],
	}

	t[`form.submit`] = map[language.Tag]string{
		language.English: `Send`,
		language.Finnish: `Lähetä`,
	}

	for key,_ := range t {
		for lang, v := range t[key] {
			message.SetString(lang, key, v )
		}
	}
}