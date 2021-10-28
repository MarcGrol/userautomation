package actions

import (
	"bytes"
	"fmt"
	"html/template"
)

func ApplyTemplate(name, templateString string, data map[string]interface{}) (string, error) {
	templ, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("Error creating new template: %s", err)
	}

	buf := bytes.Buffer{}
	err = templ.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("Error executing new template: %s", err)
	}

	return buf.String(), nil
}
