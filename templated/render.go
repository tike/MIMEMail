//Package templated implements a lightweight mail templating mechanism.
//
// Store your templates in dir "A", then Render will execute "A/base/index.html" as the main template
// for the mail body, adding any other files in "A/base" as helper templates, stylesheets, images etc.
// Other subfolders of "A", e.g. "A/en_US" should contain one file per templated message.
// These files should define 2 templates:
// "subject" and "body" (using the {{ define 'name"}} ... {{end}} syntax).
//
// See the tests and the "example" folder for a demonstration.
package templated

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
)

// Config helps you to manage your templates.
type Config struct {
	// Dir is the directory where your templates live. It's supposed to have a
	// subfolder named "base" in which the Render method will look for the main and
	// helper templates.
	Dir string

	// Lang names the language subfolder of Dir in which the file containing
	// the message specific templates named "subject" and "body" are located.
	Lang string
}

// Render renders the template c.Dir/c.Lang/name in the context of c.Dir/base/index.html.
//
// See the tests and the "example" folder for a demonstration.
func Render(c *Config, name string, data interface{}) (string, []byte, error) {
	tmpl, err := template.ParseGlob(filepath.Join(c.Dir, "base", "*"))
	if err != nil {
		return "", nil, err
	}

	tmpl, err = tmpl.ParseFiles(filepath.Join(c.Dir, c.Lang, name))
	if err != nil {
		return "", nil, err
	}

	var subj strings.Builder
	if err := tmpl.ExecuteTemplate(&subj, "subject", data); err != nil {
		return "", nil, err
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "index.html", data); err != nil {
		return "", nil, err
	}

	return subj.String(), body.Bytes(), nil
}
