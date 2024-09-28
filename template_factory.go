package joran

import (
	"fmt"
	"strings"
	"text/template"
)

type (
	Template func(ctx any) string

	TemplateFactory interface {
		GetTemplate(lang string, path []string, entry string) (Template, error)
	}

	DefaultTemplateFactory struct{}
)

func (f *DefaultTemplateFactory) GetTemplate(lang string, path []string, entry string) (Template, error) {

	if strings.Contains(entry, "{{") {
		key := fmt.Sprintf("[%s]%s", lang, strings.Join(path, "."))
		tmpl, err := template.New(key).Parse(entry)
		if err != nil {
			return nil, err
		}
		return func(ctx any) string {
			builder := &strings.Builder{}
			err := tmpl.Execute(builder, ctx)
			if err != nil {
				return "<err:" + key + ">"
			}
			return builder.String()
		}, nil
	}

	return func(ctx any) string {
		return entry
	}, nil
}
