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
	Lang string
}

// Render renders a template mail to the given account
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
