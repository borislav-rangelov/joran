package wut

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type (
	LangFactory interface {
		Lang(code string) LangSource
	}

	langFactory struct {
		sources     map[string]LangSource
		emptySource *LookupSource
	}
)

func NewLangFactory(tf TemplateFactory, files ...*LangFile) (LangFactory, error) {
	fileMap, fallbackMap, err := mergeLangFiles(files)

	if err != nil {
		return nil, err
	}

	sources, err := constructLangSourceChain(fallbackMap)

	for code, file := range fileMap {
		m, er := toConfigMap(file.Language, make([]string, 0), file.Values, tf)
		if er != nil {
			return nil, er
		}
		sources[code].(*LookupSource).configs = m
	}

	return &langFactory{
		sources:     sources,
		emptySource: &LookupSource{fallback: nil, configs: make(ConfigMap)},
	}, nil
}

func constructLangSourceChain(fallbackMap map[string]string) (map[string]LangSource, error) {
	values := maps.Clone(fallbackMap)
	valid := make(map[string]bool)
	result := make(map[string]LangSource)

	// a code->fallback is valid if there is no fallback ("")
	// or the fallback was already verified

	for {
		initial := len(values)

		for code, fallback := range values {

			if fallback == "" || valid[fallback] {
				valid[code] = true
				delete(values, code)

				ls := &LookupSource{}
				if fallbackLS, ok := result[fallback]; ok {
					ls.fallback = fallbackLS
				}
				result[code] = ls
			}
		}

		left := len(values)
		if left == 0 {
			return result, nil
		}
		if left == initial {
			return nil, errors.New(fmt.Sprintf("invalid fallback chain: %v", values))
		}
	}
}

func (l *langFactory) Lang(code string) LangSource {
	code = strings.ToLower(code)
	if source, ok := l.sources[code]; ok {
		return source
	}
	return l.emptySource
}

func mergeLangFiles(files []*LangFile) (map[string]*LangFile, map[string]string, error) {
	fileMap := make(map[string]*LangFile)
	fallbackMap := make(map[string]string)

	for _, f := range files {
		if _, ok := fileMap[f.Language]; !ok {
			fileMap[f.Language] = f
			fallbackMap[f.Language] = f.Fallback
		} else {
			err := fileMap[f.Language].Include(f)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return fileMap, fallbackMap, nil
}

func toConfigMap(lang string, path []string, values map[string]any, tf TemplateFactory) (ConfigMap, error) {
	result := make(ConfigMap)

	for k, v := range values {
		keyPath := slices.Clone(path)
		keyPath = append(keyPath, k)

		if entry, ok := v.(string); ok {
			template, err := tf.GetTemplate(lang, keyPath, entry)
			if err != nil {
				return nil, err
			}
			result[k] = &KeyConfig{
				Template: template,
			}
			continue
		}

		if cfg, ok := v.(map[string]any); ok {
			configMap, err := toConfigMap(lang, keyPath, cfg, tf)
			if err != nil {
				return nil, err
			}
			result[k] = &KeyConfig{
				Configs: configMap,
			}
			continue
		}

		return nil, fmt.Errorf("unknown value type: %T", v)
	}

	return result, nil
}
