package wut

import (
	"fmt"
	"strings"
	"text/template"
)

type (
	// Template transforms the given context into a string
	Template func(ctx any) (string, error)

	// TemplateFactory provides templates based on the lang, key and value in the config
	// This allows setting custom template resolving when needed
	TemplateFactory interface {
		// GetTemplate - provides the template
		GetTemplate(lang string, path []string, entry string) (Template, error)
	}

	// DefaultTemplateFactory - provides `text/template` templates if it detects a placeholder,
	// or directly returns the string from the config
	DefaultTemplateFactory struct{}
)

func (f *DefaultTemplateFactory) GetTemplate(lang string, path []string, entry string) (Template, error) {

	if strings.Contains(entry, "{{") {
		key := fmt.Sprintf("[%s]%s", lang, strings.Join(path, "."))
		tmpl, err := template.New(key).
			Funcs(template.FuncMap{"get": func(key string, ctx ...any) (string, error) {
				return Lang(lang).Get(key, ctx...).OrErr()
			}}).
			Parse(entry)
		if err != nil {
			return nil, err
		}
		return func(ctx any) (string, error) {
			builder := &strings.Builder{}
			err := tmpl.Execute(builder, ctx)
			if err != nil {
				return "", err
			}
			return builder.String(), nil
		}, nil
	}

	return func(ctx any) (string, error) {
		return entry, nil
	}, nil
}
